package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/JeanetteBruno/jbruno/dumbwaiter/controller"
)

// HTTPController is the structure for adding the controller's restian endpoints
type HTTPController struct {
	Controller  *controller.Controller
	ServiceName string
}

// NewHTTPController wrap the controller with http entrypoints
func NewHTTPController(controller *controller.Controller) *HTTPController {
	return &HTTPController{Controller: controller, ServiceName: "controller"}
}

// AddEndpoints adds the http endpoints to the server
func (c *HTTPController) AddEndpoints(router *mux.Router) {
	log.Info("adding controller service endpoints")
	router.HandleFunc(fmt.Sprintf("/%s/status", c.ServiceName), c.StatusEndpoint).Methods("GET")
}

// StatusEndpoint implement the http entry for status requests
func (c *HTTPController) StatusEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Info("StatusEndpoint request received")

	status := c.Controller.GetStatus()

	w.Header().Add("Content-Type", "application/json")
	log.Infof("StatusEndpoint returning: %v", status)
	json.NewEncoder(w).Encode(status)
}
