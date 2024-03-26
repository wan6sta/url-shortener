package url

import (
	"io"
	"log/slog"
	"net/http"
	"strings"
)

const localhost = "http://localhost:8080/"

type Repositories interface {
	CreateUrl(url string) (string, error)
	GetUrl(url string) (string, error)
}

type Url struct {
	r Repositories
}

func NewUrl(r Repositories) *Url {
	return &Url{r: r}
}

func (u *Url) UrlHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("charset", "utf8")

	if r.Method == http.MethodPost {
		const op = "handlers.CreateUrl"

		res, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("cannot parse body", op, err.Error())
			http.Error(w, "cannot parse body", http.StatusBadRequest)

			return
		}

		url, err := u.r.CreateUrl(string(res))
		if err != nil {
			slog.Error("key does not exists", op, err.Error())
			http.Error(w, "key does not exists", http.StatusBadRequest)

			return
		}

		w.WriteHeader(http.StatusCreated)

		_, err = w.Write([]byte(localhost + url))
		if err != nil {
			slog.Error("cannot write response", op, err.Error())
			http.Error(w, "cannot write response", http.StatusBadRequest)

			return
		}
		return
	}

	if r.Method == http.MethodGet {
		const op = "handlers.GetUrl"

		id := strings.TrimPrefix(r.URL.Path, "/")
		if id == "" {
			http.Error(w, "ID not provided", http.StatusBadRequest)
			return
		}

		url, err := u.r.GetUrl(id)
		if err != nil {
			slog.Error("cannot write response", op, err.Error())
			http.Error(w, "key does not exists", http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)

		return
	}
}
