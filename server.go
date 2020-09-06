package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

func logging(next http.Handler) http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

var templates = template.Must(template.ParseFiles("./templates/base.html", "./templates/body.html"))

func index() http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := struct {
			Title template.HTML
			BusinessName string
		}{
			Title: template.HTML("SpaceWalkers"),
			BusinessName: "Business",
		}

		err := templates.ExecuteTemplate(w, "base", &s)
		if err != nil {
			http.Error(w, fmt.Sprintf("index: could not parse template: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func public() http.Handler  {
	return http.StripPrefix("/public/", http.FileServer(http.Dir("./public")))
}




func main() {
	mux := http.NewServeMux()
	mux.Handle("/public/",logging(public()))
	mux.Handle("/", logging(index()))

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "3030"
	}
	addr := fmt.Sprintf(":%s", port)
	server := http.Server{
		Addr: addr,
		Handler: mux,
		ReadHeaderTimeout: 15 *time.Second,
		WriteTimeout: 15 *time.Second,
		IdleTimeout: 15 *time.Second,
	}
	log.Println("main: running on port", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("main: could not start server: %v\n", err)
	}

}
