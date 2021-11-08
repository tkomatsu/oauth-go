package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

type Auth struct {
	ClientID    string
	Secret      string
	RedirectURL string
	Scopes      []string
}

func GetAuthInfo() (*Auth, error) {
	if err := godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV"))); err != nil {
		return nil, err
	}
	auth := &Auth{
		ClientID:    os.Getenv("CLIENT_ID"),
		Secret:      os.Getenv("SECRET"),
		RedirectURL: "urn:ietf:wg:oauth:2.0:oob",
		Scopes:      []string{"https://www.googleapis.com/auth/drive"},
	}
	return auth, nil
}

func main() {
	auth, err := GetAuthInfo()
	if err != nil {
		return
	}
	fmt.Println("Start Execute API")

	config := &oauth2.Config{
		ClientID:     auth.ClientID,
		ClientSecret: auth.Secret,
		RedirectURL:  auth.RedirectURL,
		Scopes:       auth.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}

	url := config.AuthCodeURL("")
	fmt.Println("ブラウザで以下のURLにアクセスし、認証してください。")
	fmt.Println(url)
	fmt.Println("")

	fmt.Printf("Input auth code: ")
	var code string
	fmt.Scanf("%s\n", &code)

	token, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(err)
	}
	text, err := json.Marshal(token)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("token.json", text, 0777)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("token is saved to token.json")
}
