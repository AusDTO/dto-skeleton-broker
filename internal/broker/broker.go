package broker

import (
	"github.com/pkg/errors"
)

// A Broker represents a Cloud Foundry Service Broker
type Broker interface {

	// Provision creates an instance of the serviceid, planid pair
	// and associates it with the provided instanceid.
	Provision(instanceid, serviceid, planid string) error

	// Deprovision removes an existing service instance.
	Deprovision(instanceid, serviceid, planid string) error

	// Bind requests the creation of a service instance binding.
	Bind(instanceid, bindingid, serviceid, planid string) error

	// Unbind requests the destructions of a service instance binding
	Unbind(instanceid, bindingid, serviceid, planid string) error
}

// validatingBroker implements the Broker API and validates each parameter
// passed to the underlying Broker. This frees the real broker implementations
// from having to deal with invalid input.
type validatingBroker struct {
	Broker
}

func (b *validatingBroker) Provision(instanceid, serviceid, planid string) error {
	if instanceid == "" {
		return errors.New("instance_id is blank")
	}
	if serviceid == "" {
		return errors.New("service_id is blank")
	}
	if planid == "" {
		return errors.New("plan_id is blank")
	}
	return b.Broker.Provision(instanceid, serviceid, planid)
}

func (b *validatingBroker) Deprovision(instanceid, serviceid, planid string) error {
	if instanceid == "" {
		return errors.New("instance_id is blank")
	}
	if serviceid == "" {
		return errors.New("service_id is blank")
	}
	if planid == "" {
		return errors.New("plan_id is blank")
	}
	return b.Broker.Deprovision(instanceid, serviceid, planid)
}

func (b *validatingBroker) Bind(instanceid, bindingid, serviceid, planid string) error {
	if instanceid == "" {
		return errors.New("instance_id is blank")
	}
	if bindingid == "" {
		return errors.New("binding_id is blank")
	}
	if serviceid == "" {
		return errors.New("service_id is blank")
	}
	if planid == "" {
		return errors.New("plan_id is blank")
	}
	return b.Broker.Bind(instanceid, bindingid, serviceid, planid)
}

func (b *validatingBroker) Unbind(instanceid, bindingid, serviceid, planid string) error {
	if instanceid == "" {
		return errors.New("instance_id is blank")
	}
	if bindingid == "" {
		return errors.New("binding_id is blank")
	}
	if serviceid == "" {
		return errors.New("service_id is blank")
	}
	if planid == "" {
		return errors.New("plan_id is blank")
	}
	return b.Broker.Unbind(instanceid, bindingid, serviceid, planid)
}
