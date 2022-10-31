package main

import (
	"log"
	"net/http"

	//product _pkg
	"routing.go/Task1/product_pkg"
	//purchases _pkg
	"routing.go/Task1/purchases_pkg"
	//user _pkg
	"routing.go/Task1/user_pkg"

	"github.com/gorilla/mux"
)

func main() {

	// Router handles & endpoints
	router := mux.NewRouter()

	// User Handel functions
	router.HandleFunc("/user", user_pkg.CreateUser).Methods("POST")
	router.HandleFunc("/user/{id}", user_pkg.GetUser).Methods("GET")
	router.HandleFunc("/users", user_pkg.GetAllUsers).Methods("GET")
	router.HandleFunc("/user/{id}", user_pkg.UpdatUser).Methods("PATCH")
	router.HandleFunc("/user/{id}", user_pkg.DeleteUser).Methods("DELETE")
	router.HandleFunc("/users", user_pkg.DeleteAllUserS).Methods("DELETE")

	// Product Handle functions
	router.HandleFunc("/product", product_pkg.CreateProduct).Methods("POST")
	router.HandleFunc("/product/{id}", product_pkg.GetProduct).Methods("GET")
	router.HandleFunc("/products", product_pkg.GetAllProducts).Methods("GET")
	router.HandleFunc("/product/{id}", product_pkg.UpdateProduct).Methods("PATCH")
	router.HandleFunc("/product/{id}", product_pkg.DeleteProduct).Methods("DELETE")
	router.HandleFunc("products", product_pkg.DeleteAllProducts).Methods("DELETE")

	// Purchases Handel function
	router.HandleFunc("/buy/{id}", purchases_pkg.BuyProduct).Methods("POST")

	// Work server on the port 2020
	log.Fatal(http.ListenAndServe(":2020", router))
}
