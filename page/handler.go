package page

import (
	"fmt"
	"net/http"

	"github.com/guregu/kami"
	"github.com/satoshi03/go/redis"
	"golang.org/x/net/context"

	"github.com/satoshi03/related-article-api/common"
)

func pageHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	siteID := r.FormValue("site_id")
	if siteID == "" {
		common.Write404Response(w, map[string]interface{}{"message": "site id not found"})
		return
	}
	redirectTo := r.FormValue("redirect_to")
	if redirectTo == "" {
		common.Write404Response(w, map[string]interface{}{"message": "redirect url not found"})
		return
	}

	// Write Cookie

	// Incr Click log
	nurl, err := common.NormalizeURL(redirectTo)
	if err != nil {
		common.Write404Response(w, map[string]interface{}{"message": "redirect url is not valid"})
		return
	}
	str := common.ToMd5Hex(nurl)
	redis.Incr(ctx, common.CtxRedisKey, fmt.Sprintf("click:%s:%s", siteID, str))

	// Redirect to designated url
	http.Redirect(w, r, redirectTo, http.StatusFound)

	body := fmt.Sprintf("redirect to %s soon...", redirectTo)
	w.Write([]byte(body))
}

func InitHandler() {
	kami.Get("/v1/page", pageHandler)
}