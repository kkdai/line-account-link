package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

//CustData : Customers data for provider website.
type CustData struct {
	ID   string
	PW   string
	Desc string
}

var customers []CustData

func init() {
	//Init customer data in memory
	customers = append(customers, []CustData{
		CustData{ID: "11", PW: "pw11", Desc: "This is 11"},
		CustData{ID: "22", PW: "pw22", Desc: "This is 22"},
		CustData{ID: "33", PW: "pw33", Desc: "This is 33"},
	}...)
}

//WEB: List all user in memory
func listCust(w http.ResponseWriter, r *http.Request) {
	for i, usr := range customers {
		fmt.Fprintf(w, "%d \t ID: %s \t Desc:%s \n", i, usr.ID, usr.Desc)
	}
}

//WEB: For login (just for demo)
func login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("ParseForm() err: %v\n", err)
		return
	}
	name := r.FormValue("user")
	pass := r.FormValue("pass")
	token := r.FormValue("token")
	for _, usr := range customers {
		log.Println("usr:=", usr, " name=", name, " pass=", pass)
		if usr.ID == name {
			if pass == usr.PW {
				//generate nonce
				log.Println("generate nonce")
				http.Redirect(w, r, fmt.Sprintf("https://access.line.me/dialog/bot/accountLink?linkToken=%s&nonce=1234", token), 200)
				return
			}
		}
	}
	fmt.Fprintf(w, "Your input name or password error.")
}

//WEB: For account link
func link(w http.ResponseWriter, r *http.Request) {
	TOKEN := r.FormValue("linkToken")
	if TOKEN == "" {
		log.Println("No token.")
		return
	}

	log.Println("token = ", TOKEN)
	tmpl := template.Must(template.ParseFiles("login.tmpl"))
	if err := tmpl.Execute(w, TOKEN); err != nil {
		log.Println("Template err:", err)
	}
}
