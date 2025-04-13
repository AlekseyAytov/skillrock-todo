package app

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AlekseyAytov/skillrock-todo/internal/models/master"
	"github.com/AlekseyAytov/skillrock-todo/internal/store/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var todo *ToDoAPI

func TestMain(m *testing.M) {
	storage := mock.NewMockStore()
	todo = NewToDoAPI(master.NewTaskMaster(storage))
	m.Run()
}

func TestToDoAPI_addTask(t *testing.T) {
	type want struct {
		code int
	}

	tests := []struct {
		name          string
		requestPath   string
		requestMethod string
		requestBody   string
		want          want
	}{
		{
			name:          "first, positive",
			requestPath:   "/tasks",
			requestMethod: http.MethodPost,
			requestBody:   `{"status": "done","title": "buy", "description": "milk"}`,
			want: want{
				code: http.StatusCreated,
			},
		},
		{
			name:          "second, negative",
			requestPath:   "/tasks",
			requestMethod: http.MethodPost,
			requestBody:   `{"status": "in_progres","title": "", "description": "learn 10 words"}`,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:          "third, positive",
			requestPath:   "/tasks",
			requestMethod: http.MethodPost,
			requestBody:   `{"status": "in_progres","title": "run 100 km"}`,
			want: want{
				code: http.StatusCreated,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Подготовка входных данных для запроса
			reqBody := strings.NewReader(tt.requestBody)
			request := httptest.NewRequest(tt.requestMethod, tt.requestPath, reqBody)
			request.Header.Set("Content-Type", "application/json")

			// Обработка запроса
			res, err := todo.api.Test(request)
			require.NoError(t, err)

			// проверка статус кода
			assert.Equal(t, tt.want.code, res.StatusCode)

			// проверка тела ответа
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.NotEmpty(t, resBody, "response body must contains some information")
		})
	}
}
