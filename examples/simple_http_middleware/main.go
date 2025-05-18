package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/peick/go-htpasswd"
)

func main() {
	var handler http.Handler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("access granted\n"))
		},
	)

	// the password is 'password'
	usersReader := strings.NewReader(`username:$apr1$VfoHyKyF$EQ3gDdg7EUQB69/ppHOOU0`)
	users, err := htpasswd.NewFromReader(usersReader)
	if err != nil {
		log.Fatal(err)
	}

	basicAuthMiddleware := htpasswd.BasicAuthMiddleware("restricted", users)
	handler = basicAuthMiddleware(handler)

	log.Println("starting server")
	http.ListenAndServe(":8080", handler)
}
