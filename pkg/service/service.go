package service

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"refactoring/pkg/models"
)

type Store interface {
	SearchUsers() (models.UserStore, error)
	CreateUser(user models.UserRequest) (models.User, error)
	User(id string) (models.User, error)
	UpdateUser(id string, newName string) (models.User, error)
	DeleteUser(id string) error
}

type Service struct {
	log   *logrus.Entry
	store Store
}

func New(log *logrus.Logger, store Store) *Service {
	return &Service{
		log:   log.WithField("module", "service"),
		store: store,
	}
}

func (s *Service) Users() (models.UserStore, error) {
	users, err := s.store.SearchUsers()
	if err != nil {
		return models.UserStore{}, fmt.Errorf("service: %w", err)
	}
	return users, nil
}

func (s *Service) User(id string) (models.User, error) {
	user, err := s.store.User(id)
	if err != nil {
		return models.User{}, fmt.Errorf("service: %w", err)
	}
	return user, nil
}

func (s *Service) NewUser(user models.UserRequest) (models.User, error) {
	newUser, err := s.store.CreateUser(user)
	if err != nil {
		return models.User{}, fmt.Errorf("service: %w", err)
	}
	return newUser, nil
}

func (s *Service) UpdateUser(id string, newName string) (models.User, error) {
	user, err := s.store.UpdateUser(id, newName)
	if err != nil {
		return models.User{}, fmt.Errorf("service: %w", err)
	}
	return user, nil
}

func (s *Service) DeleteUser(id string) error {
	if err := s.store.DeleteUser(id); err != nil {
		return fmt.Errorf("service: %w", err)
	}
	return nil
}
