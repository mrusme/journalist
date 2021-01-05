package j

import (
  "crypto/md5"
  "encoding/hex"
)

func GetApiKey(username string, password string) (string) {
  hash := md5.Sum([]byte(username + ":" + password))
  return hex.EncodeToString(hash[:])
}
