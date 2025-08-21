package server

import (
	"net/http"
	"teslatrack/internal/biz"
	"teslatrack/internal/conf"
)

type Redirector struct {
	authorizeTokenUsecase *biz.AuthorizeTokenUsecase
	conf                  *conf.Server
}

func (redirect *Redirector) RedirectFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == redirect.conf.Tesla.Callback {
			// tesla callback auth code
			code := r.URL.Query().Get("code")
			redirect.authorizeTokenUsecase.ExchangeCode(r.Context(), code)
			http.ServeFile(w, r, "web/hello.html")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func NewRedirector(conf *conf.Server, authorizeTokenUsecase *biz.AuthorizeTokenUsecase) *Redirector {
	return &Redirector{conf: conf, authorizeTokenUsecase: authorizeTokenUsecase}
}
