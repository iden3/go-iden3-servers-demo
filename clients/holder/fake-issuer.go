package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func main() {
	r := mux.NewRouter().StrictSlash(true)
	api := r.PathPrefix("/api/unstable").Subrouter()

	api.HandleFunc("/claim/request", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("/claim/request")
		resp := make(map[string]interface{})
		resp["id"] = 1
		json.NewEncoder(w).Encode(resp)
	})
	api.HandleFunc("/claim/status/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("/claim/status/{id}")
		resp := make(map[string]interface{})
		resp["status"] = "approved"
		// resp["claim"] = nil
		json.NewEncoder(w).Encode(resp)
	})
	api.HandleFunc("/claim/credential", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("/claim/credential")
		resp := make(map[string]interface{})
		resp["status"] = "ready"
		// resp["credential"] = nil
		json.NewEncoder(w).Encode(resp)
	})

	fmt.Println("server listening at port :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
