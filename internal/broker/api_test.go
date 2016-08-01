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

func successfulProvision(serviceid, planid string) func(*testBroker) {
	return func(b *testBroker) {
		b.provisionfn = func(t *testing.T, instanceid, serviceid, planid string) error {
			if serviceid != serviceid {
				return errors.New("unknown service id")
			}
			if planid != planid {
				return errors.New("unknown plan id")
			}
			return nil
		}
	}
}

func successfulDeprovision(serviceid, planid string) func(*testBroker) {
	return func(b *testBroker) {
		b.deprovisionfn = func(t *testing.T, instanceid, serviceid, planid string) error {
			if serviceid != serviceid {
				return errors.New("unknown service id")
			}
			if planid != planid {
				return errors.New("unknown plan id")
			}
			return nil
		}
	}
}

func successfulBind(serviceid, planid, bindingid string) func(*testBroker) {
	return func(b *testBroker) {
		b.bindfn = func(t *testing.T, instanceid, serviceid, planid, bindingid string) error {
			if serviceid != serviceid {
				return errors.New("unknown service id")
			}
			if planid != planid {
				return errors.New("unknown plan id")
			}
			if bindingid != bindingid {
				return errors.New("unknown binding id")
			}
			return nil
		}
	}
}

func successfulUnbind(serviceid, planid, bindingid string) func(*testBroker) {
	return func(b *testBroker) {
		b.unbindfn = func(t *testing.T, instanceid, serviceid, planid, bindingid string) error {
			if serviceid != serviceid {
				return errors.New("unknown service id")
			}
			if planid != planid {
				return errors.New("unknown plan id")
			}
			if bindingid != bindingid {
				return errors.New("unknown binding id")
			}
			return nil
		}
	}
}

const (
	INST_ID    = `daa4dbef-a861-42a7-b1a3-b161df0b4eb0`
	ORG_GUID   = `74a00865-cc31-4360-98ab-728e6fd4eacd`
	PLAN_ID    = `da71b52f-a93e-48cb-968b-123e44b19320`
	SERVICE_ID = `513c3e8e-aa17-48cf-81d0-338c27c06e48`
	SPACE_GUID = `56635799-ea54-44bc-bd34-6d682ca191e0`
	BINDING_ID = `bcc24d05-e8ae-4231-b60b-55e95a44c6f5`
)

func TestProvision(t *testing.T) {
	api := testAPI(t, successfulProvision(SERVICE_ID, PLAN_ID))
	req := request("PUT", "/v2/service_instances/"+INST_ID,
		auth("admin", "admin"),
		body(map[string]interface{}{
			"organization_guid": ORG_GUID,
			"plan_id":           PLAN_ID,
			"service_id":        SERVICE_ID,
			"space_guid":        SPACE_GUID,
		}))
	resp := doRequest(api, req)
	if resp.Code != 201 {
		t.Fatalf("%s returned status: %d, expected: %d", req, resp.Code, 201)
	}
}

func TestProvisionMissingPlanId(t *testing.T) {
	api := testAPI(t, successfulProvision(SERVICE_ID, PLAN_ID))
	req := request("PUT", "/v2/service_instances/"+INST_ID,
		auth("admin", "admin"),
		body(map[string]interface{}{
			"organization_guid": ORG_GUID,
			"service_id":        SERVICE_ID,
			"space_guid":        SPACE_GUID,
		}))
	resp := doRequest(api, req)
	if resp.Code != 504 {
		t.Fatalf("%s returned status: %d, expected: %d", req, resp.Code, 504)
	}
}

func TestDeprovision(t *testing.T) {
	api := testAPI(t, successfulDeprovision(SERVICE_ID, PLAN_ID))
	req := request("DELETE", "/v2/service_instances/"+INST_ID+"?service_id="+SERVICE_ID+"&plan_id="+PLAN_ID,
		auth("admin", "admin"))
	resp := doRequest(api, req)
	if resp.Code != 200 {
		t.Fatalf("%s returned status: %d, expected: %d: %s", req, resp.Code, 200, resp.Body.String())
	}
}

func TestBind(t *testing.T) {
	api := testAPI(t, successfulBind(SERVICE_ID, PLAN_ID, BINDING_ID))
	req := request("PUT", "/v2/service_instances/"+INST_ID+"/service_bindings/"+BINDING_ID,
		auth("admin", "admin"),
		body(map[string]interface{}{
			"plan_id":    PLAN_ID,
			"service_id": SERVICE_ID,
		}))
	resp := doRequest(api, req)
	if resp.Code != 201 {
		t.Fatalf("%s returned status: %d, expected: %d: %s", req, resp.Code, 201, resp.Body.String())
	}
}

func TestUnbind(t *testing.T) {
	api := testAPI(t, successfulUnbind(SERVICE_ID, PLAN_ID, BINDING_ID))
	req := request("DELETE", "/v2/service_instances/"+INST_ID+"/service_bindings/"+BINDING_ID+
		"?service_id="+SERVICE_ID+"&plan_id="+PLAN_ID,
		auth("admin", "admin"),
	)
	resp := doRequest(api, req)
	if resp.Code != 200 {
		t.Fatalf("%s returned status: %d, expected: %d: %s", req, resp.Code, 200, resp.Body.String())
	}
}
