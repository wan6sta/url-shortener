package url

import (
	"github.com/wan6sta/url-shortener/internal/storage/postgres"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

func CreateUrlHandler(log *slog.Logger, st *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			const op = "handlers.CreateUrl"

			res, err := io.ReadAll(r.Body)
			if err != nil {
				log.Error("cannot parse body", op, err.Error())
				return
			}

			url, err := st.CreateUrl(string(res))
			if err != nil {
				log.Error("key does not exists", op, err.Error())
				return
			}

			w.WriteHeader(http.StatusCreated)

			_, err = w.Write([]byte(url))
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

			url, err := st.GetUrl(id)
			if err != nil {
				log.Error("cannot write response", op, err.Error())
				return
			}

			w.Header().Set("Location", url)
			w.WriteHeader(http.StatusTemporaryRedirect)

			/*_, err = w.Write([]byte(url))
			if err != nil {
				log.Error("cannot write response", op, err.Error())
				return
			}*/

			return
		}

		http.Error(w, "method not allowed", http.StatusBadRequest)
	}
}
