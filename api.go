// dto-skeleton-broker
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/cloudfoundry-community/types-cf"
	"github.com/gin-gonic/gin"
)

// A Broker represents a Cloud Foundry Service Broker
type Broker struct {
	env *cfenv.App
}

// API implements the Service Broker REST API
type API struct {
	*Broker
}

func (a *API) Catalog(c *gin.Context) {
	catalog := cf.Catalog{
		Services: []*cf.Service{{
			Name:        "sample-service",
			ID:          "c7067f66-3b6e-417e-bf8e-8ae317ddaafd", // https://www.guidgenerator.com/online-guid-generator.aspx
			Description: "sample-service",
			Bindable:    true,
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

func (a *API) createServiceInstance(c *gin.Context) {
	serviceID := c.Param("service_id")
	fmt.Printf("Creating service instance %s for service %s plan %s\n", serviceID)

	type serviceInstanceResponse struct {
		DashboardURL string `json:"dashboard_url"`
	}

	instance := serviceInstanceResponse{DashboardURL: fmt.Sprintf("https://%s/dashboard", a.Broker.env.ApplicationURIs[0])}
	c.JSON(201, instance)
}

func (a *API) deleteServiceInstance(c *gin.Context) {
	serviceID := c.Param("service_id")
	fmt.Printf("Deleting service instance %s for service %s plan %s\n", serviceID)
	c.JSON(200, struct{}{})
}

func (a *API) createServiceBinding(c *gin.Context) {

	type serviceBindingResponse struct {
		Credentials    map[string]interface{} `json:"credentials"`
		SyslogDrainURL string                 `json:"syslog_drain_url,omitempty"`
	}

	serviceID := c.Param("service_id")
	serviceBindingID := c.Param("binding_id")
	fmt.Printf("Creating service binding %s for service %s plan %s instance %s\n",
		serviceBindingID, serviceID)

	serviceBinding := serviceBindingResponse{
		SyslogDrainURL: os.Getenv("SYSLOG_DRAIN_URL"),
	}
	c.JSON(201, serviceBinding)
}

func (a *API) deleteServiceBinding(c *gin.Context) {
	serviceID := c.Param("service_id")
	serviceBindingID := c.Param("binding_id")
	fmt.Printf("Delete service binding %s for service %s plan %s instance %s\n",
		serviceBindingID, serviceID)
	c.JSON(200, struct{}{})
}

// newBrokerAPI returns a http.Handler which exposes the Cloud Foundry
// Service Broker API for the supplied Broker implementation.
// The broker is always protected by the user and pass basic auth
// credentials.
func newBrokerAPI(b *Broker, user, pass string) http.Handler {
	if user == "" || pass == "" {
		log.Fatal("AUTH_USER and AUTH_PASS must be set")
	}

	// create new router
	g := gin.Default()

	// apply basic auth over all endpoints
	authorized := g.Group("/", gin.BasicAuth(gin.Accounts{user: pass}))

	api := API{
		Broker: b,
	}

	authorized.GET("/v2/catalog", api.Catalog)
	authorized.PUT("/v2/service_instances/:service_id", api.createServiceInstance)
	authorized.DELETE("/v2/service_instances/:service_id", api.deleteServiceInstance)
	authorized.PUT("/v2/service_instances/:service_id/service_bindings/:binding_id", api.createServiceBinding)
	authorized.DELETE("/v2/service_instances/:service_id/service_bindings/:binding_id", api.deleteServiceBinding)

	return g
}
