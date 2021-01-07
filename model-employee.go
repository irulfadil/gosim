package main

import (
	"database/sql"
	"fmt"
)

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

func getEmployees(db *sql.DB, start, count int) ([]Employee, error) {
	query := fmt.Sprintf("SELECT id, nip, nama, tgllahir, jeniskelamin, agamaid, telfon, bagianid, statuskepegawaian, keterangan FROM pegawai LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	Employees := []Employee{}

	for rows.Next() {
		var e Employee
		if err := rows.Scan(&e.ID, &e.Nip, &e.Nama, &e.Tgllahir, &e.Jeniskelamin, &e.Agamaid, &e.Telfon, &e.Bagianid, &e.Statuskepegawaian, &e.Keterangan); err != nil {
			return nil, err
		}
		Employees = append(Employees, e)
	}
	return Employees, nil

}

func (e *Employee) getEmployee(db *sql.DB) error {
	query := fmt.Sprintf("SELECT FROM pegawai WHERE %d", e.ID)
	return db.QueryRow(query).Scan(&e.ID, &e.Nip, &e.Nama, &e.Tgllahir, &e.Jeniskelamin, &e.Agamaid, &e.Telfon, &e.Bagianid, &e.Statuskepegawaian)
}

func (e *Employee) createEmployee(db *sql.DB) error {
	query := fmt.Sprintf("INSERT INTO pegawai (nip, nama, tgllahir, jeniskelamin, agamaid,telfon, bagianid, statuskepegawaian) VALUES()", &e.ID, &e.Nip, &e.Nama, &e.Tgllahir, &e.Jeniskelamin, &e.Agamaid, &e.Telfon, &e.Bagianid, &e.Statuskepegawaian)
	_, err := db.Exec(query)

	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&e.ID)

	if err != nil {
		return err
	}
	return nil

}

func (e *Employee) updateEmployee(db *sql.DB) error {
	query := fmt.Sprintf("UPDATE pegawai SET WHERE id=%d", &e.ID)
	_, err := db.Exec(query)
	return err
}

func (e *Employee) deleteEmployee(db *sql.DB) error {
	query := fmt.Sprintf("DELETE FROM pegawai WHERE id=%d", &e.ID)
	_, err := db.Exec(query)
	return err
}
