package main

type Sala struct {
	ID               string   `json:"id"`
	Nome             string   `json:"nome"`
	CapacidadeMaxima int      `json:"capacidade_maxima"`
	Recursos         []string `json:"recursos"`
}

type Aluno struct {
	ID    string `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

type Alocacao struct {
	SalaID        string `json:"sala_id"`
	DiaSemana     string `json:"dia_semana"`
	HorarioInicio string `json:"horario_inicio"` 
	HorarioTermino string `json:"horario_termino"`
}

type Turma struct {
	ID        string    `json:"id"`
	Nome      string    `json:"nome"`
	Disciplina string   `json:"disciplina"`
	Professor string    `json:"professor"`
	AlunosIDs []string  `json:"alunos_ids"`
	Alocacao  *Alocacao `json:"alocacao,omitempty"`
}

type TurmaResponse struct {
	ID             string    `json:"id"`
	Nome           string    `json:"nome"`
	Disciplina     string    `json:"disciplina"`
	Professor      string    `json:"professor"`
	QtdAlunos      int       `json:"qtd_alunos"`
	StatusAlocacao string    `json:"status_alocacao"`
	Alocacao       *Alocacao `json:"alocacao,omitempty"`
}

type MatriculaRequest struct {
	AlunoID string `json:"aluno_id"`
}

type ErrorResponse struct {
	Mensagem string `json:"mensagem"`
}