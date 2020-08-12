package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/admin", adminHandler)

	fmt.Println("Starting server at port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported", http.StatusNotFound)
	}

	fmt.Fprintf(w, "[TODO] implement redirect if user is not logged")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	htmlBody, err := ioutil.ReadFile("login.html")
	if err != nil {
		log.Print(err)
	}

	fmt.Fprintf(w, string(htmlBody))
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	fmt.Fprintf(w, "POST request successful \n\n")
	email := r.FormValue("inputEmail")
	password := r.FormValue("inputPassword")

	fmt.Fprintf(w, "Name: %s \nPassword: %s \n", email, password)

	fmt.Fprintf(w, "\n \n[TODO] Implements redirect if user is not logged")
}
