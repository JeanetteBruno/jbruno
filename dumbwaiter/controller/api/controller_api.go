package api

type Controller interface {
	SetRequestedFloor(floor int)
	SetLastSeenFloor(floor int)
}
