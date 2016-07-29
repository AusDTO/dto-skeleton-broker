package broker

import "fmt"

type MockBroker struct {
}

func (b *MockBroker) Provision(instanceid, serviceid, planid string) error {
	fmt.Printf("Creating service instance %s for service %s plan %s\n", instanceid, serviceid, planid)

	return nil
}

func (b *MockBroker) Deprovision(instanceid string) error {
	return nil
}

func (b *MockBroker) Bind(instanceid, bindingid string) error {
	return nil
}

func (b *MockBroker) Unbind(instanceid, bindingid string) error {
	return nil
}
