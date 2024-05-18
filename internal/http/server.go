package http

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
	"vodeno.com/demo/internal/store"
)

type Server struct {
	logger *zap.Logger
	router *mux.Router
	store  store.MessengerStore
}

func NewServer(logger *zap.Logger, store store.MessengerStore) *Server {
	srv := &Server{
		logger: logger,
		router: mux.NewRouter(),
		store:  store,
	}

	apiRouter := srv.router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/messages", srv.PostMessage).Methods(http.MethodPost)

	return srv
}

func (server *Server) Serve(port string) error {
	server.logger.Info("Http server listening", zap.String("port", port))
	return http.ListenAndServe(":"+port, server.router)
}

func (server *Server) PostMessage(w http.ResponseWriter, r *http.Request) {
	server.logger.Info(r.Method, zap.String("Path", r.URL.Path))

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		server.logger.Info("Failed to read body content", zap.Error(err))
		http.Error(w, "failed to read body content", http.StatusBadRequest)
		return
	}

	message := &store.Message{}
	err = json.Unmarshal(b, message)
	if err != nil {
		server.logger.Info("Failed to unmarshall body content", zap.Any("Body", string(b)), zap.Error(err))
		http.Error(w, "failed to unmarshall body content", http.StatusBadRequest)
		return
	}

	// TODO: validate input

	err = server.store.Insert(ctx, message)
	if err != nil {
		server.logger.Error("Failed to create a message", zap.Error(err))
		http.Error(w, "Failed to create a message", http.StatusInternalServerError)
		return
	}

	// TODO: Change a way of logging body
	server.logger.Info("Message created", zap.Any("Body", string(b)))

	w.WriteHeader(http.StatusCreated)
}
