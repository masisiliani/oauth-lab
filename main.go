package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)

	fmt.Println("Starting server at port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "[TODO] implement redirect if user is not logged")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "[TODO] implement login page")
}
