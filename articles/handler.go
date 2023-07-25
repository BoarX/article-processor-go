package articles

import (
	"github.com/SkaisgirisMarius/article-processor/helper"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// InitArticlesRouter initializes the articles router using the chi package, sets up two routes for handling article requests
func InitArticlesRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/{id}", getArticleByIDHandler)
	r.Get("/list", getArticleListHandler)
	return r
}

// getArticleByIDHandler is an HTTP handler function that handles requests to get a single article by its ID.
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

// getArticleListHandler is an HTTP handler function that handles requests to get a list of articles.
// It retrieves the articles from the database using getArticleListFromDatabase and sends a JSON response containing the list of articles
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
