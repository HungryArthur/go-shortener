package url

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mock_handler "github.com/HungryArthur/go-shortener/internal/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestURLHandler_GetTextPlain(t *testing.T) {
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
		service func(test *testing.T) URLService
	}{
		{
			name: "success get",
			service: func(test *testing.T) URLService {
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

			h := NewURLHandler(tt.service(t))
			h.GetTextPlain(w, r)

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
