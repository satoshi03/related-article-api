package main

import (
	"net/http"

	"golang.org/x/net/context"

	"github.com/satoshi03/related-article-api/article"
	"github.com/satoshi03/related-article-api/common"
	"github.com/satoshi03/related-article-api/utils"
)

func articleHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// Parse Request
	siteID := r.FormValue("site_id")
	articleID := r.FormValue("article_id")

	// Validate Request
	if siteID == "" {
		// Do error process
		utils.Write404Response(w, map[string]interface{}{"message": "siteID not found"})
		return
	}

	// Get Related Artcile
	index := article.GetIndexRelated(ctx, siteID, articleID)
	if len(*index) < common.MinArticleLength {
		index = article.GetIndexRanking(ctx, siteID)
	}

	// Get Artcile Info
	articles := article.GetArticleInfo(ctx, *index, siteID)

	// Make Response
	resp := makeResponse(articles)

	// Return Response
	utils.WriteResponse(w, resp, 200)
}

func makeResponse(articles []article.Article) map[string]interface{} {
	ais := make([]map[string]interface{}, 0, len(articles))
	for _, ar := range articles {
		ai := map[string]interface{}{
			"title":     ar.Title,
			"url":       ar.URL,
			"icon_url":  ar.IconURL,
			"image_url": ar.ImageURL,
		}
		ais = append(ais, ai)
	}
	return map[string]interface{}{
		"articles": ais,
	}
}
