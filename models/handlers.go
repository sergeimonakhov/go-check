package models

import (
	"encoding/json"
	"fmt"
	"go-check/config"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//Client Hash
type Client struct {
	Hash string
}

//GetTask json get
func GetTask(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var u Client
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		fmt.Println(u.Hash)
	}
}
