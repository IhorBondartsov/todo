// +build integration

package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"to-do/api"
	"to-do/app"
	"to-do/repository"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

var url = "http://0.0.0.0:8080"

func newHttpTestService() (*httpService, error) {
	db, err := repository.NewDBClient(context.Background(), repository.StorageConfig{
		Driver: "postgres",
		DSN:    "postgresql://postgres:changeme@localhost:5432/postgres?sslmode=disable",
	})
	if err != nil {
		return nil, err
	}
	service, err := app.NewToDoService(db)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return NewHTTPService(HTTPConfig{
		Host:          "0.0.0.0",
		Port:          8080,
		InitProfiling: false,
	}, service), nil
}

func TestCreateTodo(t *testing.T) {
	assert := assert.New(t)

	service, err := newHttpTestService()
	assert.NoError(err)

	tt := []struct {
		name       string
		method     string
		todo       api.ToDo
		statusCode int
	}{
		{
			name: "Return correct todo. 200 Ok",
			todo: api.ToDo{
				Message: "to do new thing",
				UserID:  1,
			},
			method:     http.MethodPost,
			statusCode: http.StatusCreated,
		},
		{
			name: "User id does not exist in db. 500",
			todo: api.ToDo{
				Message: "todoSmt",
				UserID:  95,
			},
			method: http.MethodPost,

			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			body, err := json.Marshal(tc.todo)
			assert.NoError(err)
			request := httptest.NewRequest(tc.method, url, bytes.NewReader(body))
			responseRecorder := httptest.NewRecorder()
			params := httprouter.Params{}

			service.createToDo(responseRecorder, request, params)

			assert.Equal(tc.statusCode, responseRecorder.Code)
		})
	}
}

func TestUpdateTodo(t *testing.T) {
	assert := assert.New(t)

	service, err := newHttpTestService()
	assert.NoError(err)

	tt := []struct {
		name       string
		method     string
		todo       api.ToDo
		statusCode int
	}{
		{
			name: "Update todo. 200 Ok",
			todo: api.ToDo{
				Message: "to do",
				ID:      1,
			},
			method:     http.MethodPost,
			statusCode: http.StatusOK,
		},
		{
			name: "Update todo with empty message. 201",
			todo: api.ToDo{
				Message: "",
				ID:      1,
			},
			method:     http.MethodPost,
			statusCode: http.StatusNoContent,
		},
		{
			name: "User id does not exist in db. 500",
			todo: api.ToDo{
				Message: "todoSmt",
				ID:      95,
			},
			method: http.MethodPost,

			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			body, err := json.Marshal(tc.todo)
			assert.NoError(err)
			request := httptest.NewRequest(tc.method, url, bytes.NewReader(body))
			responseRecorder := httptest.NewRecorder()
			params := httprouter.Params{}

			service.updateToDo(responseRecorder, request, params)

			assert.Equal(tc.statusCode, responseRecorder.Code, tc.name)
		})
	}
}

func TestGetTodo(t *testing.T) {
	assert := assert.New(t)

	service, err := newHttpTestService()
	assert.NoError(err)

	tt := []struct {
		name           string
		method         string
		params         httprouter.Params
		expectedStruct api.ToDo
		statusCode     int
	}{
		{
			name: "Get exists todo. 200 Ok",
			expectedStruct: api.ToDo{
				Message: "to do",
				ID:      1,
			},
			params: httprouter.Params{
				httprouter.Param{
					Key:   "todoid",
					Value: "1",
				}},
			method:     http.MethodGet,
			statusCode: http.StatusOK,
		},
		{
			name: "todo does not exists . 404",
			expectedStruct: api.ToDo{
				Message: "",
				ID:      0,
			},
			params: httprouter.Params{
				httprouter.Param{
					Key:   "todoid",
					Value: "99",
				}},
			method:     http.MethodGet,
			statusCode: http.StatusNotFound,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.NoError(err)
			request := httptest.NewRequest(tc.method, url, nil)
			responseRecorder := httptest.NewRecorder()

			service.getToDo(responseRecorder, request, tc.params)
			actual := api.ToDo{}
			if len(tc.expectedStruct.Message) != 0 {
				err := json.NewDecoder(responseRecorder.Body).Decode(&actual)
				assert.NoError(err)
			}

			assert.Equal(tc.expectedStruct.ID, actual.ID, tc.name)
			assert.Equal(tc.expectedStruct.Message, actual.Message, tc.name)
			assert.Equal(tc.statusCode, responseRecorder.Code, tc.name)

		})
	}
}

func TestDeleteTodo(t *testing.T) {
	assert := assert.New(t)

	service, err := newHttpTestService()
	assert.NoError(err)

	tt := []struct {
		name           string
		method         string
		params         httprouter.Params
		statusCode     int
	}{
		{
			name: "Get exists todo. 200 Ok",
			params: httprouter.Params{
				httprouter.Param{
					Key:   "todoid",
					Value: "1",
				}},
			method:     http.MethodDelete,
			statusCode: http.StatusOK,
		},
		{
			name: "todo does not exists. 200 Ok",
			params: httprouter.Params{
				httprouter.Param{
					Key:   "todoid",
					Value: "99",
				}},
			method:     http.MethodGet,
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.NoError(err)
			request := httptest.NewRequest(tc.method, url, nil)
			responseRecorder := httptest.NewRecorder()

			service.deleteToDo(responseRecorder, request, tc.params)

			assert.Equal(tc.statusCode, responseRecorder.Code, tc.name)

		})
	}
}
