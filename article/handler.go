package article

import (
	"fmt"
	"net/http"
	"strconv"

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
	// TODO: 記事取得に使うオプションたちは構造体にまとめる
	articles, userGroupID := getArticles(ctx, siteID, articleID, referer, cookieUserID)

	// Make Response
	resp := makeResponse(articles, cookieUserID)

	// Log
	sendLog(ctx, articles, referer, cookieUserID, userGroupID)

	// Return Response
	fun(w, resp, 200)
}

func getIndex(ctx context.Context, siteID, referer, cookieUserID string) (*Index, int) {
	lrl := GetLogicRatioList(ctx, siteID, referer)
	cui, err := strconv.Atoi(cookieUserID)

	// cookieUserIDが数字でない場合はランキングを出す
	if err != nil {
		return GetIndexRanking(ctx, siteID), -1
	}

	// cookieUserID % 100 が ratio 未満の場合はそのユーザグループに所属する
	for _, lr := range *lrl {
		if cui%100 < lr.Ratio {
			return GetIndexRelated(ctx, siteID, referer, lr.UserGroupID), lr.UserGroupID
		}
	}

	return GetIndexRanking(ctx, siteID), -1
}

func getArticles(ctx context.Context, siteID, articleID, referer, cookieUserID string) ([]Article, int) {
	// Get Related Artcile
	index, ugid := getIndex(ctx, siteID, referer, cookieUserID)
	if len(*index) < common.MinArticleLength {
		index = GetIndexRanking(ctx, siteID)
		ugid = -1
	}
	// Get Artcile Info
	return GetArticleInfo(ctx, *index, siteID), ugid
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

func sendLog(ctx context.Context, articles []Article, referer, cookieUserID string, userGroupID int) {
	for i, ar := range articles {
		ai := map[string]interface{}{
			"article_id":     ar.ID,
			"index":          i,
			"referer":        referer,
			"cookie_user_id": cookieUserID,
			"user_group_id":  userGroupID,
		}
		fluent.Send(ctx, common.CtxFluentKey, "article.get", ai)
	}
}

// backendでやったほうがいいかも
func makeRedirectURL(a Article, cuid string) string {
	return fmt.Sprintf("%s/v1/page?site_id=%d&redirect_to=%s&cuid=%s&aid=%d", common.BASE_URL, a.SiteID, a.URL, cuid, a.ID)
}

func InitHandler() {
	kami.Get("/v1/ra/json", articleJsonHandler)
	kami.Get("/v1/ra/jsonp", articleJsonpHandler)
}
