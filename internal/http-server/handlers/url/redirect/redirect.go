package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	resp "github.com/kelar1s/url-shortener/internal/lib/api/response"
	"github.com/kelar1s/url-shortener/internal/lib/logger/sl"
	"github.com/kelar1s/url-shortener/internal/storage"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

//go:generate mockery --name=URLGetter --with-expecter=true --case=underscore --output=./mocks --outpkg=mocks
func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("alias is empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("url not found", slog.String("alias", alias))

				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, resp.Error("url not found"))

				return
			}
			log.Error("failed get url", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("got url", slog.String("url", resURL))

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
