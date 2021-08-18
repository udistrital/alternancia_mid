package models

type Sintomas struct {
	TerceroId int `json:"terceroId"`
	InfoSalud struct {
		Fiebre                 bool `json:"fiebre"`
		CongestionNasal        bool `json:"congestion_nasal"`
		DificultadRespiratoria bool `json:"dificultad_respiratoria"`
		Agotamiento            bool `json:"agotamiento"`
		MalestarGeneral        bool `json:"malestar_general"`
		EstadoEmbarazo         bool `json:"estado_embarazo"`
		ContactoCovid          bool `json:"contacto_covid"`
	} `json:"info_salud"`
	Activo bool `json:"activo"`
}
