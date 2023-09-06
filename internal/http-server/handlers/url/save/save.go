package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "go_pro/internal/lib/api/response"
	"go_pro/internal/lib/logger/sl"
	"go_pro/internal/lib/random"
	"go_pro/internal/storage"
	"golang.org/x/exp/slog"
	"net/http"
)

type Response struct {
	resp.Response
	Alias string `jsom:"alias,omitempty"`
}
type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSaver
type URLSaver interface {
	SaveUrl(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("Invalid request", sl.Err(err))
			render.JSON(w, r, resp.Error("Invalid request"))
			render.JSON(w, r, resp.ValidatorError(validateErr))
			return
		}

		alias := req.Alias
		//TODO validate if alias already exists
		if alias == "" {
			alias = random.RandomAlias(resp.AliasLenght)
		}
		id, err := urlSaver.SaveUrl(req.URL, alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Info("url alredy exist")
				render.JSON(w, r, resp.Error("url already exists"))
				return
			}
			log.Error("failed to add url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add url"))
		}
		log.Info("url added", slog.Int64("id", id))
		render.JSON(w, r, Response{
			Response: resp.Ok(),
			Alias:    alias,
		})
	}

}
