package health

import (
	"github.com/SkaisgirisMarius/article-processor/helper"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func InitHealthRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getHealthHandler())
	return r
}

func getHealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		helper.SendJsonOk(w, "Service is running")
	}
}
