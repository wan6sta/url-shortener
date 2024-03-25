package url

import (
	"github.com/wan6sta/url-shortener/internal/storage/postgres"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

const localhost = "http://localhost:8080/"

func HandleUrl(log *slog.Logger, s *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			const op = "handlers.CreateUrl"

			res, err := io.ReadAll(r.Body)
			if err != nil {
				log.Error("cannot parse body", op, err.Error())
				return
			}

			url, err := postgres.CreateUrl(string(res), s)
			if err != nil {
				log.Error("key does not exists", op, err.Error())
				return
			}

			w.WriteHeader(http.StatusCreated)

			_, err = w.Write([]byte(localhost + url))
			if err != nil {
				log.Error("cannot write response", op, err.Error())
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

			url, err := postgres.GetUrl(id, s)
			if err != nil {
				log.Error("cannot write response", op, err.Error())
				return
			}

			w.Header().Set("Location", url)
			w.WriteHeader(http.StatusTemporaryRedirect)

			return
		}

		http.Error(w, "method not allowed", http.StatusBadRequest)
	}
}
