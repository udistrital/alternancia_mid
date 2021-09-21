package models

type RegistroTraza struct {
	TerceroId         int    `json:"tercero_id"`
	EspacioId         int    `json:"oikos_id"`
	TipoEspacioId     int    `json:"tipo_espacio_id"`
	TipoEscaneo       string `json:"tipo_registro"`
	Activo            bool   `json:"activo"`
	FechaCreacion     string `json:"fecha_creacion"`
	FechaModificacion string `json:"fecha_modificacion"`
}
