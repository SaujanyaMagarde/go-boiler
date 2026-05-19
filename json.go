package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string){
	if code > 499 {
		log.Printf("5xx error code %d: %s",code,msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	// {
	// 	"error" : ""
	// }
	respondWithJson(w,code,errorResponse{Error: msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}){
	dat , err := json.Marshal(payload)
	if err != nil{
		log.Println("failes to marshel json repsonse")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type","application/json")
	w.WriteHeader(code)
	w.Write(dat)
}



