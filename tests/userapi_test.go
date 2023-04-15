package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"

	"refactoring/internal/logger"
	"refactoring/internal/server"
	"refactoring/pkg/models"
	"refactoring/pkg/service"
	"refactoring/pkg/store"
)

const (
	port      = ":4444"
	version   = "test"
	storePath = `users_test.json`
	testURL   = "http://localhost" + port
)

type IntegrationTestSuite struct {
	suite.Suite
	log         *logrus.Logger
	store       *store.Store
	service     *service.Service
	server      *server.Server
	userRequest models.UserRequest
	updateData  models.UserRequest
}

func (s *IntegrationTestSuite) SetupSuite() {
	ctx := context.Background()
	s.log = logger.New()
	s.store = store.New(s.log, storePath)
	s.service = service.New(s.log, s.store)
	s.server = server.New(s.log, s.service, port, version)
	go func() {
		_ = s.server.Run(ctx)
	}()
	time.Sleep(100 * time.Millisecond)
	data := models.UserStore{
		Increment: 0,
		List:      map[string]models.User{},
	}
	initFile, err := json.Marshal(data)
	s.Require().NoError(err)
	err = os.WriteFile(storePath, initFile, fs.ModePerm)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) SetupTest() {
	s.userRequest = models.UserRequest{
		DisplayName: "Ivan Ivanov",
		Email:       "test@email.com",
	}
	s.updateData = models.UserRequest{DisplayName: "Masha"}
}

func (s *IntegrationTestSuite) TestMainWorkFlow() {
	s.Run("create user", func() {
		ctx := context.Background()
		var respBody models.User
		resp := s.sendRequest(ctx, http.MethodPost, "/api/v1/users", s.userRequest, &respBody)
		s.Require().Equal(http.StatusCreated, resp.StatusCode)
		s.Require().NotEqual(0, respBody.ID)
		s.Require().Equal(s.userRequest.DisplayName, respBody.DisplayName)
		s.Require().Equal(s.userRequest.Email, respBody.Email)
		s.Require().Greater(respBody.Created, time.Now().Add(-time.Minute))
		s.Require().Greater(respBody.Updated, time.Now().Add(-time.Minute))
		s.Require().Equal(respBody.Created, respBody.Updated)
	})

	s.Run("update user", func() {
		ctx := context.Background()
		var respBody models.User
		resp := s.sendRequest(ctx, http.MethodPatch, "/api/v1/users/1", s.updateData, &respBody)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
		s.Require().NotEqual(0, respBody.ID)
		s.Require().Equal(s.updateData.DisplayName, respBody.DisplayName)
		s.Require().NotEqual("", respBody.Email)
		s.Require().Greater(respBody.Updated, respBody.Created)
	})

	s.Run("create one more user", func() {
		ctx := context.Background()
		var respBody models.User
		s.userRequest.DisplayName = "Kar Karich"
		s.userRequest.Email = "poop@mail.com"
		resp := s.sendRequest(ctx, http.MethodPost, "/api/v1/users", s.userRequest, &respBody)
		s.Require().Equal(http.StatusCreated, resp.StatusCode)
		s.Require().NotEqual(0, respBody.ID)
		s.Require().Equal(s.userRequest.DisplayName, respBody.DisplayName)
		s.Require().Equal(s.userRequest.Email, respBody.Email)
		s.Require().Greater(respBody.Created, time.Now().Add(-time.Minute))
		s.Require().Greater(respBody.Updated, time.Now().Add(-time.Minute))
		s.Require().Equal(respBody.Created, respBody.Updated)
	})

	s.Run("get user by id", func() {
		ctx := context.Background()
		var respBody models.User
		resp := s.sendRequest(ctx, http.MethodGet, "/api/v1/users/2", nil, &respBody)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
		s.Require().Equal(2, respBody.ID)
		s.Require().Equal(s.userRequest.DisplayName, respBody.DisplayName)
		s.Require().Equal(s.userRequest.Email, respBody.Email)
		s.Require().Greater(respBody.Created, time.Now().Add(-time.Minute))
		s.Require().Greater(respBody.Updated, time.Now().Add(-time.Minute))
	})

	s.Run("get list of users", func() {
		ctx := context.Background()
		var respBody models.UserStore
		resp := s.sendRequest(ctx, http.MethodGet, "/api/v1/users", nil, &respBody)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
		s.Require().Equal(2, len(respBody.List))
	})

	s.Run("delete user", func() {
		ctx := context.Background()
		resp := s.sendRequest(ctx, http.MethodDelete, "/api/v1/users/1", nil, nil)
		s.Require().Equal(http.StatusNoContent, resp.StatusCode)
	})

	s.Run("get list of users after deletion", func() {
		ctx := context.Background()
		var respBody models.UserStore
		resp := s.sendRequest(ctx, http.MethodGet, "/api/v1/users", nil, &respBody)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
		s.Require().Equal(1, len(respBody.List))
	})

	s.Run("user Not Found", func() {
		ctx := context.Background()
		var respBody models.User
		resp := s.sendRequest(ctx, http.MethodGet, "/api/v1/users/0", nil, &respBody)
		s.Require().Equal(http.StatusNotFound, resp.StatusCode)
	})
}

func (s *IntegrationTestSuite) sendRequest(ctx context.Context, method, endpoint string, body interface{}, dest interface{}) *http.Response {
	s.T().Helper()
	reqBody, err := json.Marshal(body)
	s.Require().NoError(err)
	req, err := http.NewRequestWithContext(ctx, method, testURL+endpoint, bytes.NewReader(reqBody))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err)
	defer func() {
		err = resp.Body.Close()
		s.Require().NoError(err)
	}()
	if dest != nil {
		err = json.NewDecoder(resp.Body).Decode(&dest)
		s.Require().NoError(err)
		s.T().Logf("%#v", &dest)
	}
	return resp
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
