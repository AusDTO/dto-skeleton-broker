package broker

type MockBroker struct {
}

func (b *MockBroker) Bind(instanceid, bindingid string) error {
	return nil
}

func (b *MockBroker) Unbind(instanceid, bindingid string) error {
	return nil
}
