package rest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/barpav/msg-sessions/internal/rest/mocks"
	"github.com/barpav/msg-sessions/internal/rest/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_getActiveSessions(t *testing.T) {
	type testService struct {
		storage Storage
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name        string
		testService testService
		args        args
		wantHeaders map[string]string
		wantBody    *models.UserSessionsV1
		wantStatus  int
	}{
		{
			name: "Invalid MIME-type (406)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/", nil)
					r.Header.Set("Accept", "*/*")
					return r
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusNotAcceptable,
		},
		{
			name: "OK (200)",
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("GetSessionsV1", mock.Anything, mock.Anything).Return(
						&models.UserSessionsV1{
							Active: 0,
							List:   make([]*models.UserSessionV1, 0),
						},
						nil,
					)
					return s
				}(),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/", nil)
					r.Header.Set("Accept", "application/vnd.userSessions.v1+json")
					return r
				}(),
			},
			wantHeaders: map[string]string{},
			wantBody: &models.UserSessionsV1{
				Active: 0,
				List:   make([]*models.UserSessionV1, 0),
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "OK without Accept (200)",
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("GetSessionsV1", mock.Anything, mock.Anything).Return(
						&models.UserSessionsV1{
							Active: 0,
							List:   make([]*models.UserSessionV1, 0),
						},
						nil,
					)
					return s
				}(),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/", nil),
			},
			wantHeaders: map[string]string{},
			wantBody: &models.UserSessionsV1{
				Active: 0,
				List:   make([]*models.UserSessionV1, 0),
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Error (500)",
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("GetSessionsV1", mock.Anything, mock.Anything).Return(
						nil,
						errors.New("Test error"),
					)
					return s
				}(),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/", nil)
					r.Header.Set("Accept", "application/vnd.userSessions.v1+json")
					r.Header.Set("request-id", "test-request-id")
					return r
				}(),
			},
			wantHeaders: map[string]string{
				"issue": "test-request-id",
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.testService.storage,
			}
			s.getActiveSessions(tt.args.w, tt.args.r)

			for k, v := range tt.wantHeaders {
				require.Equal(t, v, func() string {
					h := tt.args.w.Result().Header
					if h == nil {
						return ""
					}
					v := h[k]
					if len(v) == 0 {
						return ""
					}
					return v[0]
				}())
			}

			var body *models.UserSessionsV1
			decoded := models.UserSessionsV1{}
			err := json.NewDecoder(tt.args.w.Body).Decode(&decoded)

			if err != nil && err != io.EOF { // empty body is ok
				t.Fatal(err)
			}

			if err == nil {
				body = &decoded
			}

			require.Equal(t, body, tt.wantBody)
			require.Equal(t, tt.wantStatus, tt.args.w.Code)
		})
	}
}
