// dto-skeleton-broker
package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"

	"github.com/AusDTO/dto-skeleton-broker/internal/broker"
)

func envOr(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func main() {
	port := flag.String("p", envOr("PORT", "3000"), "port to listen")
	flag.Parse()

	addr := ":" + *port

	appEnv, err := cfenv.Current()
	if err != nil {
		log.Fatal(err)
	}

	b := new(MockBroker)
	api := broker.NewAPI(appEnv, b, os.Getenv("AUTH_USER"), os.Getenv("AUTH_PASS"))

	log.Println(os.Args[0], "listening on", addr)
	if err := http.ListenAndServe(addr, api); err != nil {
		log.Fatal(err)
	}
}
