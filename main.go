package main

import (
	"database/sql"
	"fmt"
	"github.com/bieber/barcode"
	"github.com/carlogit/phash"
	_ "github.com/mattn/go-sqlite3"
	"image/png"
	"os"
	"strings"
	"time"
)

func main1() {

	cabMng := &CabinetManager{
		portName: "/dev/ttyUSB0",
		frames:   make(chan []byte),
	}

	for i := 0; i < 3; i++ {
		cabMng.AddCabnit(i, 30)
	}
	cabMng.OpenPort()

	select {}
}

func main2() {

	client := NewClientSSH()

	err := client.Open("root", "root")
	if err != nil {
		return
	}

	client.KeyPress("esc")
	client.KeyPress("esc")
	client.KeyPress("esc")

	client.KeyPress("down")
	client.KeyPress("1")

	time.Sleep(time.Second * 5)
	client.Run("fbgrab /tmp/test.png")
	client.FileTransfer("/tmp/test.png")
	client.KeyPress("enter")
	time.Sleep(time.Second * 3)
	client.Run("fbgrab /tmp/test1.png")
	client.FileTransfer("/tmp/test1.png")

	file1, err := os.Open("/tmp/test.png")
	if err != nil {
		return
	}
	defer file1.Close()

	file2, _ := os.Open("/tmp/test1.png")
	defer file2.Close()
	fh1, err2 := phash.GetHash(file1)
	if err2 != nil {
		return
	}

	fh2, _ := phash.GetHash(file2)

	dist := phash.GetDistance(fh1, fh2)

	fmt.Printf("dist:%d", dist)

}

func main4() {
	fin, _ := os.Open("/home/feynman/Desktop/gotest/test.png")
	defer fin.Close()
	src, _ := png.Decode(fin)

	img := barcode.NewImage(src)
	scanner := barcode.NewScanner().SetEnabledAll(true)

	symbols, _ := scanner.ScanImage(img)
	for _, s := range symbols {
		fmt.Println(s.Type.Name(), s.Data, s.Quality, s.Boundary)
	}
}

func main5() {

	if err := postmanLogin("18676681017", "123456"); err != nil {
		return
	}

	for {
		if js, err := postmanBookBox("518067A362", "pkg1234568", "small", "13760268661", "18676681017"); err == nil {

			code := js.Get("code").MustInt()
			if code == 0 {

				fmt.Println(js)
				time.Sleep(time.Second * 1)

				parcelID := js.Get("data").Get("parcelId").MustString()
				if js, err := httpCancalBook("518067A362", parcelID); err == nil {

					fmt.Println(js)
					time.Sleep(time.Second * 5)
				}

			} else {

				fmt.Println(js)
				time.Sleep(time.Second)

			}

		}
	}

}

func main6() {

	db, err := sql.Open("sqlite3", "/home/feynman/Desktop/gotest/terminal_forbussiness.db")
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(`select [package_id],[book_code],[take_code],[retrieve_code] ,[postman_mobile],[take_mobile],[parcel_status] from t_parcel where parcel_id=?`)

	if err != nil {
		return
	}

	var pckID string
	var bookCode string

	var takenCode string
	var retriveCode string
	var postmanMobile string
	var takenMobile string
	var parcelStatus string

	err = stmt.QueryRow("518067A36220170911120800699322").Scan(&pckID, &bookCode, &takenCode, &retriveCode, &postmanMobile, &takenMobile, &parcelStatus)

	fmt.Println(pckID)
}

