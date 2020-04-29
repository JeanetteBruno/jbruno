package cli

// ControllerHTTPClient controller client implementing http calls to the controller
type ControllerHTTPClient struct {
	addr string // the url (with port to use when communicating with the controller)
}

// NewControllerHTTPClient instantiate an http client for communicating with the controller
func NewControllerHTTPClient(addr string) *ControllerHTTPClient {
	return &ControllerHTTPClient{addr: addr}
}

// SetRequestedFloor send a floor request to the controller
func (c *ControllerHTTPClient) SetRequestedFloor(floor int) {
	// TODO implement rest call to controller's SetRequestedFloor entry
}

// SetFloorStatus tell the controller that the platform has arrived at a floor
func (c *ControllerHTTPClient) SetFloorStatus(floor int, atFloor bool) {
	// TODO implement rest call to controller's SetFloorStatus entry
}

//SetStopRequested tell the controller to stop
func (c *ControllerHTTPClient) SetStopRequested() {
	//TODO implement send stop request controller
}
