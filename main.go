// dto-skeleton-broker
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/cloudfoundry-community/types-cf"
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
	catalog := cf.Catalog{
		Services: []*cf.Service{{
			Name:        "sample-service",
			ID:          "c7067f66-3b6e-417e-bf8e-8ae317ddaafd", // https://www.guidgenerator.com/online-guid-generator.aspx
			Description: "sample-service",
			Plans: []*cf.Plan{
				{
					ID:          "9e2d6f97-c9d9-4924-820b-593e3744ed29",
					Name:        "sample-plan",
					Description: "sample-service-plan",
					Free:        true,
				},
			},
		}},
	}
	// cf requires that the body is always JSON, even if it is empty.
	c.JSON(200, catalog)
}

func createServiceInstance(c *gin.Context) {
	serviceID := c.Param("service_id")
	fmt.Printf("Creating service instance %s for service %s plan %s\n", serviceID)

	appEnv, err := cfenv.Current()
	if err != nil {
		c.AbortWithError(504, err)
		return
	}

	type serviceInstanceResponse struct {
		DashboardURL string `json:"dashboard_url"`
	}

	instance := serviceInstanceResponse{DashboardURL: fmt.Sprintf("https://%s/dashboard", appEnv.ApplicationURIs[0])}
	c.JSON(201, instance)
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
