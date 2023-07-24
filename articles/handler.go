package articles

import (
	"github.com/SkaisgirisMarius/article-processor/helper"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func InitArticlesRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/{id}", getArticleByIDHandler)
	r.Get("/list", getArticleListHandler)
	return r
}

func getArticleByIDHandler(w http.ResponseWriter, r *http.Request) {
	articleID := chi.URLParam(r, "id")
	if articleID == "" {
		helper.SendJsonError(w, http.StatusBadRequest, "invalid request data articleID")
		return
	}

	article, err := getArticleByIDFromDatabase(articleID)
	if err != nil {
		helper.SendJsonError(w, http.StatusInternalServerError, "invalid request data articleID")
		return
	}
	var response = SingleArticleResponse{
		Status: statusSuccess,
		Data:   article,
	}
	helper.SendJsonOk(w, response)

}

func getArticleListHandler(w http.ResponseWriter, r *http.Request) {
	articles, err := getArticleListFromDatabase()
	if err != nil {
		helper.SendJsonError(w, http.StatusInternalServerError, err)
		return
	}
	var response = MultipleArticlesResponse{
		Status: statusSuccess,
		Data:   articles,
	}

	helper.SendJson(w, http.StatusOK, response)
}
