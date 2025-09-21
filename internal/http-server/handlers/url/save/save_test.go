package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kelar1s/url-shortener/internal/http-server/handlers/url/save"
	"github.com/kelar1s/url-shortener/internal/http-server/handlers/url/save/mocks"
	"github.com/kelar1s/url-shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	testCases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
		},
		{
			name:  "Empty alias",
			alias: "",
			url:   "https://google.com",
		},
		{
			name:      "Empty url",
			alias:     "some_alias",
			url:       "",
			respError: "field URL is a required field",
		},
		{
			name:      "Invalid url",
			alias:     "some_alias",
			url:       "some invalid URL",
			respError: "field URL is not a valid URL",
		},
		{
			name:      "SaveURL Error",
			alias:     "test_alias",
			url:       "https://google.com",
			respError: "failed to add url",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).Return(int64(1), tc.mockError).Once()
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))

			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			body := rr.Body.String()
			var resp save.Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			if tc.respError == "" {
				require.Equal(t, "OK", resp.Status)
				if tc.alias == "" {
					require.NotEmpty(t, resp.Alias)
				} else {
					require.Equal(t, tc.alias, resp.Alias)
				}
			} else {
				require.Equal(t, "Error", resp.Status)
				require.Equal(t, tc.respError, resp.Error)
			}
		})
	}
}
