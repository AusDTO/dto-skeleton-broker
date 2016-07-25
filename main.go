// dto-skeleton-broker
package main

import (
	"flag"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func envOr(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func brokerCatalog(c *gin.Context) {
	// cf requires that the body is always JSON, even if it is empty.
	c.JSON(200, struct{}{})
}

func createServiceInstance(c *gin.Context) {

}

func deleteServiceInstance(c *gin.Context) {

}

func createServiceBinding(c *gin.Context) {

}

func deleteServiceBinding(c *gin.Context) {

}

func main() {
	port := flag.String("p", envOr("PORT", "3000"), "port to listen")
	flag.Parse()

	addr := ":" + *port

	// create new router
	g := gin.Default()

	// apply basic auth over all endpoints
	user := os.Getenv("AUTH_USER")
	pass := os.Getenv("AUTH_PASS")
	if user == "" || pass == "" {
		log.Fatal("AUTH_USER and AUTH_PASS must be set")
	}
	authorized := g.Group("/", gin.BasicAuth(gin.Accounts{user: pass}))

	// Cloud Foundry Service API
	authorized.GET("/v2/catalog", brokerCatalog)
	authorized.PUT("/v2/service_instances/:service_id", createServiceInstance)
	authorized.DELETE("/v2/service_instances/:service_id", deleteServiceInstance)
	authorized.PUT("/v2/service_instances/:service_id/service_bindings/:binding_id", createServiceBinding)
	authorized.DELETE("/v2/service_instances/:service_id/service_bindings/:binding_id", deleteServiceBinding)

	log.Println(os.Args[0], "listening on", addr)
	err := g.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}
