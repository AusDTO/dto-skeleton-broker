package broker

type MockBroker struct {
}

func (b *MockBroker) Provision(instanceid string) error {
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
