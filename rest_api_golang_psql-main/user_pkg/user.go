package user_pkg

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"routing.go/Task1/con_pkg"

	"github.com/gorilla/mux"
)

// User schema of the user table
type User struct {
	U_id     int     `json:"U_id"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Deposit  float64 `json:"deposit"`
	Role     string  `json:"role"`
}

// handle all the db operations like Insert, Select, Update, and Delete (CRUD).

// [1] Get All Users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {

	// create the postgres db connection & execute the sql statement
	db := con_pkg.CreateConnection()
	defer db.Close()
	rows, err := db.Query(`SELECT *FROM public."User"`)
	con_pkg.CheckError(err)
	defer rows.Close()

	var user User
	var users []User
	// iterate over the rows
	for rows.Next() {
		// unmarshal the row object to user
		err = rows.Scan(&user.U_id, &user.Username, &user.Password, &user.Deposit, &user.Role)
		con_pkg.CheckError(err)
		// append the user in the users slice
		users = append(users, user)
	}
	var response = con_pkg.JsonResponse{Message: "All rows returned successfully", Data: users}
	json.NewEncoder(w).Encode(response)
}

// [2] Get one user using id
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	U_id := params["id"]
	//var response = JsonResponse{}
	var user User
	// create the postgres db connection
	db := con_pkg.CreateConnection()
	// close the db connection
	defer db.Close()
	sqlStatement := `SELECT *FROM public."User" WHERE "U_id"=$1`
	// Get One User from User table by U_id
	row := db.QueryRow(sqlStatement, U_id)
	// unmarshal the row object to user
	err := row.Scan(&user.U_id, &user.Username, &user.Password, &user.Deposit, &user.Role)

	// unmarshal the row object to user
	var response = con_pkg.JsonResponse{}
	switch err {
	case sql.ErrNoRows:
		response = con_pkg.JsonResponse{Message: "No rows were returned with ID [ " + U_id + " ]", Data: nil}

	case nil:
		response = con_pkg.JsonResponse{Message: "raw with ID [ " + U_id + " ] return successfully", Data: user}

	default:
		response = con_pkg.JsonResponse{Message: "Unable to scan the row", Data: nil}
	}
	json.NewEncoder(w).Encode(response)

}

// [3] Update user data using its id
func UpdatUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	//convert the string into an integer value.
	U_id, _ := strconv.Atoi(params["id"])
	// create an empty user of type User
	var user User
	// decode the json request to user
	err := json.NewDecoder(r.Body).Decode(&user)
	user.U_id = U_id
	con_pkg.CheckError(err)
	// create the postgres db connection
	db := con_pkg.CreateConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE public."User" SET  username =$2, password=$3, deposit=$4, role=$5  WHERE "U_id"=$1;`

	// execute the sql statement
	_, err = db.Exec(sqlStatement, user.U_id, user.Username, user.Password, user.Deposit, user.Role)
	con_pkg.CheckError(err)

	// format the response message
	response := con_pkg.JsonResponse{Message: "User updated successfully", Data: user}
	// send the response
	json.NewEncoder(w).Encode(response)

}

// [4] Creat new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	Username := r.FormValue("username")
	Password := r.FormValue("password")
	Deposit := r.FormValue("deposit")
	Role := r.FormValue("role")

	var response = con_pkg.JsonResponse{}

	if Username == "" || Password == "" || Deposit == "" || Role == "" {
		response = con_pkg.JsonResponse{Message: "Complete Missing Data!!!", Data: nil}
	} else {
		db := con_pkg.CreateConnection()
		// close the db connection
		defer db.Close()

		// create the insert sql query
		sqlStatement := `INSERT INTO public."User" (username, password, deposit, role) VALUES ($1,$2,$3,$4)  RETURNING  "U_id",username, password, deposit, role;`
		// the inserted data will store in this user
		var user User
		err := db.QueryRow(sqlStatement, Username, Password, Deposit, Role).Scan(&user.U_id, &user.Username, &user.Password, &user.Deposit, &user.Role)
		con_pkg.CheckError(err)
		response = con_pkg.JsonResponse{Message: "The User has been inserted successfully!", Data: user}

	}
	json.NewEncoder(w).Encode(response)

}

// [5] Delete user by id
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	U_id := params["id"]
	// create the postgres db connection
	db := con_pkg.CreateConnection()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM public."User" WHERE "U_id"=$1;`

	// execute the sql statement
	_, err := db.Exec(sqlStatement, U_id)
	con_pkg.CheckError(err)
	// format the response message
	response := con_pkg.JsonResponse{Message: "User with id = [" + U_id + "] is deleted successfully", Data: nil}
	// send the response
	json.NewEncoder(w).Encode(response)

}

// [6] Delete All Users
func DeleteAllUserS(w http.ResponseWriter, r *http.Request) {
	// create the postgres db connection
	db := con_pkg.CreateConnection()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM public."User";`

	// execute the sql statement
	_, err := db.Exec(sqlStatement)
	con_pkg.CheckError(err)
	var response = con_pkg.JsonResponse{Message: "All Users have been deleted successfully!", Data: nil}

	json.NewEncoder(w).Encode(response)

}
