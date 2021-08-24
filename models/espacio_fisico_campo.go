package models

type EspacioFisicoCampo struct {
	Id                int
	Valor             string
	EspacioFisicoId   *EspacioFisico
	CampoId           *Campo
	FechaInicio       string
	FechaFin          string
	FechaCreacion     string
	FechaModificacion string
}
