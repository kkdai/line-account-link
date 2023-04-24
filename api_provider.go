package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// CustData : Customers data for provider website.
type CustData struct {
	ID     string
	PW     string
	Name   string
	Nounce string
	UserId string
}

// var db *sql.DB
var customers []CustData

func init() {
	//Init customer data in memory
}

// WEB: List all user in memory
func listCust(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bookstore customer list as follow:\n")
	for i, usr := range customers {
		fmt.Fprintf(w, "%d \tID: %s \tName: %s \n", i, usr.ID, usr.Name)
	}
}

// WEB: For login (just for demo)
func login(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", "canis:vz3s10cdDtkU1BRv@tcp(103.200.113.92)/foodler")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	results, err := db.Query("SELECT `username`, `password`, `fname` FROM `users` WHERE `identity` = 'customer'")
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()
	var user CustData
	for results.Next() {
		results.Scan(&user.ID, &user.PW, &user.Name)
		customers = append(customers, user)
	}
	//7. The user enters his/her credentials.
	if err := r.ParseForm(); err != nil {
		log.Printf("ParseForm() err: %v\n", err)
		return
	}

	name := r.FormValue("user")
	pw := r.FormValue("pass")
	token := r.FormValue("token")

	PW := MD5(pw)
	for _, usr := range customers {
		if usr.ID == name {
			if PW == usr.PW {
				//8. The web server acquires the user ID from the provider's service and uses that to generate a nonce.
				sNonce := generateNounce(token, name, pw)
				rs, err := db.Exec("UPDATE `linebot` SET `nounce`= ? WHERE `username` = ?", sNonce, usr.ID)
				if err != nil {
					log.Println("exec failed:", err)
					return
				}
				idAff, err := rs.RowsAffected()
				if err != nil {
					log.Println("RowsAffected failed:", err)
					return
				}
				if idAff == 0 {
					usr.UserId = ""
					log.Println("Nounce:" + usr.Nounce + "\nUsername:" + usr.ID + "\nPassword:" + usr.PW + "\nUserId:" + usr.UserId + "\nName:" + usr.Name)
					_, err := db.Exec("INSERT INTO `linebot`(`nounce`, `username`, `password`, `userId`, `name`) VALUES (?, ?, ?, ?, ?)", sNonce, usr.ID, usr.PW, usr.UserId, usr.Name)
					if err != nil {
						log.Println("exec failed", err)
					}
				}

				results, err := db.Query("SELECT `nounce` FROM linebot WHERE `username` = ?", usr.ID)
				if err != nil {
					panic(err.Error())
				}
				defer results.Close()
				for results.Next() {
					results.Scan(&user.Nounce)
					customers = append(customers, user)
				}

				targetURL := fmt.Sprintf("https://access.line.me/dialog/bot/accountLink?linkToken=%s&nonce=%s", token, sNonce)

				log.Println("generate nonce, targetURL=", targetURL)
				tmpl := template.Must(template.ParseFiles("link.tmpl"))
				if err := tmpl.Execute(w, targetURL); err != nil {
					log.Println("Template err:", err)
				}
				return
			}
		}
	}
	fmt.Fprintf(w, "您輸入的帳號有誤!")

}

func link(w http.ResponseWriter, r *http.Request) {
	//5. The user accesses the linking URL.
	TOKEN := r.FormValue("linkToken")
	if TOKEN == "" {
		log.Println("No token.")
		return
	}

	log.Println("token = ", TOKEN)
	tmpl := template.Must(template.ParseFiles("login.tmpl"))
	//6. The web server displays the login screen.
	if err := tmpl.Execute(w, TOKEN); err != nil {
		log.Println("Template err:", err)
	}
}

// generate nonce (currently nounce combine by token + name + pw)
func generateNounce(token, name, pw string) string {
	return base64.StdEncoding.EncodeToString([]byte(token + name + pw))
}

func MD5(pw string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(pw))
	return hex.EncodeToString(algorithm.Sum(nil))
}
