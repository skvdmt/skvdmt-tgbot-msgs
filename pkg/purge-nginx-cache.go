package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
)

// PurgeNginxCache Очистка кэша Nginx.
func PurgeNginxCache(scheme, host, requestURI string) error {
	c := fmt.Sprintf("%s%s%s", scheme, host, requestURI)
	s := md5.Sum([]byte(c))
	h := hex.EncodeToString(s[:])
	filePath := fmt.Sprintf("/usr/local/nginx/cache/%c/%s/%s", h[31], h[29:31], h)
	err := os.Remove(filePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
