package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
)

func NormalizeURL(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}

	path := u.Path
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	return fmt.Sprintf("%s%s", u.Host, u.Path), nil
}

func ToMd5Hex(s string) string {
	md5sum := md5.Sum([]byte(s))
	return hex.EncodeToString(md5sum[:])
}
