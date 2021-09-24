package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"journal/auth"
	"journal/controllers"
	"journal/db"
	"log"
	"net/http"
	"time"
)

func main() {

	r := chi.NewRouter()

	db.Connect()

	r.Use(middleware.Timeout(time.Second * 60))
	r.Use(timingMiddleware)
	r.Get("/", baseIndexHandler)
	r.Get("/token", controllers.GetToken)

	r.Route("/subjects", func(r chi.Router) {
		r.Use(auth.Authenticator)
		r.Get("/", controllers.GetSubjects)
		r.Get("/{firstname}", controllers.GetSubject)
		r.Put("/", controllers.UpdateSubject)
		r.Post("/", controllers.InsertSubject)
		r.Delete("/{subjectid}", controllers.DeleteSubject)
	})

	r.Route("/sessions", func(r chi.Router) {
		r.Use(auth.Authenticator)
		r.Get("/", controllers.GetSessions)
		r.Get("/{subjectid}", controllers.GetSessionForSubject)
		r.Get("/{sessionId}", controllers.GetSessionForSessionId)
		r.Put("/", controllers.UpdateSession)
		r.Post("/", controllers.InsertSession)
		r.Delete("/{sessionid}", controllers.DeleteSession)
	})

	r.Route("/tags", func(r chi.Router) {
		r.Use(auth.Authenticator)
		r.Get("/", controllers.GetTags)
		r.Put("/{name}", controllers.InsertTag)
		r.Delete("/{name}", controllers.DeleteTag)
	})
	r.Post("/", postHandler)

	http.ListenAndServe(":8080", r)
}

func baseIndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the base page..."))
}

type person struct {
	Name    string `json:"name"`
	Surname string `json:"lastname"`
}

func postHandler(w http.ResponseWriter, r *http.Request) {

	p := person{}

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		errText := err.Error()
		http.Error(w, errText, 400)
		return
	}

	fmt.Fprint(w, "Object was name: "+p.Name+", surname: "+p.Surname)
}

func timingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeNow := time.Now()
		log.Println("Request started at:" + timeNow.String())
		next.ServeHTTP(w, r)

		log.Println("Request completed at: " + time.Now().String())

		timeAtTheEnd := time.Now()
		log.Println("Request lasted: " + timeAtTheEnd.Sub(timeNow).String())
	})
}
