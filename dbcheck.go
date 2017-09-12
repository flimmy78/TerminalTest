package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
)

func checkOrder(js *simplejson.Json, orderPackID string) error {

	parcelID := js.Get("data").Get("parcelId").MustString()
	// orderBookCode := js.Get("data").Get("bookCode").MustString()
	// orderRetrieveCode := js.Get("data").Get("retrieveCode").MustString()
	// orderTakenCode := js.Get("data").Get("takeCode").MustString()

	db, err := sql.Open("sqlite3", "/tmp/terminal_forbussiness.db")
	if err != nil {
		return errors.New("")
	}
	defer db.Close()

	stmt, err := db.Prepare(`select [package_id],[book_code],[take_code],[retrieve_code] ,[postman_mobile],[take_mobile],[parcel_status] from t_parcel where parcel_id=?`)

	if err != nil {
		return errors.New("")
	}

	var pckID string
	var bookCode string

	var takenCode string
	var retriveCode string
	var postmanMobile string
	var takenMobile string
	var parcelStatus string

	err = stmt.QueryRow(parcelID).Scan(&pckID, &bookCode, &takenCode, &retriveCode, &postmanMobile, &takenMobile, &parcelStatus)

	fmt.Println(pckID)

	// md5TakenCode := fmt.Sprintf("%x", md5.Sum([]byte(orderTakenCode)))

	// md5BookCode := fmt.Sprintf("%x", md5.Sum([]byte(orderBookCode)))
	// md5RetrieveCode := fmt.Sprintf("%x", md5.Sum([]byte(orderRetrieveCode)))
	if orderPackID != pckID {
		return errors.New("")
	}

	// if md5RetrieveCode != retriveCode {
	// 	return errors.New("")
	// }

	// if md5BookCode != bookCode {
	// 	return errors.New("")
	// }

	// if md5TakenCode != takenCode {
	// 	return errors.New("")
	// }

	if parcelStatus != "book" {
		return errors.New("")
	}

	return nil
}

func checkOrderCancel(js *simplejson.Json) error {

	return nil
}
