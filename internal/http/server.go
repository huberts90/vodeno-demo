package http

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
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
	apiRouter.HandleFunc("/messages/{id:[0-9]+}", srv.DeleteMessage).Methods(http.MethodDelete)
	apiRouter.HandleFunc("/messages/send", srv.PostOrderMailing).Methods(http.MethodPost)

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
	err = server.store.InsertMessage(ctx, message)
	if err != nil {
		server.logger.Error("Failed to create a message", zap.Error(err))
		http.Error(w, "Failed to create a message", http.StatusInternalServerError)
		return
	}

	// TODO: Change a way of logging body
	server.logger.Info("Message created", zap.Any("Body", string(b)))

	w.WriteHeader(http.StatusCreated)
}

func (server *Server) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	server.logger.Info(r.Method, zap.String("Path", r.URL.Path))

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	params := mux.Vars(r)
	idStr, ok := params["id"]
	if !ok {
		server.logger.Info("Missing message id")
		http.Error(w, "missing message id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if !ok {
		server.logger.Info("Message id is not a number")
		http.Error(w, "message id is not a number", http.StatusBadRequest)
		return
	}

	err = server.store.DeleteMessage(ctx, id)
	if err != nil {
		server.logger.Error("Failed to delete a message", zap.Error(err))
		http.Error(w, "failed to delete a message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (server *Server) PostOrderMailing(w http.ResponseWriter, r *http.Request) {
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

	job := &store.MailingJob{}
	err = json.Unmarshal(b, job)
	if err != nil {
		server.logger.Info("Failed to unmarshall body content", zap.Any("Body", string(b)), zap.Error(err))
		http.Error(w, "failed to unmarshall body content", http.StatusBadRequest)
		return
	}

	// TODO: validate input
	err = server.store.OrderMailing(ctx, job)
	if err != nil {
		server.logger.Error("Failed to order a mailing job", zap.Error(err))
		http.Error(w, "Failed to order a mailing job", http.StatusInternalServerError)
		return
	}

	// TODO: Change a way of logging body
	server.logger.Info("Mailing job scheduled", zap.Any("Body", string(b)))

	w.WriteHeader(http.StatusCreated)
}
