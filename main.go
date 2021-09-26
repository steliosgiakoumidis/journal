package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"journal/auth"
	"journal/controllers"
	"journal/db"
	customMiddleware "journal/middleware"
	"net/http"
	"time"
)

var tagHandler *controllers.TagHandler
var sessionHandler *controllers.SessionHandler
var subjectHandler *controllers.SubjectHandler
var loginHandler *controllers.LoginHandler
var middlewareHandler *customMiddleware.CustomMiddleware

func init() {

	db.Connect()
	middlewareHandler = customMiddleware.NewCustomMiddleware(auth.Auth{})
	loginHandler = controllers.NewLoginHandler(auth.Auth{})
	tagHandler = controllers.NewTagHandler(db.TagRepository{})
	sessionHandler = controllers.NewSessionHandler(db.SessionRepository{db.TagRepository{}})
	subjectHandler = controllers.NewSubjectHandler(db.SubjectRepository{})
}

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Timeout(time.Second * 60))
	r.Use(middlewareHandler.TimingMiddleware)
	r.Get("/", baseIndexHandler)
	r.Get("/login", loginHandler.GetToken)

	r.Route("/subjects", func(r chi.Router) {
		r.Use(middlewareHandler.Authenticator)
		r.Get("/", subjectHandler.GetSubjects)
		r.Get("/{subjectId}", subjectHandler.GetSubject)
		r.Put("/", subjectHandler.UpdateSubject)
		r.Post("/", subjectHandler.InsertSubject)
		r.Delete("/{subjectId}", subjectHandler.DeleteSubject)
	})

	r.Route("/sessions", func(r chi.Router) {
		r.Use(middlewareHandler.Authenticator)
		r.Get("/", sessionHandler.GetSessions)
		r.Get("/{subjectId}", sessionHandler.GetSessionsForSubject)
		r.Get("/{sessionId}", sessionHandler.GetSessionById)
		r.Put("/", sessionHandler.UpdateSession)
		r.Post("/", sessionHandler.InsertSession)
		r.Delete("/{sessionId}", sessionHandler.DeleteSession)
	})

	r.Route("/tags", func(r chi.Router) {
		r.Use(middlewareHandler.Authenticator)
		r.Get("/", tagHandler.GetTags)
		r.Post("/{name}", tagHandler.InsertTag)
		r.Delete("/{tagId}", tagHandler.DeleteTag)
	})

	http.ListenAndServe(":8080", r)
}

func baseIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Journal!!"))
}
