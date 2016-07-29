package broker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testBroker struct {
	*testing.T
}

func (t *testBroker) Provision(instanceid, serviceid, planid string) error {
	t.Logf("Creating service instance %s for service %s plan %s\n", instanceid, serviceid, planid)
	return nil
}

func (t *testBroker) Deprovision(instanceid, serviceid, planid string) error {
	t.Logf("Deleting service instance %s for service %s plan %s\n", instanceid, serviceid, planid)
	return nil
}

func (t *testBroker) Bind(instanceid, bindingid, serviceid, planid string) error {
	t.Logf("Creating service binding %s for service %s plan %s instance %s\n",
		bindingid, serviceid, planid, instanceid)
	return nil
}

func (t *testBroker) Unbind(instanceid, bindingid, serviceid, planid string) error {
	t.Logf("Delete service binding %s for service %s plan %s instance %s\n",
		bindingid, serviceid, planid, instanceid)
	return nil
}

func testAPI(t *testing.T) http.Handler {
	b := testBroker{
		T: t,
	}
	return NewAPI(nil, &b, "admin", "admin")
}

func doRequest(t *testing.T, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	testAPI(t).ServeHTTP(rec, req)
	return rec
}

func validJson(t *testing.T, response []byte, url string) {
	m := make(map[string]interface{})
	if json.Unmarshal(response, &m) != nil {
		t.Error(url, "should return a valid json")
	}
}

func TestRespondHasApplicationJSONContentType(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v2/catalog", nil)
	req.SetBasicAuth("admin", "admin")
	resp := doRequest(t, req)
	want, got := "application/json", resp.Header().Get("Content-Type")
	if !strings.Contains(got, want) {
		t.Fatalf("GET /v2/catalog: Content-Type: got: %v, want: %v", got, want)
	}
}
