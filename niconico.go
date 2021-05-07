package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)


const URL = `https://account.nicovideo.jp/login/redirector`

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	mail := os.Getenv("NICONICO_MAIL")
	pass := os.Getenv("NICONICO_PASS")

	client := &http.Client{
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return http.ErrUseLastResponse
        },
    }

	form := url.Values{}
	form.Add("mail_tel", mail)
	form.Add("password", pass)

	body := strings.NewReader(form.Encode())

	req, err := http.NewRequest("POST", URL, body)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)	
	if err != nil {
		log.Fatal(err)
	}

	userSession := ""
	userSessionSecure := ""

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "user_session" && cookie.Value != "deleted" {
			userSession = cookie.Value
		}
		if cookie.Name == "user_session_secure" {
			userSessionSecure = cookie.Value
		}
	}

	fmt.Println(userSession)
	fmt.Println(userSessionSecure)
}