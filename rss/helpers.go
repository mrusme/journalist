package rss

import (
  "crypto/md5"
  "encoding/hex"
)

func GetGUID(src string) (string) {
  hash := md5.Sum([]byte(src))
  return hex.EncodeToString(hash[:])
}
