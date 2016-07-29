package broker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/cloudfoundry-community/types-cf"
	"github.com/gin-gonic/gin"
)

// API implements the Service Broker REST API
type API struct {
	Env *cfenv.App
	Broker
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

func (a *API) Provision(c *gin.Context) {
	instanceid := c.Param("instance_id")

	type ProvisionDetails struct {
		ServiceID        string          `json:"service_id"`
		PlanID           string          `json:"plan_id"`
		OrganizationGUID string          `json:"organization_guid"`
		SpaceGUID        string          `json:"space_guid"`
		RawParameters    json.RawMessage `json:"parameters,omitempty"`
	}

	var details ProvisionDetails

	if err := json.NewDecoder(c.Request.Body).Decode(&details); err != nil {
		c.JSON(422, struct {
			Description string
		}{Description: err.Error()})
		return
	}

	if err := a.Broker.Provision(instanceid, details.ServiceID, details.PlanID); err != nil {
		c.JSON(504, struct {
			Description string
		}{Description: err.Error()})
		return
	}

	type serviceInstanceResponse struct {
		DashboardURL string `json:"dashboard_url"`
	}

	instance := serviceInstanceResponse{DashboardURL: fmt.Sprintf("https://%s/dashboard", a.Env.ApplicationURIs[0])}
	c.JSON(201, instance)
}

func (a *API) Deprovision(c *gin.Context) {
	instanceid := c.Param("instance_id")
	serviceid := c.Query("service_id")
	planid := c.Query("plan_id")

	if err := a.Broker.Deprovision(instanceid, serviceid, planid); err != nil {
		c.JSON(504, struct {
			Description string
		}{Description: err.Error()})
		return
	}
	c.JSON(200, struct{}{})
}

func (a *API) Bind(c *gin.Context) {

	type serviceBindingResponse struct {
		Credentials    map[string]interface{} `json:"credentials"`
		SyslogDrainURL string                 `json:"syslog_drain_url,omitempty"`
	}

	type BindResource struct {
		AppGuid string `json:"app_guid,omitempty"`
		Route   string `json:"route,omitempty"`
	}

	type BindDetails struct {
		AppGUID      string                 `json:"app_guid"`
		PlanID       string                 `json:"plan_id"`
		ServiceID    string                 `json:"service_id"`
		BindResource *BindResource          `json:"bind_resource,omitempty"`
		Parameters   map[string]interface{} `json:"parameters,omitempty"`
	}

	instanceid := c.Param("instance_id")
	bindingid := c.Param("binding_id")

	var details BindDetails
	if err := json.NewDecoder(c.Request.Body).Decode(&details); err != nil {
		c.JSON(422, struct {
			Description string
		}{Description: err.Error()})
		return
	}

	if err := a.Broker.Bind(instanceid, bindingid, details.ServiceID, details.PlanID); err != nil {
		c.JSON(504, struct {
			Description string
		}{Description: err.Error()})
		return
	}

	serviceBinding := serviceBindingResponse{
		Credentials: map[string]interface{}{
			"user":     "scott",
			"password": "tiger",
		},
		SyslogDrainURL: os.Getenv("SYSLOG_DRAIN_URL"),
	}
	c.JSON(201, serviceBinding)
}

func (a *API) Unbind(c *gin.Context) {
	instanceid := c.Param("instance_id")
	bindingid := c.Param("binding_id")
	serviceid := c.Query("service_id")
	planid := c.Query("plan_id")

	if err := a.Broker.Unbind(instanceid, bindingid, serviceid, planid); err != nil {
		c.JSON(504, struct {
			Description string
		}{Description: err.Error()})
		return
	}
	c.JSON(200, struct{}{})
}

// NewAPI returns a http.Handler which exposes the Cloud Foundry
// Service Broker API using the supplied Broker.
// The broker is always protected by the user and pass basic auth
// credentials.
func NewAPI(env *cfenv.App, b Broker, user, pass string) http.Handler {
	if user == "" || pass == "" {
		log.Fatal("AUTH_USER and AUTH_PASS must be set")
	}

	// create new router
	g := gin.Default()

	// apply basic auth over all endpoints
	authorized := g.Group("/", gin.BasicAuth(gin.Accounts{user: pass}))

	api := API{
		Env: env,
		Broker: &validatingBroker{
			Broker: b,
		},
	}

	authorized.GET("/v2/catalog", api.Catalog)
	authorized.PUT("/v2/service_instances/:instance_id", api.Provision)
	authorized.DELETE("/v2/service_instances/:instance_id", api.Deprovision)
	authorized.PUT("/v2/service_instances/:instance_id/service_bindings/:binding_id", api.Bind)
	authorized.DELETE("/v2/service_instances/:instance_id/service_bindings/:binding_id", api.Unbind)

	return g
}
