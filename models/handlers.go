package models

import (
	"encoding/json"
	"go-check/config"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

//Client Hash
type Client struct {
	Hash string
}

//GetTask json get
func GetTask(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id, _ := strconv.Atoi(p.ByName("id"))

		if r.Method != "GET" {
			http.Error(w, http.StatusText(405), 405)
			return
		}

		bks, err := GetTasksReq(env.DB, id)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err = json.NewEncoder(w).Encode(bks); err != nil {
			w.WriteHeader(500)
		}
	}
}

//Activate post json
func Activate(env *config.Env) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var u Client

		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), 405)
			return
		}

		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		err = InsertHash(env.DB, u.Hash)
		if err != nil {
			http.Error(w, "hash key not found", 500)
			return
		}
		w.WriteHeader(200)
	}
}
