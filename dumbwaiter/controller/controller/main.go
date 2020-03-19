package main

import (
	"flag"

	http "github.com/jbruno/dumbwaiter/common/httpservice"
	"github.com/jbruno/dumbwaiter/controller"
	"github.com/jbruno/dumbwaiter/controller/api"
)

var (
	httpAddrFlag = flag.String("http_addr", "localhost:9090", "host:port to serve http api on")
	logLevelFlag = flag.String("log_level", "info", "log at this level and above (debug|info|error)")
	numFloors    = flag.Int("num_floors", 3, "the number of floors the dumbwaiter will serve")
)

// start the service.
func main() {
	// parse flags and initialize logging
	flag.Parse()

	s := newHTTPService(*httpAddrFlag, *numFloors)

	s.RunService()
}

// construct an http-based pusher service
func newHTTPService(httpAddr string, numFloors int) *http.Service {
	controllerDO := controller.NewController(numFloors)     // construct the domain object
	service := &api.HTTPController{Controller: controllerDO} // wrap the controller with http entry points
	s := http.NewService(service, httpAddr)                 // create the http service object

	return s
}
