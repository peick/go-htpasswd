package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/peick/go-htpasswd"
)

func main() {
	var handler http.Handler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("access granted\n"))
		},
	)

	users, err := htpasswd.New("examples/reload_sighup/.htpasswd")
	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)
	go func() {
		for {
			sig := <-sigs
			switch sig {
			case syscall.SIGHUP:
				log.Println("reloading users...")
				reloadErr := users.Reload()
				if reloadErr != nil {
					log.Fatal("failed to reload the htpasswd file", reloadErr)
				}
			}
		}
	}()

	basicAuthMiddleware := htpasswd.BasicAuthMiddleware("restricted", users)
	handler = basicAuthMiddleware(handler)

	log.Println("starting server")
	http.ListenAndServe(":8080", handler)
}
