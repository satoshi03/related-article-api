package article

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/satoshi03/go/redis"
	"github.com/satoshi03/related-article-api/common"
)

type LogicRatioList []LogicRatio

type LogicRatio struct {
	UserGroupID int `msgpack:"user_group_id"`
	RecType     int `msgpack:"rec_type"`
	Ratio       int `msgpack:"ratio"`
}

func (l *LogicRatioList) makeKey(siteID, articleURL string) string {
	return fmt.Sprintf("logic_ratio:%s:%s", siteID, articleURL)
}

func GetLogicRatioList(ctx context.Context, siteID, articleURL string) *LogicRatioList {
	var lrl LogicRatioList
	if articleURL == "" {
		return &lrl
	}
	articleURLHex := common.ToMd5Hex(articleURL)
	key := lrl.makeKey(siteID, articleURLHex)
	out, err := redis.GetPackedValue(ctx, common.CtxRedisKey, key, &lrl)
	if err != nil {
		return &lrl
	}
	return out.(*LogicRatioList)
}
