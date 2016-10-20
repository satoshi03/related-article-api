package article

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/satoshi03/go/redis"

	"github.com/satoshi03/related-article-api/common"
)

type Index []Element

type Element struct {
	ID         int     `msgpack:"aid"`
	Similarity float64 `msgpack:"sim"`
}

func (i *Index) makeKey(siteID, articleID string) string {
	return fmt.Sprintf("index:%s:%s:1", siteID, articleID)
}

func GetIndexRelated(ctx context.Context, siteID, articleURL string) *Index {
	var index Index
	if articleURL == "" {
		return &index
	}
	articleURLHex := common.ToMd5Hex(articleURL)
	key := index.makeKey(siteID, articleURLHex)
	out, err := redis.GetPackedValue(ctx, common.CtxRedisKey, key, &index)
	if err != nil {
		return &index
	}
	return out.(*Index)
}

func GetIndexRanking(ctx context.Context, siteID string) *Index {
	var index Index
	key := fmt.Sprintf("ranking:%s", siteID)
	out, err := redis.GetPackedValue(ctx, common.CtxRedisKey, key, &index)
	if err != nil {
		return &index
	}
	return out.(*Index)
}
