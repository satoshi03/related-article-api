package main

import (
	"github.com/satoshi03/related-article-api/article"
	"github.com/satoshi03/related-article-api/page"
)

func init() {
	article.InitHandler()
	page.InitHandler()
}
