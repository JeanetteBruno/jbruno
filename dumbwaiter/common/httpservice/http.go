package httpservice

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type serviceObject interface {
	AddEndpoints(router *mux.Router)
}

// Service an http service wrapper for a domain object
type Service struct {
	domainObject serviceObject
	httpAddr     string
	srv          *http.Server
}

// NewService constructs a Service
func NewService(domainObject serviceObject, httpAddr string) *Service {
	s := &Service{
		domainObject: domainObject,
		httpAddr:     httpAddr,
	}
	return s
}

// RunService starts the http service
func (s *Service) RunService() {
	router := mux.NewRouter()

	s.domainObject.AddEndpoints(router)

	s.srv = &http.Server{
		Handler: router,
		Addr:    s.httpAddr,
	}

	err := s.srv.ListenAndServe()
	log.Error(err)
}
