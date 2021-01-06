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
)

var db *gorm.DB
var err error

type Employee struct {
	ID                int    `json:"id"`
	Nip               int    `json:"nip"`
	Nama              string `json:"nama"`
	Tgllahir          string `json:"tgllahir"`
	Jeniskelamin      string `json:"jeniskelamin"`
	Agamaid           int    `json:"agamaid"`
	Telfon            string `json:"telfon"`
	Bagianid          int    `json:"bagianid"`
	Statuskepegawaian string `json:"statuskepegawaian"`
	Keterangan        string `json:"keterangan"`
}

type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	db, err = gorm.Open("mysql", "root:@/db_sim?charset=utf8&parseTime=True")

	if err != nil {
		log.Println("Connection failed", err)
	} else {
		log.Println("Connection established")
	}

	db.AutoMigrate(&Employee{})
	handleRequests()
}

func handleRequests() {
	log.Println("Started Server at http://127.0.0.1:9000")

	route := mux.NewRouter().StrictSlash(true) //fungsi untuk membuat router baru

	route.HandleFunc("/", Homepage)
	route.HandleFunc("/api/employee", CreateEmployee).Methods("POST")
	route.HandleFunc("/api/employee", GetEmployees).Methods("GET")
	route.HandleFunc("/api/employee/{id}", GetEmployee).Methods("GET")
	route.HandleFunc("/api/employee/{id}", UpdateEmployee).Methods("PUT")
	route.HandleFunc("/api/employee/{id}", DeleteEmployee).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":9000", route))

}

func Homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello website sim")
}

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body) //untuk menangkap data dari POST/ JSON (contoh ini kirim dari postman)

	var employee Employee
	json.Unmarshal(payloads, &employee)

	db.Create(&employee) //create data ke tabel employee

	res := Response{Data: employee, Message: "Success create employee"}
	Response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(Response)
}

func GetEmployees(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: get employee")

	employee := []Employee{}
	db.Find(&employee)

	res := Response{Data: employee, Message: "Success get employee"}
	Response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(Response)
}

func GetEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeID := vars["id"]

	var employee Employee

	db.First(&employee, employeeID)

	res := Response{Data: employee, Message: "Success get detail employee"}
	Response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(Response)
}

func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeID := vars["id"]

	payloads, _ := ioutil.ReadAll(r.Body)

	var employeeUpdates Employee
	json.Unmarshal(payloads, &employeeUpdates)

	var employee Employee
	db.First(&employee, employeeID)
	db.Model(&employee).Updates(employeeUpdates)

	res := Response{Data: employee, Message: "Success update employee"}
	Response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(Response)
}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeID := vars["id"]

	var employee Employee

	db.First(&employee, employeeID)
	db.Delete(&employee)

	res := Response{Message: "Success delete employee"}
	Response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(Response)
}
