package main

import (
	"encoding/json"
	"net/http"
)

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func respondError(w http.ResponseWriter, status int, mensagem string) {
	respondJSON(w, status, ErrorResponse{Mensagem: mensagem})
}

func temSobreposicao(inicio1, termino1, inicio2, termino2 string) bool {
	return inicio1 < termino2 && termino1 > inicio2
}