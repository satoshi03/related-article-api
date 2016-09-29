package page

import (
	"fmt"
	"net/http"

	"github.com/guregu/kami"
	"github.com/satoshi03/related-article-api/common"
	"golang.org/x/net/context"
)

func pageHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	redirect_to := r.FormValue("redirect_to")
	if redirect_to == "" {
		common.Write404Response(w, map[string]interface{}{"message": "redirect url not found"})
	}

	// Write Cookie

	// Incr Click log

	// Redirect to designated url
	http.Redirect(w, r, redirect_to, http.StatusFound)

	body := fmt.Sprintf("redirect to %s soon...", redirect_to)
	w.Write([]byte(body))
}

func InitHandler() {
	kami.Get("/v1/page", pageHandler)
}
