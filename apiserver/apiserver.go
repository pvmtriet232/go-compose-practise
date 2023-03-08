// apiserver/apiserver.go
package apiserver

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pvmtriet232/go-compose-practise/storage"
	"github.com/sirupsen/logrus"
)

var defaultStopTimeout = time.Second * 30

type Endpoint struct {
	handler EndpointFunc
}

type EndpointFunc func(w http.ResponseWriter, req *http.Request) error

func (e Endpoint) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := e.handler(w, req); err != nil {
		logrus.WithError(err).Error("could not process request")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}
}

type APIServer struct {
	addr    string
	storage *storage.Storage
}

func NewAPIServer(addr string, storage *storage.Storage) (*APIServer, error) {
	if addr == "" {
		return nil, errors.New("addr cannot be blank")
	}

	return &APIServer{
		addr:    addr,
		storage: storage,
	}, nil
}

// Start starts a server with a stop channel
func (s *APIServer) Start(stop <-chan struct{}) error {
	srv := &http.Server{
		Addr:    s.addr,
		Handler: s.router(),
	}

	go func() {
		logrus.WithField("addr", srv.Addr).Info("starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), defaultStopTimeout)
	defer cancel()

	logrus.WithField("timeout", defaultStopTimeout).Info("stopping server")
	return srv.Shutdown(ctx)
}

func (s *APIServer) router() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", s.defaultRoute)
	router.Methods("POST").Path("/items").Handler(Endpoint{s.createItem})
	router.Methods("GET").Path("/items").Handler(Endpoint{s.listItems})
	return router
}

func (s *APIServer) defaultRoute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello tao ne "))
}
