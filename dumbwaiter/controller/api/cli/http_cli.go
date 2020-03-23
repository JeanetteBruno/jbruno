package cli

type ControllerHttpCli struct {
	addr string // the url (with port to use when communicating with the controller)
}

func NewControllerHttpCli(addr string) *ControllerHttpCli {
	return &ControllerHttpCli{addr: addr}
}

func (c *ControllerHttpCli) SetRequestedFloor(floor int) {
	// TODO implement rest call to controller's SetRequestedFloor entry
}

func (c *ControllerHttpCli) SetLastSeenFloor(floor int) {
	// TODO implement rest call to controller's SetLastSeenFloor entry
}
