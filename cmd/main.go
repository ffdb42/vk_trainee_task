package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/ffdb42/vk_trainee_task/docs"
	"github.com/ffdb42/vk_trainee_task/internal/api/middleware"
	"github.com/ffdb42/vk_trainee_task/internal/api/server"
	"github.com/ffdb42/vk_trainee_task/internal/db"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title           VK backend trainee task API
// @version         1.0

// @securityDefinitions.basic  BasicAuth

// @host      localhost:8888
// @BasePath  /
func main() {
	db.Init()

	mux := http.NewServeMux()
	server := server.Server{}
	port := os.Getenv("API_INT_PORT")

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello world!"))
	})
	mux.HandleFunc("/sign-up/", server.SignUp)
	mux.Handle("/actor/", middleware.Authenticate(http.HandlerFunc(server.ActorHandler)))
	mux.Handle("/film/", middleware.Authenticate(http.HandlerFunc(server.FilmHandler)))
	mux.Handle("/search/", middleware.Authenticate(http.HandlerFunc(server.SearchHandler)))
	mux.HandleFunc("/swagger/", httpSwagger.Handler(httpSwagger.URL(fmt.Sprintf("http://localhost:%v/swagger/doc.json", port))))

	handler := middleware.Logger(mux)

	log.Printf("starting server on port %v", port)
	err := http.ListenAndServe(":"+port, handler)
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server closed due error: %v", err)
	}
}
