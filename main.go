package main

import (
	"flag"
	"runtime"

	"github.com/guregu/kami"
	"golang.org/x/net/context"

	"github.com/satoshi03/go-dsp-api/common/consts"
	"github.com/satoshi03/go/config"
	"github.com/satoshi03/go/fluent"
	"github.com/satoshi03/go/redis"

	"github.com/satoshi03/related-article-api/common"
)

func main() {
	flag.Parse()
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	config := config.Read(`config.yml`)

	ctx := context.Background()

	// init redis connection
	ctx = redis.Open(ctx, config.Redis[common.CtxRedisKey], common.CtxRedisKey)
	defer redis.Close(ctx, common.CtxRedisKey)

	ctx = fluent.Open(ctx, config.Fluent, common.CtxFluentKey)
	defer fluent.Close(ctx, consts.CtxFluentKey)

	kami.Context = ctx
	kami.Get("/v1/ra/json", articleJsonHandler)
	kami.Get("/v1/ra/jsonp", articleJsonpHandler)
	kami.Serve()
}
