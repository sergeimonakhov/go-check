package main

import (
	"fmt"
	"go-check/config"
	"go-check/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {

	db, err := config.NewDB("postgres://checker:checker@localhost/cheker?sslmode=disable")
	if err != nil {
		fmt.Println(err.Error())
	}
	env := &config.Env{DB: db}

	r := httprouter.New()
	r.POST("/api/v1/activate", models.Activate(env))
	r.GET("/api/v1/gettask/:id", models.GetTask(env))
	r.POST("/api/v1/statusupdate", models.StatusUpdate(env))

	http.ListenAndServe(":3000", r)
}
