package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "UP",
		"horario": time.Now().Format(time.RFC3339),
		"versao":  "1.0.0",
	}
	respondJSON(w, http.StatusOK, response)
}

func CriarSalaHandler(w http.ResponseWriter, r *http.Request) {
	var sala Sala
	if err := json.NewDecoder(r.Body).Decode(&sala); err != nil {
		respondError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	if sala.ID == "" || sala.Nome == "" {
		respondError(w, http.StatusBadRequest, "ID e Nome são obrigatórios")
		return
	}

	if sala.CapacidadeMaxima <= 0 {
		respondError(w, http.StatusBadRequest, "Capacidade máxima deve ser maior que zero")
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.Salas[sala.ID]; exists {
		respondError(w, http.StatusConflict, "Sala com este ID já existe")
		return
	}

	if sala.Recursos == nil {
		sala.Recursos = []string{}
	}

	db.Salas[sala.ID] = sala
	respondJSON(w, http.StatusCreated, sala)
}

func ListarSalasHandler(w http.ResponseWriter, r *http.Request) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	salas := make([]Sala, 0, len(db.Salas))
	for _, sala := range db.Salas {
		salas = append(salas, sala)
	}

	respondJSON(w, http.StatusOK, salas)
}

func CriarAlunoHandler(w http.ResponseWriter, r *http.Request) {
	var aluno Aluno
	if err := json.NewDecoder(r.Body).Decode(&aluno); err != nil {
		respondError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	if aluno.ID == "" || aluno.Nome == "" || aluno.Email == "" {
		respondError(w, http.StatusBadRequest, "ID, Nome e Email são obrigatórios")
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.Alunos[aluno.ID]; exists {
		respondError(w, http.StatusConflict, "Aluno com este ID já existe")
		return
	}

	db.Alunos[aluno.ID] = aluno
	respondJSON(w, http.StatusCreated, aluno)
}

func ListarAlunosHandler(w http.ResponseWriter, r *http.Request) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	alunos := make([]Aluno, 0, len(db.Alunos))
	for _, aluno := range db.Alunos {
		alunos = append(alunos, aluno)
	}

	respondJSON(w, http.StatusOK, alunos)
}

func BuscarAlunoPorIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	db.mu.RLock()
	defer db.mu.RUnlock()

	aluno, exists := db.Alunos[id]
	if !exists {
		respondError(w, http.StatusNotFound, "Aluno não encontrado")
		return
	}

	respondJSON(w, http.StatusOK, aluno)
}

func CriarTurmaHandler(w http.ResponseWriter, r *http.Request) {
	var turma Turma
	if err := json.NewDecoder(r.Body).Decode(&turma); err != nil {
		respondError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	if turma.ID == "" || turma.Nome == "" || turma.Disciplina == "" || turma.Professor == "" {
		respondError(w, http.StatusBadRequest, "ID, Nome, Disciplina e Professor são obrigatórios")
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.Turmas[turma.ID]; exists {
		respondError(w, http.StatusConflict, "Turma com este ID já existe")
		return
	}

	turma.AlunosIDs = []string{}
	turma.Alocacao = nil

	db.Turmas[turma.ID] = turma
	respondJSON(w, http.StatusCreated, turma)
}

func ListarTurmasHandler(w http.ResponseWriter, r *http.Request) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	response := make([]TurmaResponse, 0, len(db.Turmas))
	for _, turma := range db.Turmas {
		status := "NAO_ALOCADA"
		if turma.Alocacao != nil {
			status = "ALOCADA"
		}

		tr := TurmaResponse{
			ID:             turma.ID,
			Nome:           turma.Nome,
			Disciplina:     turma.Disciplina,
			Professor:      turma.Professor,
			QtdAlunos:      len(turma.AlunosIDs),
			StatusAlocacao: status,
			Alocacao:       turma.Alocacao,
		}
		response = append(response, tr)
	}

	respondJSON(w, http.StatusOK, response)
}

func MatricularAlunoHandler(w http.ResponseWriter, r *http.Request) {
	turmaID := r.PathValue("id")

	var req MatriculaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	turma, turmaExists := db.Turmas[turmaID]
	if !turmaExists {
		respondError(w, http.StatusNotFound, "Turma não encontrada")
		return
	}

	_, alunoExists := db.Alunos[req.AlunoID]
	if !alunoExists {
		respondError(w, http.StatusNotFound, "Aluno não encontrado")
		return
	}

	for _, id := range turma.AlunosIDs {
		if id == req.AlunoID {
			respondError(w, http.StatusConflict, "Aluno já está matriculado nesta turma")
			return
		}
	}

	if turma.Alocacao != nil {
		sala, salaExists := db.Salas[turma.Alocacao.SalaID]
		if salaExists && len(turma.AlunosIDs)+1 > sala.CapacidadeMaxima {
			respondError(w, http.StatusUnprocessableEntity, "Capacidade máxima da sala excedida para esta turma")
			return
		}
	}

	if turma.Alocacao != nil {
		for _, outraTurma := range db.Turmas {
			if outraTurma.ID == turma.ID || outraTurma.Alocacao == nil {
				continue
			}

			alunoNaOutra := false
			for _, id := range outraTurma.AlunosIDs {
				if id == req.AlunoID {
					alunoNaOutra = true
					break
				}
			}

			if alunoNaOutra && outraTurma.Alocacao.DiaSemana == turma.Alocacao.DiaSemana {
				if temSobreposicao(turma.Alocacao.HorarioInicio, turma.Alocacao.HorarioTermino,
					outraTurma.Alocacao.HorarioInicio, outraTurma.Alocacao.HorarioTermino) {
					respondError(w, http.StatusConflict, "Aluno já possui aula em outra turma no mesmo horário e dia")
					return
				}
			}
		}
	}

	turma.AlunosIDs = append(turma.AlunosIDs, req.AlunoID)
	db.Turmas[turmaID] = turma

	respondJSON(w, http.StatusOK, map[string]string{"mensagem": "Aluno matriculado com sucesso"})
}

func ListarAlunosDaTurmaHandler(w http.ResponseWriter, r *http.Request) {
	turmaID := r.PathValue("id")

	db.mu.RLock()
	defer db.mu.RUnlock()

	turma, exists := db.Turmas[turmaID]
	if !exists {
		respondError(w, http.StatusNotFound, "Turma não encontrada")
		return
	}

	alunos := make([]Aluno, 0, len(turma.AlunosIDs))
	for _, id := range turma.AlunosIDs {
		if aluno, ok := db.Alunos[id]; ok {
			alunos = append(alunos, aluno)
		}
	}

	respondJSON(w, http.StatusOK, alunos)
}

func AlocarSalaHandler(w http.ResponseWriter, r *http.Request) {
	turmaID := r.PathValue("id")

	var alocacao Alocacao
	if err := json.NewDecoder(r.Body).Decode(&alocacao); err != nil {
		respondError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	if alocacao.SalaID == "" || alocacao.DiaSemana == "" || alocacao.HorarioInicio == "" || alocacao.HorarioTermino == "" {
		respondError(w, http.StatusBadRequest, "Todos os campos da alocação são obrigatórios")
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	turma, turmaExists := db.Turmas[turmaID]
	if !turmaExists {
		respondError(w, http.StatusNotFound, "Turma não encontrada")
		return
	}

	sala, salaExists := db.Salas[alocacao.SalaID]
	if !salaExists {
		respondError(w, http.StatusNotFound, "Sala não encontrada")
		return
	}

	if len(turma.AlunosIDs) > sala.CapacidadeMaxima {
		respondError(w, http.StatusUnprocessableEntity, "Capacidade da sala é menor que o número atual de alunos da turma")
		return
	}

	for _, outraTurma := range db.Turmas {
		if outraTurma.ID != turma.ID && outraTurma.Alocacao != nil {
			if outraTurma.Alocacao.SalaID == alocacao.SalaID && outraTurma.Alocacao.DiaSemana == alocacao.DiaSemana {
				if temSobreposicao(alocacao.HorarioInicio, alocacao.HorarioTermino,
					outraTurma.Alocacao.HorarioInicio, outraTurma.Alocacao.HorarioTermino) {
					respondError(w, http.StatusConflict, "A sala já está alocada para outra turma neste dia e horário")
					return
				}
			}
		}
	}

	turma.Alocacao = &alocacao
	db.Turmas[turmaID] = turma

	respondJSON(w, http.StatusOK, turma)
}