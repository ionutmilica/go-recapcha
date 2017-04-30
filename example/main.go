package main

import (
	"fmt"
	"github.com/ionutmilica/go-recaptcha"
	"html/template"
	"net/http"
)

const (
	secretKey = "--SECRET-KEY--"
	siteKey   = "--SITE-KEY--"
)

// Naive example to demonstrate the usage of the library
func main() {

	re := recaptcha.ReCaptcha{
		SecretKey: secretKey,
	}

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		// In production you should parse the templates outside of the http handlers
		t, err := template.ParseFiles("app.tmpl")
		if err != nil {
			fmt.Println(err)
			return
		}

		t.ExecuteTemplate(res, "app", map[string]string{
			"siteKey": siteKey,
		})
	})

	http.HandleFunc("/captcha", func(res http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			fmt.Fprintf(res, "Error: %s\n", err)
			return
		}

		response, err := re.Verify(req.FormValue("g-recaptcha-response"), req.RemoteAddr)
		if err != nil {
			fmt.Fprintf(res, "Http Error: %s\n", err)
			return
		}

		fmt.Fprintf(res, "Is captcha valid? %t\n", response.Success)
	})

	http.ListenAndServe(":8080", nil)
}