func main7() {

	cabMng := GetCabMng()
	cabMng.portName = "/dev/ttyUSB0"

	cabMng.AddCabnit(0, 32)
	cabMng.AddCabnit(1, 32)
	cabMng.AddCabnit(2, 32)

	cabMng.OpenPort()

	scan := GetScanner()
	scan.OpenPort("/dev/ttyUSB1")

	client := NewClientSSH()

	if err := client.Open("root", "root"); err != nil {
		return
	}

	//step 1 order a box
	if err := postmanLogin("18676681017", "123456"); err != nil {
		return
	}

	for {

		//step 1 order a box
		timeNow := time.Now()
		orderPckID := fmt.Sprintf("EZ%04d%02d%02d%02d%02d%02d", timeNow.Year(), timeNow.Month(), timeNow.Day(), timeNow.Hour(), timeNow.Minute(), timeNow.Second())
		if js, err := postmanBookBox("518067A362", orderPckID, "small", "13760268661", "18676681017"); err == nil {

			code := js.Get("code").MustInt()
			if code == 0 {

				fmt.Println(js)
				time.Sleep(time.Second * 5) //等5秒，让数据下发到终端

				//parcelID := js.Get("data").Get("parcelId").MustString()

				takenCode := js.Get("data").Get("takeCode").MustString()
				client.FileTransfer("/home/root/terminal_forbussiness.db")

				checkOrder(js, orderPckID)

				client.KeyPress("esc")
				client.KeyPress("esc")
				client.KeyPress("esc")
				client.KeyPress("esc")

				time.Sleep(2 * time.Second)
				client.Run("fbgrab /tmp/HomePage.png")
				client.FileTransfer("/tmp/HomePage.png")

				client.KeyPress("down")
				client.KeyPress("1")

				time.Sleep(2 * time.Second)
				client.Run("fbgrab /tmp/PostmanHomePage.png")
				client.FileTransfer("/tmp/PostmanHomePage.png")
				client.KeyPress("down")
				client.KeyPress("enter")

				client.Run("fbgrab /tmp/InputOrderPkgId.png")
				client.FileTransfer("/tmp/InputOrderPkgId.png")

				//输入预约单号

				scan.SendCode(orderPckID)
				// for i := 0; i < len(orderPckID); i++ {

				// 	client.KeyPress(orderPckID[i : i+1])
				// }

				//client.KeyPress("enter")
				time.Sleep(3 * time.Second)

				//open box OK

				//time.Sleep(3 * time.Second)

				client.Run("fbgrab /tmp/DeliveryOpenOK.png")
				client.FileTransfer("/tmp/DeliveryOpenOK.png")

				cabMng.SetAllBox(0)

				client.KeyPress("enter")

				time.Sleep(2 * time.Second)

				client.KeyPress("1")

				//scanner.SendCode(takenCode)

				for i := 0; i < len(takenCode); i++ {
					client.KeyPress(takenCode[i : i+1])
				}

				time.Sleep(3 * time.Second)

				client.KeyPress("enter")

			} else {

				fmt.Println(js)
				time.Sleep(time.Second)

			}

		}
	}
}

func main() {

	cabMng := GetCabMng()
	cabMng.portName = "/dev/ttyUSB0"

	cabMng.AddCabnit(0, 32)
	cabMng.AddCabnit(1, 32)
	cabMng.AddCabnit(2, 32)

	cabMng.OpenPort()

	scan := GetScanner()
	scan.OpenPort("/dev/ttyUSB1")

	adminLogin("13410324304", "123456")

	client := NewClientSSH()

	if err := client.Open("root", "root"); err != nil {
		return
	}

	for {

		client.KeyPress("esc")
		client.KeyPress("esc")
		client.KeyPress("esc")
		client.KeyPress("esc")
		client.KeyPress("1")
		//
		client.KeyPress("1")
		client.KeyPress("1")
		client.KeyPress("1")
		client.KeyPress("1")
		client.KeyPress("1")
		client.KeyPress("1")
		client.KeyPress("shift")

		time.Sleep(time.Second)

		client.Run("fbgrab /tmp/AdminLogin.png")

		client.FileTransfer("/tmp/AdminLogin.png")

		fin, _ := os.Open("/tmp/AdminLogin.png")

		src, _ := png.Decode(fin)

		fin.Close()

		img := barcode.NewImage(src)
		scanner := barcode.NewScanner().SetEnabledAll(true)

		symbols, _ := scanner.ScanImage(img)

		var qrCode string
		for _, s := range symbols {

			qrCode = s.Data
			fmt.Println(qrCode)
			break
			//fmt.Println(s.Type.Name(), s.Data, s.Quality, s.Boundary)

		}

		if qrCode != "" {

			index := strings.Index(qrCode, "?code=")

			if js, err := adminTerminalScanIn(qrCode[index+6:]); err == nil {

				fmt.Print(js)
			}
		}

		time.Sleep(time.Second * 5)
	}

}
