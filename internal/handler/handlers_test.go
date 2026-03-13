package handler_test

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HungryArthur/go-shortener/internal/handler"
	mock_handler "github.com/HungryArthur/go-shortener/internal/mocks"
	"github.com/HungryArthur/go-shortener/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestURLHandler_Create(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name         string
		want         want
		shortenedURL string
		sourceURL    string
		service      func(test *testing.T) handler.URLService
	}{
		// норм тест
		{
			name: "success",
			service: func(test *testing.T) handler.URLService {
				controller := gomock.NewController(test)
				mock := mock_handler.NewMockURLService(controller)
				mock.EXPECT().Save("https://www.perplexity.ai/").Return("random", nil).Times(1)
				return mock
			},
			want: want{
				code:        201,
				response:    "http://localhost:8080/random",
				contentType: "text/plain",
			},
			shortenedURL: "random",
			sourceURL:    "https://www.perplexity.ai/",
		},
		// не норм тест
		{
			name: "not success",
			service: func(test *testing.T) handler.URLService {
				controller := gomock.NewController(test)
				mock := mock_handler.NewMockURLService(controller)
				mock.EXPECT().Save("lol").Return("", service.ErrShortenedURLDoesntExist).Times(1)
				return mock
			},
			want: want{
				code:        500,
				response:    "can't create url",
				contentType: "text/plain",
			},
			shortenedURL: "",
			sourceURL:    "lol",
		},
		// длинная ссылка тест
		{
			name: "long URL",
			service: func(test *testing.T) handler.URLService {
				controller := gomock.NewController(test)
				mock := mock_handler.NewMockURLService(controller)
				return mock
			},
			want: want{
				code:        413,
				response:    "request's body too big",
				contentType: "text/plain",
			},
			shortenedURL: "long",
			sourceURL: func() string {
				buf := make([]byte, 10000)
				rand.Read(buf)
				return base64.RawStdEncoding.EncodeToString(buf)
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			var body io.Reader
			if tt.sourceURL != "" {
				body = strings.NewReader(tt.sourceURL)
			} else {
				body = strings.NewReader("")
			}

			request := httptest.NewRequest(http.MethodPost, "/", body)

			h := handler.NewURLHandler(tt.service(t))
			h.Create(w, request)

			res := w.Result()
			defer res.Body.Close()

			rBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			require.Equal(t, tt.want.code, res.StatusCode)
			require.Equal(t, tt.want.response, string(rBody))
			require.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestURLHandler_Get(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
		location    string
	}

	tests := []struct {
		name    string
		urlPath string
		want    want
		service func(test *testing.T) handler.URLService
	}{
		{
			name: "success get",
			service: func(test *testing.T) handler.URLService {
				controller := gomock.NewController(test)
				mock := mock_handler.NewMockURLService(controller)
				mock.EXPECT().Get("random").Return("https://www.perplexity.ai/", nil).Times(1)
				return mock
			},
			want: want{
				code:        307,
				response:    "",
				contentType: "",
				location:    "https://www.perplexity.ai/",
			},
			urlPath: "/random",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("shortenedURL", strings.TrimPrefix(tt.urlPath, "/"))

			r := httptest.NewRequest(http.MethodGet, tt.urlPath, nil)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

			w := httptest.NewRecorder()

			h := handler.NewURLHandler(tt.service(t))
			h.Get(w, r)

			res := w.Result()
			defer res.Body.Close()

			rBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			require.Equal(t, tt.want.code, res.StatusCode)
			require.Equal(t, tt.want.response, string(rBody))
			require.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
