package product_pkg

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"routing.go/Task1/con_pkg"

	"github.com/gorilla/mux"
)

// Product schema of the product table
type Product struct {
	P_id            int     `json:"P_id"`
	AmountAvailable int     `json:"AmountAvailable"`
	Cost            float64 `json:"Cost"`
	ProductName     string  `json:"ProductName"`
	Sellerid        int     `json:"Sellerid"`
}

// [1] Creat new Product
func CreateProduct(w http.ResponseWriter, r *http.Request) {

	AmountAvailable := r.FormValue("amountAvailable")
	Cost := r.FormValue("cost")
	ProductName := r.FormValue("productName")
	SellerId, _ := strconv.Atoi(r.FormValue("sellerId"))

	db := con_pkg.CreateConnection()
	defer db.Close()

	// check if  inserted seller_id is valid or no
	rows, err := db.Query(`SELECT "U_id" FROM public."User" WHERE role ='seller'`)
	con_pkg.CheckError(err)
	var id int
	var available_id []int
	// iterate over the rows
	for rows.Next() {
		err = rows.Scan(&id)
		// append the avilable IDs for selleres in  slice
		available_id = append(available_id, id)
	}
	var response = con_pkg.JsonResponse{}
	for i := 0; i < len(available_id); i++ {

		if available_id[i] == SellerId {
			// create the postgres db connection & execute the sql statement
			sqlStatement := `INSERT INTO public."Product" ("amountAvailable", cost, "productName", "sellerId") VALUES ($1, $2, $3, $4)
					RETURNING  "P_id" ,"amountAvailable", cost, "productName", "sellerId";`
			var product Product
			err := db.QueryRow(sqlStatement, AmountAvailable, Cost, ProductName, SellerId).Scan(&product.P_id,
				&product.AmountAvailable, &product.Cost, &product.ProductName, &product.Sellerid)
			con_pkg.CheckError(err)

			//json response
			response = con_pkg.JsonResponse{Message: "The Product has been inserted successfully!", Data: product}
			break
		} else {
			response = con_pkg.JsonResponse{Message: "These seller id isnont valid ,Available sellers in db are :", Data: available_id}
		}
	}
	//send response
	json.NewEncoder(w).Encode(response)

}

// [2] Get All products
func GetAllProducts(w http.ResponseWriter, r *http.Request) {

	// create the postgres db connection & execute the sql statement
	db := con_pkg.CreateConnection()
	defer db.Close()
	rows, err := db.Query(`SELECT *FROM public."Product"`)
	con_pkg.CheckError(err)
	defer rows.Close()

	var product Product
	var products []Product
	// iterate over the rows
	for rows.Next() {

		// unmarshal the row object to user
		err = rows.Scan(&product.P_id, &product.AmountAvailable, &product.Cost, &product.ProductName, &product.Sellerid)
		con_pkg.CheckError(err)

		// append the product in the products slice
		products = append(products, product)
	}
	var response = con_pkg.JsonResponse{Message: "All rows returned successfully", Data: products}
	json.NewEncoder(w).Encode(response)
}

// [3] Get specific product by its id
func GetProduct(w http.ResponseWriter, r *http.Request) {

	// get the product_id from the request params, key is "id"
	params := mux.Vars(r)
	P_id := params["id"]

	// create the postgres db connection & execute the sql statement
	db := con_pkg.CreateConnection()
	defer db.Close()
	sqlStatement := `SELECT *FROM public."Product" WHERE "P_id"=$1`
	row := db.QueryRow(sqlStatement, P_id)

	// unmarshal the row object to user
	var product Product
	err := row.Scan(&product.P_id, &product.AmountAvailable, &product.Cost, &product.ProductName, &product.Sellerid)

	// json response
	var response = con_pkg.JsonResponse{}
	switch err {
	case sql.ErrNoRows:
		response = con_pkg.JsonResponse{Message: "No rows were returned with ID [ " + P_id + " ]", Data: nil}

	case nil:
		response = con_pkg.JsonResponse{Message: "raw with ID [ " + P_id + " ] return successfully", Data: product}

	default:
		response = con_pkg.JsonResponse{Message: "Unable to scan the row", Data: nil}
	}
	json.NewEncoder(w).Encode(response)

}

// [4] Update product data using its id
func UpdateProduct(w http.ResponseWriter, r *http.Request) {

	// get the product_id from the request params, key is "id"
	params := mux.Vars(r)
	P_id, _ := strconv.Atoi(params["id"])
	var product Product
	// decode the json request to product
	err := json.NewDecoder(r.Body).Decode(&product)
	product.P_id = P_id

	// create the postgres db connection & execute the sql statement
	db := con_pkg.CreateConnection()
	defer db.Close()
	sqlStatement := `UPDATE public."Product" SET  "amountAvailable"=$2, cost=$3, "productName"=$4, "sellerId"=$5  WHERE "P_id"=$1;`
	_, err = db.Exec(sqlStatement, product.P_id, product.AmountAvailable, product.Cost, product.ProductName, product.Sellerid)
	con_pkg.CheckError(err)

	// format the response message
	response := con_pkg.JsonResponse{Message: "Product updated successfully", Data: product}
	// send the response
	json.NewEncoder(w).Encode(response)

}

// [5] Delete Product by id
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	P_id := params["id"]
	// create the postgres db connection & execute the sql statement
	db := con_pkg.CreateConnection()
	defer db.Close()
	sqlStatement := `DELETE FROM public."Product" WHERE "P_id"=$1; `
	_, err := db.Exec(sqlStatement, P_id)
	con_pkg.CheckError(err)

	// format the response message
	response := con_pkg.JsonResponse{Message: "User with id = [" + P_id + "] is deleted successfully", Data: nil}
	// send the response
	json.NewEncoder(w).Encode(response)

}

// [6] Delete All products
func DeleteAllProducts(w http.ResponseWriter, r *http.Request) {

	// create the postgres db connection & execute the sql statement
	db := con_pkg.CreateConnection()
	defer db.Close()
	sqlStatement := `DELETE FROM public."Product";`
	_, err := db.Exec(sqlStatement)
	con_pkg.CheckError(err)

	// format the response message
	var response = con_pkg.JsonResponse{Message: "All Products have been deleted successfully!", Data: nil}
	// send the response
	json.NewEncoder(w).Encode(response)

}
