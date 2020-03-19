package api

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/jbruno/dumbwaiter/controller"
)

// HTTPController is the structure that implements the HTTP api for the controller service
type HTTPController struct {
	Controller *controller.Controller
}

// AddEndpoints adds the http endpoints to the server
func (c *HTTPController) AddEndpoints(router *mux.Router) {
	log.Info("adding endpoints")
}
