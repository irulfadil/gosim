package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shopspring/decimal"
)

var db *gorm.DB
var err error

type Product struct {
	ID    int             `json:"id"`
	Code  string          `json:"code"`
	Name  string          `json:"name"`
	Price decimal.Decimal `json:"price" sql:"type:decimal(16,2);"`
}

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	db, err = gorm.Open("mysql", "root:@/go_restapi_crud?charset=utf8&parseTime=True")

	if err != nil {
		log.Println("Connection failed", err)
	} else {
		log.Println("Connection established")
	}

	db.AutoMigrate(&Product{})
	handleRequests()
}

func handleRequests() {
	log.Println("Started Server at http://127.0.0.1:9000")

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homepage)
	myRouter.HandleFunc("/api/product", createProduct).Methods("POST")
	myRouter.HandleFunc("/api/product", getProducts).Methods("GET")
	myRouter.HandleFunc("/api/product/{id}", getProduct).Methods("GET")
	myRouter.HandleFunc("/api/product/{id}", updateProduct).Methods("PUT")
	myRouter.HandleFunc("/api/product/{id}", deleteProduct).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":9000", myRouter))

	// myRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusNotFound)

	// 	res := Result{Code: 404, Message: "Method not found"}
	// 	response, _ := json.Marshal(res)
	// 	w.Write(response)
	// })
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello Product")
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body) // untuk menangkap data dari POST/ JSON (contoh ini kirim dari postman)

	var product Product
	json.Unmarshal(payloads, &product)

	db.Create(&product) //create data ke tabel product

	res := Response{Code: 200, Data: product, Message: "Success create product"}
	Response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(Response)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: get products")

	products := []Product{}
	db.Find(&products)

	res := Response{Code: 200, Data: products, Message: "Success get products"}
	Response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(Response)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	var product Product

	db.First(&product, productID)

	res := Response{Code: 200, Data: product, Message: "Success get detail product"}
	Response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(Response)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	payloads, _ := ioutil.ReadAll(r.Body)

	var productUpdates Product
	json.Unmarshal(payloads, &productUpdates)

	var product Product
	db.First(&product, productID)
	db.Model(&product).Updates(productUpdates)

	res := Response{Code: 200, Data: product, Message: "Success update product"}
	Response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(Response)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	var product Product

	db.First(&product, productID)
	db.Delete(&product)

	res := Response{Code: 200, Message: "Success delete product"}
	Response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(Response)
}
