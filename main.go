package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const (
	//AuthorizeURL is the URL we'll send the user to first to get their authorization
	AuthorizeURL = "https://github.com/login/oauth/authorize"
	//TokenURL is the endpoint we'll request an access token from
	TokenURL = "https://github.com/login/oauth/access_token"
	//APIURLBase is the GitHub base URL for API requests
	APIURLBase = "https://api.github.com/"
	//BaseURL for this script, used as the redirect URL
	BaseURL = "https://localhost"
	//PortAddress is where the server is listening
	PortAddress = "8080"
	//
)

var sessionManager *scs.SessionManager

func main() {

	// Initialize a new session manager and configure the session lifetime.
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/admin", adminHandler)

	fmt.Printf("Starting server at port %s...", PortAddress)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", PortAddress), sessionManager.LoadAndSave(mux)); err != nil {
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

	if sessionManager.Exists(r.Context(), "access_token") == false {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "POST request successful \n\n")
	email := r.FormValue("inputEmail")
	password := r.FormValue("inputPassword")

	fmt.Fprintf(w, "Name: %s \nPassword: %s \n", email, password)
}

func apiRequest(r *http.Request, urlGit string) {

	req, err := http.NewRequest("POST", urlGit, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Add("Accept", "application/vnd.github.v3+json, application/json")
	req.Header.Add("User-Agent", fmt.Sprintf("http://localhost:%s", PortAddress))

	if sessionManager.Exists(r.Context(), "access_token") {
		req.Header.Add("Authorization", "Bearer")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(body)

}
