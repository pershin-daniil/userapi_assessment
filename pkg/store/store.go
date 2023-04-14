package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"refactoring/pkg/models"
)

var ErrUserNotFound = errors.New("user_not_found")

type Store struct {
	log   *logrus.Entry
	store string
}

func New(log *logrus.Logger, storePath string) *Store {
	return &Store{
		log:   log.WithField("module", "store"),
		store: storePath,
	}
}

func (s *Store) SearchUsers() (models.UserStore, error) {
	file, err := os.ReadFile(s.store)
	if err != nil {
		return models.UserStore{}, fmt.Errorf("serch users faild: %w", err)
	}
	var store models.UserStore
	if err = json.Unmarshal(file, &store); err != nil {
		return models.UserStore{}, fmt.Errorf("serch users faild: %w", err)
	}
	return store, nil
}

func (s *Store) CreateUser(user models.UserRequest) (models.User, error) {
	store, err := s.connectStore(s.store)
	if err != nil {
		return models.User{}, fmt.Errorf("create user faild: %w", err)
	}
	store.Increment++
	t := time.Now()
	newUser := models.User{
		ID:          store.Increment,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Created:     t,
		Updated:     t,
	}
	id := strconv.Itoa(newUser.ID)
	store.List[id] = newUser
	result, err := json.Marshal(&store)
	if err != nil {
		return models.User{}, fmt.Errorf("create user faild: %w", err)
	}
	if err = os.WriteFile(s.store, result, fs.ModePerm); err != nil {
		return models.User{}, fmt.Errorf("create user faild: %w", err)
	}
	return newUser, nil
}

func (s *Store) User(id string) (models.User, error) {
	store, err := s.connectStore(s.store)
	if err != nil {
		return models.User{}, fmt.Errorf("get user faild: %w", err)
	}
	var user models.User
	user, ok := store.List[id]
	if !ok {
		return models.User{}, fmt.Errorf("get user faild: %w", ErrUserNotFound)
	}
	return user, nil
}

func (s *Store) UpdateUser(id string, newName string) (models.User, error) {
	store, err := s.connectStore(s.store)
	if err != nil {
		return models.User{}, fmt.Errorf("update user faild: %w", err)
	}
	user, ok := store.List[id]
	if !ok {
		return models.User{}, fmt.Errorf("update user faild: %w", ErrUserNotFound)
	}
	user.DisplayName = newName
	user.Updated = time.Now()
	store.List[id] = user
	result, err := json.Marshal(&store)
	if err != nil {
		return models.User{}, fmt.Errorf("update user faild: %w", err)
	}
	if err = os.WriteFile(s.store, result, fs.ModePerm); err != nil {
		return models.User{}, fmt.Errorf("update user faild: %w", err)
	}
	return user, nil
}

func (s *Store) DeleteUser(id string) error {
	store, err := s.connectStore(s.store)
	if err != nil {
		return fmt.Errorf("delete user faild: %w", err)
	}
	if _, ok := store.List[id]; !ok {
		return fmt.Errorf("delete user faild: %w", ErrUserNotFound)
	}
	delete(store.List, id)
	result, err := json.Marshal(&store)
	if err != nil {
		return fmt.Errorf("delete user faild: %w", err)
	}
	if err = os.WriteFile(s.store, result, fs.ModePerm); err != nil {
		return fmt.Errorf("dele user faild: %w", err)
	}
	return nil
}

func (s *Store) connectStore(path string) (models.UserStore, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return models.UserStore{}, fmt.Errorf("connect store faild: %w", err)
	}
	var store models.UserStore
	if err = json.Unmarshal(file, &store); err != nil {
		return models.UserStore{}, fmt.Errorf("connect store faild: %w", err)
	}
	return store, nil
}
