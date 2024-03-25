package url

import (
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type Creater interface {
	CreateUrl(url string) (string, error)
	GetUrl(url string) (string, error)
}

func CreateUrl(log *slog.Logger, creater Creater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			const op = "handlers.CreateUrl"

			res, err := io.ReadAll(r.Body)
			if err != nil {
				log.Error("cannot parse body", op, err.Error())
				return
			}

			url, err := creater.CreateUrl(string(res))
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

			url, err := creater.GetUrl(id)
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
