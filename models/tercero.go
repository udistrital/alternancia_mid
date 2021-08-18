package models

import "time"

type Tercero struct {
	Id                  int
	NombreCompleto      string
	PrimerNombre        string
	SegundoNombre       string
	PrimerApellido      string
	SegundoApellido     string
	LugarOrigen         int
	FechaNacimiento     time.Time
	Activo              bool
	TipoContribuyenteId *TipoContribuyente
	UsuarioWSO2         string
}
