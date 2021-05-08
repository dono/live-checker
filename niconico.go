package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
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

	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		log.Fatal(err)
	}

	mail := os.Getenv("NICONICO_MAIL")
	pass := os.Getenv("NICONICO_PASS")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		 Jar:     jar,
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

	_, err = client.Do(req)	
	if err != nil {
		log.Fatal(err)
	}

	req, err = http.NewRequest("GET", "https://com.nicovideo.jp/community/co3000390", body)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(bytes))

	

	// userSession := ""
	// userSessionSecure := ""

	// for _, cookie := range resp.Cookies() {
	// 	if cookie.Name == "user_session" && cookie.Value != "deleted" {
	// 		userSession = cookie.Value
	// 	}
	// 	if cookie.Name == "user_session_secure" {
	// 		userSessionSecure = cookie.Value
	// 	}
	// }

	// fmt.Println(userSession)
	// fmt.Println(userSessionSecure)
}