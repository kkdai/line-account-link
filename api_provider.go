package main

import (
	b64 "encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// CustData : Customers data for provider website.
type CustData struct {
	ID    string
	PW    string
	Name  string
	Age   int
	Desc  string
	Nonce string
}

var customers []CustData

func init() {
	//Init customer data in memory
	customers = append(customers, []CustData{
		CustData{ID: "11", PW: "pw11", Name: "Tom", Age: 43, Desc: "He is from A corp. likes to read comic books."},
		CustData{ID: "22", PW: "pw22", Name: "John", Age: 25, Desc: "He is from B corp. likes to read news paper"},
		CustData{ID: "33", PW: "pw33", Name: "Mary", Age: 13, Desc: "She is a student, like to read science books"},
	}...)
}

// WEB: List all user in memory
func listCust(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bookstore customer list as follow:\n")
	for i, usr := range customers {
		fmt.Fprintf(w, "%d \tID: %s \tName: %s \tPW: %s \tDesc:%s \n", i, usr.ID, usr.Name, usr.PW, usr.Desc)
	}
}

// WEB: For login (just for demo)
func login(w http.ResponseWriter, r *http.Request) {
	//7. The user enters his/her credentials.
	if err := r.ParseForm(); err != nil {
		log.Printf("ParseForm() err: %v\n", err)
		return
	}
	name := r.FormValue("user")
	pw := r.FormValue("pass")
	token := r.FormValue("token")
	for i, usr := range customers {
		if usr.ID == name {
			if pw == usr.PW {
				//8. The web server acquires the user ID from the provider's service and uses that to generate a nonce.
				sNonce := generateNonce(token, name, pw)

				//update nonce to provider DB to store it.
				customers[i].Nonce = sNonce

				//9. The web server redirects the user to the account-linking endpoint.
				//10. The user accesses the account-linking endpoint.
				//Print link to user to click it.
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
	fmt.Fprintf(w, "Your input name or password error.")
}

// WEB: For account link
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

// generate nonce (currently nonce combine by token + name + pw)
func generateNonce(token, name, pw string) string {
	return b64.StdEncoding.EncodeToString([]byte(token + name + pw))
}
