package main

/*
 *https://utahta.hatenablog.com/blog/2016/08/16/Google_OAuth2_%E3%83%88%E3%83%BC%E3%82%AF%E3%83%B3%E3%82%92%E5%8F%96%E5%BE%97%E3%81%99%E3%82%8B_with_Go
 */

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauthapi "google.golang.org/api/oauth2/v2"
)

func main() {
	if err := godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV"))); err != nil {
		return
	}
	confGoogle.ClientID = os.Getenv("GOOGLE_CLIENT_ID")
	confGoogle.ClientSecret = os.Getenv("GOOGLE_SECRET")
	confIntra.ClientID = os.Getenv("INTRA_CLIENT_ID")
	confIntra.ClientSecret = os.Getenv("INTRA_SECRET")

	mux := http.NewServeMux()
	mux.HandleFunc("/login/google", GoogleLoginHandler)
	mux.HandleFunc("/login/google/redirect", GoogleLoginRHandler)
	mux.HandleFunc("/login/intra", IntraLoginHandler)
	mux.HandleFunc("/login/intra/redirect", IntraLoginRHandler)

	log.Println("Server has started")
	fmt.Println("Pleas access: http://localhost:5001/login/google")
	fmt.Println("Pleas access: http://localhost:5001/login/intra")
	http.ListenAndServe(":5001", mux)
}

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	var url = confGoogle.AuthCodeURL("yourStateUUID", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleLoginRHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query()["code"]
	if code == nil || len(code) == 0 {
		fmt.Fprint(w, "Invalid Parameter")
	}
	ctx := context.Background()
	tok, err := confGoogle.Exchange(ctx, code[0])
	if err != nil {
		fmt.Fprintf(w, "OAuth Error:%v", err)
	}
	client := confGoogle.Client(ctx, tok)
	svr, err := oauthapi.New(client)
	ui, err := svr.Userinfo.Get().Do()
	if err != nil {
		fmt.Fprintf(w, "OAuth Error:%v", err)
	} else {
		fmt.Fprintf(w, "Your are logined as : %s", ui.Email)
	}
}

var confGoogle = &oauth2.Config{
	ClientID:     "",
	ClientSecret: "",
	Scopes:       []string{oauthapi.UserinfoEmailScope},
	Endpoint:     google.Endpoint,
	RedirectURL:  "http://localhost:5001/login/google/redirect",
}

func IntraLoginHandler(w http.ResponseWriter, r *http.Request) {
	var url = confIntra.AuthCodeURL("")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func IntraLoginRHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("IntraLoginRHandler")
	code := r.URL.Query()["code"]
	if code == nil || len(code) == 0 {
		fmt.Fprint(w, "Invalid Parameter")
	}
	log.Println("Code is valid")
	ctx := context.Background()
	tok, err := confIntra.Exchange(ctx, code[0])
	if err != nil {
		fmt.Fprintf(w, "OAuth Error:%v", err)
	}
	log.Println("Token exchange success")
	client := confIntra.Client(ctx, tok)
	res, err := client.Get("https://api.intra.42.fr/v2/me/projects")
	if err != nil {
		log.Println("/me/projects failed")
		fmt.Fprintln(w, "Error: ", err)
	} else {
		log.Println("/me/projects SUCCEEDED!!!!!!!!")
		fmt.Fprintln(w, res.Body)
	}
}

var confIntra = &oauth2.Config{
	ClientID:     "",
	ClientSecret: "",
	Scopes:       []string{"public", "projects", "profile", "elearning", "tig", "forum"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://api.intra.42.fr/oauth/authorize",
		TokenURL: "https://api.intra.42.fr/oauth/token",
	},
	RedirectURL: "http://localhost:5001/login/intra/redirect",
}
