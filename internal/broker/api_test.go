package broker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testBroker struct {
	*testing.T

	provisionfn   func(t *testing.T, instanceid, serviceid, planid string) error
	deprovisionfn func(t *testing.T, instanceid, serviceid, planid string) error
	bindfn        func(t *testing.T, instanceid, bindingid, serviceid, planid string) error
	unbindfn      func(t *testing.T, instanceid, bindingid, serviceid, planid string) error
}

func (t *testBroker) Provision(instanceid, serviceid, planid string) error {
	return t.provisionfn(t.T, instanceid, serviceid, planid)
}

func (t *testBroker) Deprovision(instanceid, serviceid, planid string) error {
	return t.deprovisionfn(t.T, instanceid, serviceid, planid)
}

func (t *testBroker) Bind(instanceid, bindingid, serviceid, planid string) error {
	return t.bindfn(t.T, instanceid, bindingid, serviceid, planid)
}

func (t *testBroker) Unbind(instanceid, bindingid, serviceid, planid string) error {
	return t.unbindfn(t.T, instanceid, bindingid, serviceid, planid)
}

type option func(t *testBroker)

func testAPI(t *testing.T, opts ...option) http.Handler {
	b := testBroker{
		T: t,
	}
	for _, o := range opts {
		o(&b)
	}
	return NewAPI(nil, &b, "admin", "admin")
}

type eofReadCloser struct{}

func (r *eofReadCloser) Read([]byte) (int, error) { return 0, io.EOF }
func (r *eofReadCloser) Close() error             { return io.EOF }

func doRequest(api http.Handler, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()

	// ensure req.Body is not nil. This is guarenteed by the
	// net/http api, but as we're calling the handler's ServeHTTP
	// method directly we bypass that guarantee.
	if req.Body == nil {
		req.Body = new(eofReadCloser)
	}

	api.ServeHTTP(rec, req)
	return rec
}

func validJson(t *testing.T, response []byte, url string) {
	m := make(map[string]interface{})
	if json.Unmarshal(response, &m) != nil {
		t.Error(url, "should return a valid json")
	}
}

func auth(user, pass string) func(*http.Request) {
	return func(req *http.Request) {
		req.SetBasicAuth(user, pass)
	}
}

func body(v interface{}) func(*http.Request) {
	var r io.Reader
	switch v := v.(type) {
	case string:
		r = strings.NewReader(v)
	case []byte:
		r = bytes.NewReader(v)
	case map[string]interface{}:
		buf, err := json.Marshal(v)
		if err != nil {
			panic(err.Error())
		}
		r = bytes.NewReader(buf)
	default:
		panic(fmt.Sprintf("body cannot handle type %T", v))
	}

	return func(req *http.Request) {
		req.Body = ioutil.NopCloser(r)
	}
}

func request(method, url string, opts ...func(*http.Request)) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err.Error())
	}
	for _, o := range opts {
		o(req)
	}
	return req
}

func TestRespondHasApplicationJSONContentType(t *testing.T) {
	api := testAPI(t)
	req := request("GET", "/v2/catalog", auth("admin", "admin"))
	resp := doRequest(api, req)
	want, got := "application/json", resp.Header().Get("Content-Type")
	if !strings.Contains(got, want) {
		t.Fatalf("GET /v2/catalog: Content-Type: got: %v, want: %v", got, want)
	}
}

func successfulProvisioning(serviceid string) func(*testBroker) {
	return func(b *testBroker) {
		b.provisionfn = func(t *testing.T, instanceid, serviceid, planid string) error {
			if serviceid != serviceid {
				return errors.New("unknown service id")
			}
			return nil
		}
	}
}

func TestProvision(t *testing.T) {
	api := testAPI(t, successfulProvisioning("a_service_guid"))
	req := request("PUT", "/v2/service_instances/a_guid",
		auth("admin", "admin"),
		body(map[string]interface{}{
			"organization_guid": "org-guid-here",
			"plan_id":           "plan-guid-here",
			"service_id":        "service-guid-here",
			"space_guid":        "space-guid-here",
		}))
	resp := doRequest(api, req)
	if resp.Code != 201 {
		t.Fatalf("%s returned status: %d, expected: %d", req, resp.Code, 201)
	}
}

func TestProvisionMissingPlanId(t *testing.T) {
	api := testAPI(t, successfulProvisioning("a_service_guid"))
	req := request("PUT", "/v2/service_instances/a_guid",
		auth("admin", "admin"),
		body(map[string]interface{}{
			"organization_guid": "org-guid-here",
			"service_id":        "service-guid-here",
			"space_guid":        "space-guid-here",
		}))
	resp := doRequest(api, req)
	if resp.Code != 504 {
		t.Fatalf("%s returned status: %d, expected: %d", req, resp.Code, 504)
	}
}
