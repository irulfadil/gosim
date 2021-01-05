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

type Pegawai struct {
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
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	// db, err = gorm.Open("mysql", "root:@/db_sim")
	db, err = gorm.Open("mysql", "root:@/db_sim?charset=utf8&parseTime=True")

	if err != nil {
		log.Println("Connection failed", err)
	} else {
		log.Println("Connection established")
	}

	// db.AutoMigrate(&Pegawai{})
	handleRequests()
}

func handleRequests() {
	log.Println("Started Server at http://127.0.0.1:9000")

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homepage)
	myRouter.HandleFunc("/api/pegawai", createPegawai).Methods("POST")
	myRouter.HandleFunc("/api/pegawai", getPegawais).Methods("GET")
	myRouter.HandleFunc("/api/pegawai/{id}", getPegawai).Methods("GET")
	myRouter.HandleFunc("/api/pegawai/{id}", updatePegawai).Methods("PUT")
	myRouter.HandleFunc("/api/pegawai/{id}", deletePegawai).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":9000", myRouter))
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello website sim")
}

func createPegawai(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body) //untuk menangkap data dari POST/ JSON (contoh ini kirim dari postman)

	var pegawai Pegawai
	json.Unmarshal(payloads, &pegawai)

	db.Create(&pegawai) //create data ke tabel pegawai

	res := Response{Code: 200, Data: pegawai, Message: "Success create pegawai"}
	Response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(Response)
}

func getPegawais(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: get pegawai")

	pegawai := []Pegawai{}
	db.Find(&pegawai)

	res := Response{Code: 200, Data: pegawai, Message: "Success get pegawai"}
	Response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(Response)
}

func getPegawai(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pegawaiID := vars["id"]

	var pegawai Pegawai

	db.First(&pegawai, pegawaiID)

	res := Response{Code: 200, Data: pegawai, Message: "Success get detail pegawai"}
	Response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(Response)
}

func updatePegawai(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pegawaiID := vars["id"]

	payloads, _ := ioutil.ReadAll(r.Body)

	var pegawaiUpdates Pegawai
	json.Unmarshal(payloads, &pegawaiUpdates)

	var pegawai Pegawai
	db.First(&pegawai, pegawaiID)
	db.Model(&pegawai).Updates(pegawaiUpdates)

	res := Response{Code: 200, Data: pegawai, Message: "Success update pegawai"}
	Response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(Response)
}

func deletePegawai(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pegawaiID := vars["id"]

	var pegawai Pegawai

	db.First(&pegawai, pegawaiID)
	db.Delete(&pegawai)

	res := Response{Code: 200, Message: "Success delete pegawai"}
	Response, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(Response)
}
