package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", HealthcheckHandler)

	mux.HandleFunc("POST /salas", CriarSalaHandler)
	mux.HandleFunc("GET /salas", ListarSalasHandler)

	mux.HandleFunc("POST /alunos", CriarAlunoHandler)
	mux.HandleFunc("GET /alunos", ListarAlunosHandler)
	mux.HandleFunc("GET /alunos/{id}", BuscarAlunoPorIDHandler)

	mux.HandleFunc("POST /turmas", CriarTurmaHandler)
	mux.HandleFunc("GET /turmas", ListarTurmasHandler)
	mux.HandleFunc("POST /turmas/{id}/matriculas", MatricularAlunoHandler)
	mux.HandleFunc("GET /turmas/{id}/alunos", ListarAlunosDaTurmaHandler)

	mux.HandleFunc("POST /turmas/{id}/alocacao", AlocarSalaHandler)

	fmt.Println("Servidor SGA rodando na porta :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Erro ao iniciar o servidor:", err)
	}
}
