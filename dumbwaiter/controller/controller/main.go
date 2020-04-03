package main

import (
	"flag"

	"github.com/JeanetteBruno/jbruno/dumbwaiter/common/httpservice"
	"github.com/JeanetteBruno/jbruno/dumbwaiter/controller"
	"github.com/JeanetteBruno/jbruno/dumbwaiter/controller/api"
)

var (
	httpAddrFlag = flag.String("http_addr", "localhost:9090", "host:port to serve http api on")
	numFloors    = flag.Int("num_floors", 3, "the number of floors the dumbwaiter will serve")
)

// start the service.
func main() {
	// parse flags
	flag.Parse()

	s := newControllerHTTPService(*httpAddrFlag, *numFloors) // create the controller with http nature

	s.RunService() // start the controller listening for requests
}

func newControllerHTTPService(httpAddr string, numFloors int) *httpservice.Service {
	// construct controller object and start its processing loop
	controller := controller.NewController(numFloors)
	controller.StartProcessingLoop()

	// add the http endpoints
	httpController := api.NewHTTPController(controller)

	// add the final (common) http nature
	s := httpservice.NewService(httpController, httpAddr, httpController.ServiceName)

	return s
}
