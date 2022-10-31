package con_pkg

import (
	"database/sql"
	"fmt"

	// postgres golang driver
	_ "github.com/lib/pq"
)

// response format

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "ayazaky44"
	dbname   = "API_DB"
)

// create connection with postgres db
func CreateConnection() *sql.DB {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	CheckError(err)
	//defer db.Close()

	err = db.Ping()
	CheckError(err)

	fmt.Println("Successfully connected!")
	return db
}

// Function for handling errors
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

// Struct JsonResponse will display the JSON response once the data is fetched.
type JsonResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
