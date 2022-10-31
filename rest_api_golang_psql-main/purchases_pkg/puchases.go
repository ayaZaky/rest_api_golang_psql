package purchases_pkg

import (
	"net/http"
	"strconv"

	"routing.go/Task1/con_pkg"     //connection_pkg
	"routing.go/Task1/product_pkg" //product_pkg
	"routing.go/Task1/user_pkg"    //user_pkg

	"encoding/json"

	"github.com/gorilla/mux"
)

// purchases schema of the   table
type Purchases struct {
	User_id        int `json:"user_id"`
	Product_Id     int `json:"product_id"`
	Product_amount int `json:"product_amount"`
}

// Buy endpoint
func BuyProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	U_id, _ := strconv.Atoi(params["id"])
	// get productid and product amount from r
	ProductId, _ := strconv.Atoi(r.FormValue("productId"))
	ProductAmount, _ := strconv.Atoi(r.FormValue("productAmount"))

	var response = con_pkg.JsonResponse{}
	var user user_pkg.User
	var product product_pkg.Product

	// create the postgres db connection & execute the sql statement
	db := con_pkg.CreateConnection()
	defer db.Close()
	sqlStatement1 := `SELECT * FROM public."User" WHERE "U_id"=$1;`
	sqlStatement2 := `SELECT *FROM public."Product" WHERE "P_id"=$1;`
	row1 := db.QueryRow(sqlStatement1, U_id)
	row2 := db.QueryRow(sqlStatement2, ProductId)
	// unmarshal the row1 object to user
	err := row1.Scan(&user.U_id, &user.Username, &user.Password, &user.Deposit, &user.Role)
	con_pkg.CheckError(err)
	// unmarshal the row2 object to product
	err1 := row2.Scan(&product.P_id, &product.AmountAvailable, &product.Cost, &product.ProductName, &product.Sellerid)
	//con_pkg.CheckError(err1)

	/// handel expections
	switch {
	case user.Role != "buyer":
		response = con_pkg.JsonResponse{Message: "This id doesnot belong to buyer !!!", Data: nil}
	case ProductId != product.P_id:
		response = con_pkg.JsonResponse{Message: "This Product Not Found !!!", Data: nil}

	case (ProductAmount > product.AmountAvailable):
		response = con_pkg.JsonResponse{Message: "Quantity required is greater than available !!!", Data: nil}

	case (user.Deposit) < product.Cost*float64(ProductAmount):
		response = con_pkg.JsonResponse{Message: "Deposit of this buyer isnot suffcient!!!", Data: nil}

	case err == nil && err1 == nil:

		//update data after buying process:

		new_Deposit := user.Deposit - (product.Cost * float64(ProductAmount))
		new_AvailableAmount := product.AmountAvailable - ProductAmount

		// create the update sql query
		sqlStatement1 := `UPDATE public."Product" SET  "amountAvailable"=$2 WHERE "P_id"=$1;`
		sqlStatement2 := `UPDATE public."User" SET  deposit=$2 WHERE "U_id"=$1 ;`

		// execute the sql statements
		_, err := db.Exec(sqlStatement1, ProductId, new_AvailableAmount)
		con_pkg.CheckError(err)
		_, err1 := db.Exec(sqlStatement2, U_id, new_Deposit)
		con_pkg.CheckError(err1)

		// insert purchases data in purchases Table in postgresql
		sqlStatement := `INSERT INTO public.purchases(user_id, product_id, product_amount) VALUES ($1,$2,$3)  RETURNING  user_id, product_id, product_amount;`
		// the inserted data will store in this purchases struct
		var purchases Purchases
		err2 := db.QueryRow(sqlStatement, U_id, ProductId, ProductAmount).Scan(&purchases.User_id, &purchases.Product_Id, &purchases.Product_amount)
		con_pkg.CheckError(err2)

		// format the response message
		response = con_pkg.JsonResponse{Message: "Purchase Completed successfully!", Data: purchases}

	default:
		response = con_pkg.JsonResponse{Message: "Unable to scan the row", Data: nil}
	}
	// send the response
	json.NewEncoder(w).Encode(response)

}
