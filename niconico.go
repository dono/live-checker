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

	client := new(http.Client)

	form := url.Values{}
	form.Add("mail_tel", mail)
	form.Add("password", pass)

	body := strings.NewReader(form.Encode())

	req, err := http.NewRequest("POST", URL, body)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; MAFSJS; rv:11.0) like Gecko")

	resp, err := client.Do(req)	
	if err != nil {
		log.Fatal(err)
	}

	// 302 ga kaette kuru hazu nanda kedo...

	fmt.Println(resp)
}