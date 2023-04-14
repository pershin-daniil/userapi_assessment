package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"refactoring/pkg/models"
	"refactoring/pkg/store"
)

type Service interface {
	Users() (models.UserStore, error)
	User(id string) (models.User, error)
	NewUser(user models.UserRequest) (models.User, error)
	UpdateUser(id string, newName string) (models.User, error)
	DeleteUser(id string) error
}

func (s *Server) rootHandler(w http.ResponseWriter, _ *http.Request) {
	s.writeResponse(w, http.StatusOK, []byte(time.Now().String()))
}

func (s *Server) usersHandler(w http.ResponseWriter, _ *http.Request) {
	users, err := s.service.Users()
	if err != nil {
		s.log.Warnf("err during getting users: %v", err)
		s.writeResponse(w, http.StatusInternalServerError, err)
	}
	s.writeResponse(w, http.StatusOK, users)
}

func (s *Server) userHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParamFromCtx(ctx, "id")
	user, err := s.service.User(id)
	switch {
	case errors.Is(err, store.ErrUserNotFound):
		s.writeResponse(w, http.StatusNotFound, err)
		return
	case err != nil:
		s.log.Warnf("err during getting user: %v", err)
		s.writeResponse(w, http.StatusInternalServerError, err)
		return
	}
	s.writeResponse(w, http.StatusOK, user)
}

func (s *Server) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		s.writeResponse(w, http.StatusBadRequest, err)
		return
	}
	newUser, err := s.service.NewUser(user)
	if err != nil {
		s.log.Warnf("err during creating user: %v", err)
		s.writeResponse(w, http.StatusInternalServerError, err)
		return
	}
	s.writeResponse(w, http.StatusCreated, newUser)
}

func (s *Server) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParamFromCtx(ctx, "id")
	var user models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		s.writeResponse(w, http.StatusBadRequest, err)
		return
	}
	updatedUser, err := s.service.UpdateUser(id, user.DisplayName)
	switch {
	case errors.Is(err, store.ErrUserNotFound):
		s.writeResponse(w, http.StatusNotFound, err)
		return
	case err != nil:
		s.log.Warnf("err during updating user: %v", err)
		s.writeResponse(w, http.StatusInternalServerError, err)
		return
	}
	s.writeResponse(w, http.StatusOK, updatedUser)
}

func (s *Server) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParamFromCtx(ctx, "id")
	err := s.service.DeleteUser(id)
	switch {
	case errors.Is(err, store.ErrUserNotFound):
		s.writeResponse(w, http.StatusNotFound, err)
	case err != nil:
		s.log.Warnf("err during deleting user: %v", err)
		s.writeResponse(w, http.StatusInternalServerError, err)
		return
	}
	s.writeResponse(w, http.StatusNoContent, err)
}

func (s *Server) writeResponse(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	if x, ok := data.(error); ok {
		if err := json.NewEncoder(w).Encode(ErrorResponse{Error: x.Error()}); err != nil {
			s.log.Warnf("write response failed: %v", err)
		}
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.log.Warnf("write response failed: %v", err)
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}
