package main

import (
	"flag"
	"go-check/config"
	"go-check/models"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	var (
		server = flag.Bool("server", true, "use -server for start restapi")
	)
	flag.Parse()

	db, err := config.NewDB("postgres://checker:checker@localhost/cheker?sslmode=disable")
	if err != nil {
		log.Print(err)
	}
	env := &config.Env{DB: db}

	if *server == true {
		r := httprouter.New()
		r.POST("/api/v1/activate", models.Activate(env))
		r.GET("/api/v1/gettask/:id", models.GetTask(env))
		http.ListenAndServe(":3000", r)
	}
}
