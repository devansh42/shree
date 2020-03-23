package shree

import (
	"crypto/md5"
	"encoding/hex"
)

type user struct {
	Email    string `json:email`
	Password string `json:password`
	UID      int64  `json:int`
}

func (u *user) setPassword(pa []byte) {
	n := md5.New()
	u.Password = hex.EncodeToString(n.Sum(pa))
}
