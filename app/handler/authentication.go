package handler

import (
	"github.com/erpe/image_service_go/app/config"
	"log"
	"net/http"
	"strings"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		used := strings.TrimPrefix(token, "Bearer ")

		atoken := config.GetConfig().Server.Token

		if used == atoken {
			next.ServeHTTP(w, r)
		} else {
			log.Println("Wrong token: \"" + token + "\"")
			respondError(w, http.StatusForbidden, "wrong token")
		}
	})
}
