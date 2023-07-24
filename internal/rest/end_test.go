package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/barpav/msg-sessions/internal/rest/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_endSessions(t *testing.T) {
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
		wantStatus  int
	}{
		{
			name: "One session ended",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("DELETE", "/", nil)
					r.URL.RawQuery = "id=5"
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("EndSession", mock.Anything, mock.Anything, mock.Anything).Return(nil)
					return s
				}(),
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "Incorrect id",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("DELETE", "/", nil)
					r.URL.RawQuery = "id=iAmIncorrect"
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("EndSession", mock.Anything, mock.Anything, mock.Anything).Return(nil)
					return s
				}(),
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "All sessions ended",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("DELETE", "/", nil),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("EndAllSessions", mock.Anything, mock.Anything).Return(nil)
					return s
				}(),
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "Error",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("DELETE", "/", nil)
					r.Header.Set("request-id", "test-request-id")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("EndAllSessions", mock.Anything, mock.Anything).Return(errors.New("test error"))
					return s
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
			s.endSessions(tt.args.w, tt.args.r)

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

			require.Equal(t, tt.wantStatus, tt.args.w.Code)
		})
	}
}
