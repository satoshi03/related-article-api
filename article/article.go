package article

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/satoshi03/go/redis"

	"github.com/satoshi03/related-article-api/common"
)

type Article struct {
	SiteID   int    `msgpack:"site_id"`
	ID       int    `msgpack:"article_id"`
	URL      string `msgpack:"url"`
	Title    string `msgpack:"title"`
	IconURL  string `msgpack:"icon_url"`
	ImageURL string `msgpack:"image_url"`
}

func (a *Article) makeKey(siteID string, articleID int) string {
	return fmt.Sprintf("info:%s:%d", siteID, articleID)
}

func GetArticleInfo(ctx context.Context, index Index, siteID string) []Article {
	var key string
	articles := make([]Article, 0, len(index))
	for i, ai := range index {
		// XXX: redisへの問い合わせ回数が増えるのでmgetにしたほうがいい
		//      redigoだとあんまいい感じできないっぽい
		// 最大数取得したらループを抜ける
		// TODO: サイトごとに設定できるようにしたほうがいいかも
		if i >= common.MaxArticleLength {
			break
		}
		var article Article
		key = article.makeKey(siteID, ai.ID)
		redis.GetPackedValue(ctx, common.CtxRedisKey, key, &article)
		articles = append(articles, article)
	}
	return articles
}
