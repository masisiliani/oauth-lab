package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
)

const (
	//AuthorizeURL is the URL we'll send the user to first to get their authorization
	AuthorizeURL = `https://github.com/login/oauth/authorize`
	//TokenURL is the endpoint we'll request an access token from
	TokenURL = `https://github.com/login/oauth/access_token`
	//APIURLBase is the GitHub base URL for API requests
	APIURLBase = `https://api.github.com/`
	//RedirectURL for this script, used as the redirect URL
	RedirectURL = `http://localhost:8080/login/github/authorization/callback`
	//AdminURL for this script, used as the Admin URL
	AdminURL = `http://localhost:8080/admin`
	//PortAddress is where the server is listening
	PortAddress = "8080"
)

type Configuration struct {
	Envs EnvironmentsVariables
}

type EnvironmentsVariables struct {
	GithubClientID     string
	GithubClientSecret string
}

var Config Configuration

var sessionManager *scs.SessionManager

///{"access_token":"e72e16c7e42f292c6912e7710c838347ae178b4a", "scope":"repo,gist", "token_type":"bearer"}
type AuthenticationToken struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func main() {

	initConfiguration()
	// Initialize a new session manager and configure the session lifetime.
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/login/github/authorization", authorizationHandler)
	mux.HandleFunc("/login/github/authorization/callback", authorizationCallbackHandler)
	mux.HandleFunc("/admin", adminHandler)

	fmt.Printf("Starting server at port %s...", PortAddress)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", PortAddress), sessionManager.LoadAndSave(mux)); err != nil {
		log.Fatal(err)
	}
}

func initConfiguration() {
	Config.Envs.GithubClientID = os.Getenv("GITHUB_CLIENT_ID")
	Config.Envs.GithubClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
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
	fmt.Fprintf(w, `<a href="/login/github/authorization">LOGIN</a>`)
	// htmlBody, err := ioutil.ReadFile("login.html")
	// if err != nil {
	// 	log.Print(err)
	// }

	// fmt.Fprintf(w, string(htmlBody))
}

func authorizationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := sessionManager.Clear(ctx)
	if err != nil {
		log.Print(err)
		return
	}

	uToken := uuid.New().String()

	githubURL, err := url.Parse(AuthorizeURL)
	if err != nil {
		log.Print(err)
		return
	}

	parameters := githubURL.Query()
	parameters.Add("client_id", Config.Envs.GithubClientID)
	parameters.Add("redirect_uri", RedirectURL)
	parameters.Add("scope", "public_repo")
	parameters.Add("state", uToken)

	githubURL.RawQuery = parameters.Encode()

	http.Redirect(w, r, githubURL.String(), http.StatusFound)
}

func authorizationCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	code := r.URL.Query().Get("code")
	uToken := r.URL.Query().Get("state")

	requestBodyMap := map[string]string{
		"client_id":     Config.Envs.GithubClientID,
		"client_secret": Config.Envs.GithubClientSecret,
		"code":          code,
		"state":         fmt.Sprintf("%v", uToken),
	}

	requestJSON, _ := json.Marshal(requestBodyMap)

	req, reqerr := http.NewRequest("POST", TokenURL, bytes.NewBuffer(requestJSON))
	if reqerr != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	responseAuthToken, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic("Request failed")
	}

	responseBytes, err := ioutil.ReadAll(responseAuthToken.Body)
	if err != nil {
		log.Print(err)
		return
	}

	authToken := AuthenticationToken{}
	json.Unmarshal(responseBytes, &authToken)
	log.Println(string(responseBytes))

	sessionManager.Put(r.Context(), "access_token", authToken.AccessToken)
	sessionManager.Put(r.Context(), "state", uToken)

	http.Redirect(w, r, AdminURL, http.StatusFound)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	if sessionManager.Exists(r.Context(), "access_token") {
		fmt.Fprint(w, "FUNCIONOU CARAIA!")
	}
}

func apiRequest(r *http.Request, urlGit string) {
	req, err := http.NewRequest("POST", urlGit, nil)
	if err != nil {
		log.Fatalln(err)
		return
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
