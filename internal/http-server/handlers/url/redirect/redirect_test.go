package redirect_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/kelar1s/url-shortener/internal/http-server/handlers/url/redirect"
	"github.com/kelar1s/url-shortener/internal/http-server/handlers/url/redirect/mocks"
	"github.com/kelar1s/url-shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/kelar1s/url-shortener/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestRedirectHandler(t *testing.T) {
	testCases := []struct {
		name           string
		alias          string
		mockURL        string
		mockError      error
		expectedStatus int
		expectedLoc    string
	}{
		{
			name:           "Success",
			alias:          "github",
			mockURL:        "https://github.com",
			expectedStatus: http.StatusFound,
			expectedLoc:    "https://github.com",
		},
		{
			name:           "URL not found",
			alias:          "missing",
			mockError:      storage.ErrURLNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Internal error",
			alias:          "oops",
			mockError:      errors.New("db failure"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Empty alias",
			alias:          "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			getterMock := mocks.NewURLGetter(t)
			if tc.alias != "" {
				getterMock.
					On("GetURL", tc.alias).
					Return(tc.mockURL, tc.mockError).
					Once()
			}

			handler := redirect.New(slogdiscard.NewDiscardLogger(), getterMock)

			req := httptest.NewRequest(http.MethodGet, "/"+tc.alias, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("alias", tc.alias)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.expectedStatus, rr.Code)

			if tc.expectedLoc != "" {
				require.Equal(t, tc.expectedLoc, rr.Header().Get("Location"))
			}
		})
	}
}
