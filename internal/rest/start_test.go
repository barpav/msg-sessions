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

func TestService_startNewSession(t *testing.T) {
	type testService struct {
		storage Storage
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name        string
		args        args
		testService testService
		wantHeaders map[string]string
		wantStatus  int
	}{
		{
			name: "Session started",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/", nil),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("StartNewSession", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
						int64(5), "test-session-key", nil,
					)
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"session-id":  "5",
				"session-key": "test-session-key",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Too many sessions",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/", nil),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("StartNewSession", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
						int64(0), "", &ErrTooManySessionsTest{},
					)
					return s
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusTooManyRequests,
		},
		{
			name: "Session creation error",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("POST", "/", nil)
					r.Header.Set("request-id", "test-request-id")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("StartNewSession", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
						int64(0), "", errors.New("test error"),
					)
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
			s.startNewSession(tt.args.w, tt.args.r)

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

type ErrTooManySessionsTest struct{}

func (e *ErrTooManySessionsTest) Error() string {
	return "ErrTooManySessionsTest"
}

func (e *ErrTooManySessionsTest) ImplementsTooManySessionsError() {
}
