package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Welcome home!")
	json.NewEncoder(w).Encode(map[string]string{"sdfsd": "sdewrwe"})
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	log.Fatal(http.ListenAndServe(":8080", router))
	// json.Unmarshal()
	// json.NewEncoder(w).Encode(map[string]string{"sdfsd": "sdewrwe"})
}
