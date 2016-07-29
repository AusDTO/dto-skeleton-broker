package broker

import "fmt"

type MockBroker struct {
}

func (b *MockBroker) Provision(instanceid, serviceid, planid string) error {
	fmt.Printf("Creating service instance %s for service %s plan %s\n", instanceid, serviceid, planid)

	return nil
}

func (b *MockBroker) Deprovision(instanceid, serviceid, planid string) error {
	fmt.Printf("Deleting service instance %s for service %s plan %s\n", instanceid, serviceid, planid)
	return nil
}

func (b *MockBroker) Bind(instanceid, bindingid, serviceid, planid string) error {
	fmt.Printf("Creating service binding %s for service %s plan %s instance %s\n",
		bindingid, serviceid, planid, instanceid)

	return nil
}

func (b *MockBroker) Unbind(instanceid, bindingid, serviceid, planid string) error {
	fmt.Printf("Delete service binding %s for service %s plan %s instance %s\n",
		bindingid, serviceid, planid, instanceid)
	return nil
}
