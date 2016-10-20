package article

import (
	"fmt"
	"net/http"

	"github.com/guregu/kami"
	"golang.org/x/net/context"

	"github.com/satoshi03/go/fluent"
	"github.com/satoshi03/related-article-api/common"
)

func articleJsonHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	articleHandler(ctx, w, r, common.JsonWriter)
}

func articleJsonpHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	articleHandler(ctx, w, r, common.JsonpWriter)
}

func articleHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, fun common.ResponseWriter) {
	// Parse Request
	siteID := r.FormValue("site_id")
	articleID := r.FormValue("article_id")
	cookieUserID := r.FormValue("cuid")

	// Validate Request
	if siteID == "" {
		// Do error process
		common.Write404Response(w, map[string]interface{}{"message": "siteID not found"})
		return
	}

	// referからURLを取得
	referer := r.Referer()
	if referer == "" {
		common.Write404Response(w, map[string]interface{}{"message": "referer url not found"})
		return
	}
	referer, _ = common.NormalizeURL(referer)

	// Get Articles related with designated article
	articles := getArticles(ctx, siteID, articleID, referer)

	// Make Response
	resp := makeResponse(articles, cookieUserID)

	// Log
	sendLog(ctx, articles, referer, cookieUserID)

	// Return Response
	fun(w, resp, 200)
}

func getArticles(ctx context.Context, siteID, articleID, referer string) []Article {
	// Get Related Artcile
	index := GetIndexRelated(ctx, siteID, referer)
	if len(*index) < common.MinArticleLength {
		index = GetIndexRanking(ctx, siteID)
	}
	// Get Artcile Info
	return GetArticleInfo(ctx, *index, siteID)
}

func makeResponse(articles []Article, cuid string) map[string]interface{} {
	ais := make([]map[string]interface{}, 0, len(articles))
	for i, ar := range articles {
		ai := map[string]interface{}{
			"title":     ar.Title,
			"url":       makeRedirectURL(ar, cuid),
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

func sendLog(ctx context.Context, articles []Article, referer, cookieUserID string) {
	for i, ar := range articles {
		ai := map[string]interface{}{
			"article_id":     ar.ID,
			"index":          i,
			"referer":        referer,
			"cookie_user_id": cookieUserID,
		}
		fluent.Send(ctx, common.CtxFluentKey, "article.get", ai)
	}
}

// backendでやったほうがいいかも
func makeRedirectURL(a Article, cuid string) string {
	return fmt.Sprintf("%s/v1/page?site_id=%d&redirect_to=%s&cuid=%s", common.BASE_URL, a.SiteID, a.URL, cuid)
}

func InitHandler() {
	kami.Get("/v1/ra/json", articleJsonHandler)
	kami.Get("/v1/ra/jsonp", articleJsonpHandler)
}
