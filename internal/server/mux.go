package server

import (
	"net/http"
	v1 "teslatrack/api/helloworld/v1"
	"teslatrack/internal/service"
)

type Redirector struct {
	greeter *service.GreeterService
}

func (redirect *Redirector) RedirectFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirect.greeter.SayHello(r.Context(), &v1.HelloRequest{Name: "kratos"})

		if r.URL.Path == "/helloworld/kratos" {
			http.Redirect(w, r, "https://go-kratos.dev/", http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func NewRedirector(greeter *service.GreeterService) *Redirector {
	return &Redirector{greeter}
}
