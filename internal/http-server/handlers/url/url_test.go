package url

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wan6sta/url-shortener/internal/storage/postgres"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const googleUrl = "https://google.com/"

func TestUrlHandler(t *testing.T) {
	st, err := postgres.NewStorage()
	if err != nil {
		t.Errorf("cannot create storage: %v", err)
	}
	u := NewUrl(st)

	var resStr string

	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name   string
		method string
		want   want
	}{
		{
			name:   "[POST] positive test #1",
			method: http.MethodPost,
			want: want{
				code:        http.StatusCreated,
				response:    localhost,
				contentType: "text/plain",
			},
		},
		{
			name:   "[GET] positive test #1",
			method: http.MethodGet,
			want: want{
				code:        http.StatusTemporaryRedirect,
				response:    resStr,
				contentType: "text/plain",
			},
		},
		{
			name:   "[GET] negative test #1",
			method: http.MethodGet,
			want: want{
				code:        http.StatusBadRequest,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	test1 := tests[0]
	t.Run(test1.name, func(t *testing.T) {
		request := httptest.NewRequest(test1.method, "/", bytes.NewBufferString(googleUrl))
		w := httptest.NewRecorder()
		u.UrlHandler(w, request)

		res := w.Result()

		assert.Equal(t, res.StatusCode, test1.want.code)
		defer res.Body.Close()

		resBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		contains := strings.Contains(string(resBody), localhost)

		resStr = string(resBody)

		assert.True(t, contains)
		assert.Equal(t, res.Header.Get("Content-Type"), test1.want.contentType)

	})

	test2 := tests[1]
	t.Run(test2.name, func(t *testing.T) {
		ids := strings.Split(resStr, "/")
		id := ids[3]

		request := httptest.NewRequest(test2.method, fmt.Sprintf("/%s", id), nil)
		w := httptest.NewRecorder()
		u.UrlHandler(w, request)

		res := w.Result()

		assert.Equal(t, res.StatusCode, test2.want.code)
		require.NoError(t, err)

		resID := res.Header.Get("Location")

		assert.True(t, resID == googleUrl)
		assert.Equal(t, res.Header.Get("Content-Type"), test2.want.contentType)
	})

	test3 := tests[2]
	t.Run(test3.name, func(t *testing.T) {
		request := httptest.NewRequest(test3.method, fmt.Sprintf("/%s", "123123123"), nil)
		w := httptest.NewRecorder()
		u.UrlHandler(w, request)

		res := w.Result()

		assert.Equal(t, res.StatusCode, test3.want.code)
		require.NoError(t, err)
		assert.Equal(t, res.Header.Get("Content-Type"), test3.want.contentType)
	})
}
