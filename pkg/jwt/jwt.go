package jwt

import "time"

// LoginUser
type LoginUser struct {
	ID        int64     `json:"id"`
	Mobile    string    `json:"mobile"`
	OpenID    string    `json:"open_id"`
	NickName  string    `json:"nick_name"`
	Avatar    string    `json:"avatar"`
	LoginTime time.Time `json:"login_time"`
}
