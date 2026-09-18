package main

import "sync"

type MemoryDB struct {
	mu     sync.RWMutex
	Salas  map[string]Sala
	Alunos map[string]Aluno
	Turmas map[string]Turma
}

func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		Salas:  make(map[string]Sala),
		Alunos: make(map[string]Aluno),
		Turmas: make(map[string]Turma),
	}
}

var db = NewMemoryDB()