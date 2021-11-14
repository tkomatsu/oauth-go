package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauthapi "google.golang.org/api/oauth2/v2"
)

var (
	confGoogle *oauth2.Config
	confIntra  *oauth2.Config
)

func main() {
	if err := setConfig(); err != nil {
		return
	}

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

func setConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Print("[ERROR] viper: ", err)
		return err
	}

	key := viper.GetStringMapString("google")
	confGoogle = &oauth2.Config{
		ClientID:     key["client_id"],
		ClientSecret: key["client_secret"],
		Scopes:       []string{oauthapi.UserinfoEmailScope},
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:5001/login/google/redirect",
	}

	key = viper.GetStringMapString("intra")
	confIntra = &oauth2.Config{
		ClientID:     key["client_id"],
		ClientSecret: key["client_secret"],
		Scopes:       []string{"public", "projects", "profile", "elearning", "tig", "forum"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.intra.42.fr/oauth/authorize",
			TokenURL: "https://api.intra.42.fr/oauth/token",
		},
		RedirectURL: "http://localhost:5001/login/intra/redirect",
	}

	return nil
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
		confmap := viper.GetStringMapString("google")
		confmap["access_token"] = tok.AccessToken
		viper.Set("google", confmap)
		viper.WriteConfig()
	}
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
	if err != nil || res.StatusCode != http.StatusOK {
		log.Println("/me/projects failed")
		fmt.Fprintln(w, "Error: ", err)
	} else {
		log.Println("/me/projects SUCCEEDED!!!!!!!!")
		fmt.Fprintln(w, res.Body)
		confmap := viper.GetStringMapString("intra")
		confmap["access_token"] = tok.AccessToken
		viper.Set("intra", confmap)
		viper.WriteConfig()
	}
}
