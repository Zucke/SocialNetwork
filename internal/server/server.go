package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Zucke/SocialNetwork/internal/handlers"
	"github.com/Zucke/SocialNetwork/pkg/authentication"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

//Server is a struct that contains the param for the terver  like the port and the server of object
type Server struct {
	server *http.Server
	port   string
}

//Start put the server to listen
func (serv *Server) Start() {

	log.Printf("Escuchando en http://localhost:%s", serv.port)
	log.Fatal(serv.server.ListenAndServe())

}

func (serv *Server) getRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Post("/login", handlers.LoginUser)
	r.Post("/register", handlers.RegisterSocialUser)
	LogGroup := r.Group(nil)
	LogGroup.Use(authentication.ValidateMiddleware)

	LogGroup.Get("/my-friend", handlers.GetUserFriends)
	LogGroup.Get("/blocklist", handlers.GetUserBlockedUser)

	LogGroup.Post("/publicate", handlers.NewPublication)
	LogGroup.Get("/my-publications", handlers.GetAllUserPublication)
	LogGroup.Get("/my-publications/{publication_id}", handlers.GetUserPublication)
	LogGroup.Put("/my-publications/{publication_id}", handlers.ChangePublication)
	LogGroup.Delete("/my-publications/{publication_id}", handlers.DeletePublication)

	LogGroup.Get("/{user_id}", handlers.GetUserByID)
	LogGroup.Post("/{user_id}/add-to-friend", handlers.UserToFriendList)
	LogGroup.Post("/{user_id}/add-to-blocklist", handlers.UserToBlockList)

	LogGroup.Get("/{user_id}/publication/{publication_id}", handlers.GetUserPublication)

	LogGroup.Post("/{user_id}/publication/{publication_id}/like", handlers.NewPublicationLike)

	LogGroup.Post("/{user_id}/publication/{publication_id}/comment", handlers.NewPublicationComment)
	LogGroup.Put("/{user_id}/publication/{publication_id}/comment", handlers.ChangePublicationComment)
	LogGroup.Delete("/{user_id}/publication/{publication_id}/comment", handlers.DeletePublicationComment)
	LogGroup.Post("/{user_id}/publication/{publication_id}/comment/like", handlers.CommendLiked)

	return r

}

//New initialize the params for the server
func New(port string) *Server {
	serv := &Server{
		port: port,
	}

	r := serv.getRoutes()

	serv.server = &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return serv
}

//Close kill the server
func (serv *Server) Close(ctx context.Context) {
	log.Fatal(serv.server.Shutdown(ctx))
}
