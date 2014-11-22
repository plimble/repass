package repass

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/plimble/errs"
	"strings"
	"time"
)

func encode(email string) string {
	hash := md5.New()
	hash.Write([]byte(time.Now().UTC().Format(time.RFC3339) + "," + email))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func decode(tokenID string) (*Token, error) {
	result, err := base64.URLEncoding.DecodeString(tokenID)

	texts := strings.Split(string(result), ",")
	if len(texts) != 2 {
		return nil, errs.NewErrors("Token is invlid")
	}

	expire, err := time.Parse(time.RFC3339, texts[0])
	if err != nil {
		return nil, err
	}

	return &Token{tokenID, texts[1], expire, false}, nil
}
