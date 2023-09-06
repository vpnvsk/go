package redirect

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "go_pro/internal/lib/api/response"
	"go_pro/internal/lib/logger/sl"
	"go_pro/internal/storage"
	"golang.org/x/exp/slog"
	"net/http"
)

type URLGetter interface {
	GetUrl(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.redirect"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("the alias is empty")
			render.JSON(w, r, resp.Error("alias is empty"))
			return

		}
		resUrl, err := urlGetter.GetUrl(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("not found", "alias", alias)
				render.JSON(w, r, resp.Error("not found"))
				return
			}
			log.Error("failed to get url", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))

			return
		}
		log.Info("get url", slog.String("url", resUrl))
		http.Redirect(w, r, resUrl, http.StatusFound)
	}

}
