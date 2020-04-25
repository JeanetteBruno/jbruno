package api

// Controller clients should use this interface when interacting with the controller
// an http client that implements this interface will be provided.
type Controller interface {
	SetRequestedFloor(floor int)
	SetLastSeenFloor(floor int)
	SetStopRequested()
}
