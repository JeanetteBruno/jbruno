package httpservice

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// serviceObject is the standard interface domain's http wrappers must implement
// adding the service specific rest endpoints that will call into domain object
// functions
type serviceObject interface {
	AddEndpoints(router *mux.Router)
}

// Service is the common object that wraps all of the applications domain object adding http entry points
// and turning the domain object into a restian service
type Service struct {
	domainObject serviceObject
	httpAddr     string
	srv          *http.Server
	serviceName  string
}

// NewService create an restian Service object to wrap the domain logic
func NewService(domainObject serviceObject, httpAddr string, serviceName string) *Service {
	s := &Service{
		domainObject: domainObject,
		httpAddr:     httpAddr,
		serviceName:  serviceName,
	}

	return s
}

// RunService starts the service listening for requests
func (s *Service) RunService() {
	// add the request handler and endpoints
	router := mux.NewRouter()
	s.addCommonEndpoints(router)
	s.domainObject.AddEndpoints(router)

	// create the http service object
	s.srv = &http.Server{
		Handler: router,
		Addr:    s.httpAddr,
	}

	// start the object listening for requests
	err := s.srv.ListenAndServe()
	log.Fatal(err)
}

// addCommonEndpoints add endpoints common to all of the services
func (s *Service) addCommonEndpoints(router *mux.Router) {
	log.Info("adding standard endpoints")
	router.HandleFunc(fmt.Sprintf("%s/health", s.serviceName), s.HealthEndpoint).Methods("GET")
}

// HealthEndpoint will return ok if the service is running, otherwise the caller should get 404
func (s *Service) HealthEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Info("Health request received")

	status := "ok"

	w.Header().Add("Content-Type", "application/json")
	log.Infof("StatusEndpoint returning: %v", status)
	json.NewEncoder(w).Encode(status)
}
