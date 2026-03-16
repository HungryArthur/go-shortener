package url

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mock_handler "github.com/HungryArthur/go-shortener/internal/mocks"
	"github.com/HungryArthur/go-shortener/internal/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestURLHandler_Create(t *testing.T) {
	type want struct {
		code                  int
		plainTextResponseBody string
		jsonResponseBody      string

		contentType *string
	}

	tests := []struct {
		name string

		want want

		shortenedURL string
		sourceURL    string
		service      func(test *testing.T) URLService
	}{
		// норм тест
		{
			name: "success",
			service: func(test *testing.T) URLService {
				controller := gomock.NewController(test)
				mock := mock_handler.NewMockURLService(controller)
				mock.EXPECT().Save("https://www.perplexity.ai/").Return("http://localhost:8080/random", nil).Times(1)
				return mock
			},
			want: want{
				code:                  201,
				plainTextResponseBody: "http://localhost:8080/random",
				jsonResponseBody: `
				{
					"result": "http://localhost:8080/random"
				}
				`,
				contentType: nil,
			},
			shortenedURL: "random",
			sourceURL:    "https://www.perplexity.ai/",
		},
		// не норм тест
		{
			name: "not success",
			service: func(test *testing.T) URLService {
				controller := gomock.NewController(test)
				mock := mock_handler.NewMockURLService(controller)
				mock.EXPECT().Save("lol").Return("", service.ErrShortenedURLDoesntExist).Times(1)
				return mock
			},
			want: want{
				code:                  500,
				plainTextResponseBody: "can't create url",
				jsonResponseBody: `
				{
					"error": "can't create url"
				}
				`,
				contentType: nil,
			},
			shortenedURL: "",
			sourceURL:    "lol",
		},
		// длинная ссылка тест
		{
			name: "long URL",
			service: func(test *testing.T) URLService {
				controller := gomock.NewController(test)
				mock := mock_handler.NewMockURLService(controller)
				return mock
			},
			want: want{
				code:                  413,
				plainTextResponseBody: "request's body too big",
				jsonResponseBody: `
				{
					"error": "request's body too big"
				}
				`,
				contentType: nil,
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
			// test for plain text handler
			w := httptest.NewRecorder()

			var body io.Reader
			if tt.sourceURL != "" {
				body = strings.NewReader(tt.sourceURL)
			} else {
				body = strings.NewReader("")
			}

			request := httptest.NewRequest(http.MethodPost, "/", body)

			h := NewURLHandler(tt.service(t))
			h.CreateTextPlain(w, request)

			res := w.Result()
			defer res.Body.Close()

			rBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			require.Equal(t, tt.want.code, res.StatusCode)
			require.Equal(t, tt.want.plainTextResponseBody, string(rBody))

			if tt.want.contentType != nil {
				require.Equal(t, *tt.want.contentType, res.Header.Get("Content-Type"))
			} else {
				require.Equal(t, "text/plain", res.Header.Get("Content-Type"))
			}

			// test for json handler
			w = httptest.NewRecorder()

			bodyBytes, _ := json.Marshal(map[string]string{"url": tt.sourceURL})

			body = bytes.NewReader(bodyBytes)

			request = httptest.NewRequest(http.MethodPost, "/", body)

			h = NewURLHandler(tt.service(t))
			h.CreateJSON(w, request)

			res = w.Result()
			defer res.Body.Close()

			rBody, err = io.ReadAll(res.Body)
			require.NoError(t, err)

			require.Equal(t, tt.want.code, res.StatusCode)

			replaceFn := func(r rune) rune {
				if r == '\t' || r == '\n' || r == ' ' {
					return -1
				} else {
					return r
				}
			}

			require.Equal(t,
				strings.Map(replaceFn, tt.want.jsonResponseBody),
				strings.Map(replaceFn, string(rBody)),
			)

			if tt.want.contentType != nil {
				require.Equal(t, *tt.want.contentType, res.Header.Get("Content-Type"))
			} else {
				require.Equal(t, "application/json", res.Header.Get("Content-Type"))
			}
		})
	}
}
