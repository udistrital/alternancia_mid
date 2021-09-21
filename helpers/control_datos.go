package helpers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/alternancia_mid/models"
)

var idEstudiante string
var persona models.Tercero

func ConsultarTraza(idPersona string) (trazaRes models.TrazaEstudiante, outputError map[string]interface{}) {
	idEstudiante = idPersona

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/GetTraza", "err": err, "status": "502"}
			panic(outputError)
		}
	}()

	if status, err := getJsonTest(beego.AppConfig.String("UrlCrudTerceros")+"tercero/"+idEstudiante, &persona); status != 200 || err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/GetTraza/GetTercero", "err": err, "responseStatus": status, "status": "502"}
		return models.TrazaEstudiante{}, outputError
	}

	traza, err := RegistrarTraza()
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/GetTraza/RegistrarTraza", "err": err, "status": "502"}
		return models.TrazaEstudiante{}, outputError
	}

	return traza, nil
}

func RegistrarTraza() (traza models.TrazaEstudiante, outputError map[string]interface{}) {
	var res map[string]interface{}
	var res_seguimiento []models.RegistroTraza
	if status, err := getJsonTest(beego.AppConfig.String("UrlCrudSeguimiento")+"seguimiento/?limit=0&query=tercero_id:"+idEstudiante+"&sortby=fecha_creacion&order=asc", &res); status != 200 || err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/GetTraza/GetSeguimiento", "err": err, "responseStatus": status, "status": "502"}
		return models.TrazaEstudiante{}, outputError
	}
	LimpiezaRespuestaRefactor(res, &res_seguimiento)
	traza.Estudiante = persona.NombreCompleto
	var temp int
	var tempAula int
	for _, reg := range res_seguimiento {
		var espacio models.EspacioFisico
		if status, err := getJsonTest(beego.AppConfig.String("UrlCrudOikos")+"espacio_fisico/"+strconv.Itoa(reg.EspacioId), &espacio); status != 200 || err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "/GetTraza/GetOikos/id:" + strconv.Itoa(reg.EspacioId), "err": err, "responseStatus": status, "status": "502"}
			return models.TrazaEstudiante{}, outputError
		}
		if reg.TipoEspacioId == 1 {
			if reg.TipoEscaneo == "I" {
				traza.Sedes = append(traza.Sedes, struct {
					Nombre       string
					FechaEntrada string
					FechaSalida  string
					Aulas        []struct {
						Nombre       string
						FechaEntrada string
						FechaSalida  string
					}
				}{
					Nombre:       espacio.Nombre,
					FechaEntrada: reg.FechaCreacion,
				})
				temp = len(traza.Sedes) - 1
			} else if reg.TipoEscaneo == "S" {
				traza.Sedes[temp].FechaSalida = reg.FechaCreacion
				temp = -1
			}
		} else if reg.TipoEspacioId == 3 {
			if reg.TipoEscaneo == "I" {
				traza.Sedes[temp].Aulas = append(traza.Sedes[temp].Aulas, struct {
					Nombre       string
					FechaEntrada string
					FechaSalida  string
				}{
					Nombre:       espacio.Nombre,
					FechaEntrada: reg.FechaCreacion,
				})
				tempAula = len(traza.Sedes[temp].Aulas) - 1
			} else if reg.TipoEscaneo == "S" {
				traza.Sedes[temp].Aulas[tempAula].FechaSalida = reg.FechaCreacion
				tempAula = -1
			}

		}
	}
	return
}
