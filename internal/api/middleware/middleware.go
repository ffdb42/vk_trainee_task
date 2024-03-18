package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/ffdb42/vk_trainee_task/internal/constants"
	"github.com/ffdb42/vk_trainee_task/internal/db"
	"golang.org/x/crypto/bcrypt"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%v %v: %vms", r.Method, r.RequestURI, time.Since(start).Milliseconds())
	})
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name, pass, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("auth data was not provided"))
			return
		}

		user, err := db.Instance().GetUser(name)

		if err != nil || user == nil {
			log.Printf("ERROR %v %v: cannot get user from db: %v", r.Method, r.RequestURI, err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
		if err != nil || user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}

		if r.Method != http.MethodGet && user.Role != constants.AdminRole {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("forbidden"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
