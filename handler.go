package main

import (
	"net/http"

	"golang.org/x/net/context"

	"github.com/satoshi03/related-article-api/article"
	"github.com/satoshi03/related-article-api/common"
	"github.com/satoshi03/related-article-api/utils"
)

type responseWriter func(W http.ResponseWriter, resp map[string]interface{}, statusCode int)

func jsonpWriter(w http.ResponseWriter, resp map[string]interface{}, statusCode int) {
	utils.WriteJsonpResponse(w, resp, statusCode)
}

func jsonWriter(w http.ResponseWriter, resp map[string]interface{}, statusCode int) {
	utils.WriteResponse(w, resp, statusCode)
}

func articleJsonHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	articleHandler(ctx, w, r, jsonWriter)
}

func articleJsonpHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	articleHandler(ctx, w, r, jsonpWriter)
}

func articleHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, fun responseWriter) {
	// Parse Request
	siteID := r.FormValue("site_id")
	articleID := r.FormValue("article_id")

	// Validate Request
	if siteID == "" {
		// Do error process
		utils.Write404Response(w, map[string]interface{}{"message": "siteID not found"})
		return
	}

	// Get Articles related with designated article
	articles := getArticles(ctx, siteID, articleID)

	// Make Response
	resp := makeResponse(articles)

	// Return Response
	fun(w, resp, 200)
}

func getArticles(ctx context.Context, siteID, articleID string) []article.Article {
	// Get Related Artcile
	index := article.GetIndexRelated(ctx, siteID, articleID)
	if len(*index) < common.MinArticleLength {
		index = article.GetIndexRanking(ctx, siteID)
	}
	// Get Artcile Info
	return article.GetArticleInfo(ctx, *index, siteID)
}

func makeResponse(articles []article.Article) map[string]interface{} {
	ais := make([]map[string]interface{}, 0, len(articles))
	for i, ar := range articles {
		ai := map[string]interface{}{
			"title":     ar.Title,
			"url":       ar.URL,
			"icon_url":  ar.IconURL,
			"image_url": ar.ImageURL,
			"index":     i,
		}
		ais = append(ais, ai)
	}
	return map[string]interface{}{
		"articles": ais,
	}
}
