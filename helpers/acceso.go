package helpers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/alternancia_mid/models"
)

func Autorizacion(idQr string, idScan string, salon string, idEdificio string, idSede string, tipoScan string) (persona models.Persona, outputError map[string]interface{}) {

	var respuesta_peticion []models.InfoComplementariaTercero

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "/Autorizacion", "err": err, "status": "502"}
			panic(outputError)
		}
	}()
	if response, err := getJsonTest(beego.AppConfig.String("UrlCrudTerceros")+"info_complementaria_tercero/?limit=-1&query=tercero_id:"+idQr, &respuesta_peticion); (err == nil) && (response == 200) {
		if len(respuesta_peticion) != 0 {

			//Declaracion de variables a usar
			var comorbilidad bool
			var vacuna bool
			var sintomas bool
			var cupoDisponible bool
			var permiso bool
			var clase bool
			var aforo int
			var cupo models.EspacioFisicoCampo
			id := idSede
			var idRol int
			//var materiasDia []models.CargaAcademica

			//Consulta de comorbilidades
			comorbilidad, err1 := ConsultarComorbilidades(idQr)
			if err1 != nil {
				logs.Error(err1)
				return models.Persona{}, err1
			}

			//Consulta de vacunacion
			vacuna, err1 = ConsultarVacuna(idQr)
			if err1 != nil {
				logs.Error(err1)
				return models.Persona{}, err1
			}

			//Consulta de sintomas
			sintomas, err1 = ConsultarSintomas(idQr)
			if err1 != nil {
				logs.Error(err1)
				return models.Persona{}, err1
			}

			if idEdificio != "" {
				id = idEdificio
			}
			if salon != "" {
				//Busqueda de id de espacio
				var respuesta_peticion_salones []models.EspacioFisicoPadre
				salon = strings.ToUpper(salon)
				coincidencias := 0
				if response, err := getJsonTest(beego.AppConfig.String("UrlCrudOikos")+"espacio_fisico_padre/?limit=-1&query=PadreId.Id:"+idEdificio, &respuesta_peticion_salones); (err == nil) && (response == 200) {
					if len(respuesta_peticion_salones) > 0 {
						for _, espacio := range respuesta_peticion_salones {
							hijo := espacio.Hijo
							if strings.Contains(hijo.Nombre, salon) {
								id = strconv.Itoa(hijo.Id)
								coincidencias = coincidencias + 1
							}
						}
						if coincidencias > 1 {
							logs.Error("Se encontró más de un salon")
							outputError = map[string]interface{}{"funcion": "/ConsultaIdSalon", "err": "Se encontró más de un salon que coincide con la busqueda", "status": "502"}
							return models.Persona{}, outputError
						} else if coincidencias == 0 {
							logs.Error("No se encontró ningun salon")
							outputError = map[string]interface{}{"funcion": "/ConsultaIdSalon", "err": "No se encontró ningun salon", "status": "502"}
							return models.Persona{}, outputError
						}
					} else {
						logs.Error(err)
						outputError = map[string]interface{}{"funcion": "/ConsultaIdSalon", "err": "Se encontró más de un salon que coincide con la busqueda", "status": "502"}
						return models.Persona{}, outputError
					}
				} else {
					logs.Error(err)
					outputError = map[string]interface{}{"funcion": "/ConsultaIdSalon", "err": err, "status": "502"}
					return models.Persona{}, outputError
				}
			}

			//Consulta de aforo
			aforo, err1 = ConsultarAforo(id)
			if err1 != nil {
				logs.Error(err1)
				outputError = map[string]interface{}{"funcion": "/ConsultarAforo", "err": err1, "status": "502"}
				return models.Persona{}, outputError
			}

			//Consulta de cupo
			err1 = ConsultarCupo(id, &cupo)
			if err1 != nil {
				logs.Error(err1)
				outputError = map[string]interface{}{"funcion": "/ConsultarCupo", "err": err1, "status": "502"}
				return models.Persona{}, outputError
			}
			persona.Cupo, err = strconv.Atoi(cupo.Valor)
			if err != nil {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/ConsultarCupo", "err": err, "status": "502"}
				return models.Persona{}, outputError
			}
			if tipoScan == "in" {
				cupoDisponible = persona.Cupo > 0
			} else if tipoScan == "out" {
				cupoDisponible = persona.Cupo < aforo
			}

			//Consulta de roles para conceder permisos
			var respuesta_peticion_permisos []models.Vinculacion
			if response, err := getJsonTest(beego.AppConfig.String("UrlCrudTerceros")+"vinculacion/?query=Activo:true,TerceroPrincipalId.Id:"+idScan, &respuesta_peticion_permisos); (err == nil) && (response == 200) {
				for _, vinculacion := range respuesta_peticion_permisos {
					idRol = vinculacion.TipoVinculacionId
					if idRol == 377 || idRol == 292 || idRol == 294 || (idRol >= 296 && idRol <= 299) {
						permiso = true
					}
				}
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/ConsultarComorbilidades", "err": err, "status": "502"}
				return models.Persona{}, outputError
			}

			//Consulta de horario
			/*respuesta := append(materiasDia,
				models.CargaAcademica{
					CODIGO_DIA: "7",
					HORA:       "15",
					SALON:      "AULA 201",
				}, models.CargaAcademica{
					CODIGO_DIA: "4",
					HORA:       "10",
					SALON:      "AULA 501",
				})
			var minutosAHora float64 = float64(time.Now().Local().Minute()) / 60
			var horaActual float64 = math.Round(float64(time.Now().Local().Hour()) + minutosAHora)
			for _, materia := range respuesta {
				if hora, err := strconv.ParseFloat(materia.HORA, 64); hora == horaActual && err == nil {
					clase = true
				}
			}*/
			clase = true

			//Rellenado de modelo persona
			persona.Nombre = respuesta_peticion[0].TerceroId.NombreCompleto
			persona.Fecha = time.Now()
			if (vacuna || !comorbilidad) && !sintomas && cupoDisponible && clase {
				persona.Acceso = "Autorizado"
				if permiso {
					valor, err1 := ActualizarCupo(cupo, tipoScan)
					if err1 != nil {
						logs.Error(err1)
						outputError = map[string]interface{}{"funcion": "/Autorizacion", "err": err1, "status": "502"}
						return models.Persona{}, outputError
					}
					persona.Cupo = valor
				}
			} else {
				persona.Acceso = "No autorizado"
				if comorbilidad && !vacuna {
					persona.Causa = "Presenta comorbilidad(es) y no tiene vacunacion"
				} else if sintomas {
					persona.Causa = "Presenta sintomas"
				} else if !cupoDisponible {
					persona.Causa = "No hay cupo disponible en el espacio"
				} else if !clase {
					persona.Causa = "No tiene clase en este momento"
				}
			}
		} else {
			logs.Error("No hay datos de caracterización registrados para el usuario")
			outputError = map[string]interface{}{"funcion": "/Autorizacion", "err": "No hay datos de caracterización registrados para el usuario", "status": "502"}
			return models.Persona{}, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/Autorizacion", "err": err, "status": "502"}
		return models.Persona{}, outputError
	}

	return
}

func ActualizarAforo(idPersona string, idEspacio string, tipoQr string) (persona models.Persona, outputError map[string]interface{}) {
	var cupo models.EspacioFisicoCampo
	err := ConsultarCupo(idEspacio, &cupo)
	var respuesta_peticion []models.InfoComplementariaTercero
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ActualizarAforo", "err": err, "status": "502"}
		return models.Persona{}, outputError
	}
	if response, err := getJsonTest(beego.AppConfig.String("UrlCrudTerceros")+"info_complementaria_tercero/?limit=-1&query=tercero_id:"+idPersona, &respuesta_peticion); (err == nil) && (response == 200) {
		if len(respuesta_peticion) != 0 {
			persona.Nombre = respuesta_peticion[0].TerceroId.NombreCompleto
			persona.Fecha = time.Now()
		} else {
			logs.Error("No hay datos de caracterización registrados para el usuario")
			outputError = map[string]interface{}{"funcion": "/Autorizacion", "err": "No hay datos de caracterización registrados para el usuario", "status": "502"}
			return models.Persona{}, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/Autorizacion", "err": err, "status": "502"}
		return models.Persona{}, outputError
	}
	if tipoQr == "in" {
		valorCupo, err1 := strconv.Atoi(cupo.Valor)
		if err1 != nil {
			logs.Error(err1)
			outputError = map[string]interface{}{"funcion": "/ActualizarAforo", "err": err1, "status": "502"}
			return models.Persona{}, outputError
		}
		comorbilidades, err := ConsultarComorbilidades(idPersona)
		if err != nil {
			logs.Error(err)
			return models.Persona{}, err
		}
		vacuna, err := ConsultarVacuna(idPersona)
		if err != nil {
			logs.Error(err)
			return models.Persona{}, err
		}
		sintomas, err := ConsultarSintomas(idPersona)
		if err != nil {
			logs.Error(err)
			return models.Persona{}, err
		}
		if valorCupo > 0 && (vacuna || !comorbilidades) && !sintomas {
			persona.Acceso = "Autorizado"
			_, err = ActualizarCupo(cupo, tipoQr)
			if err != nil {
				logs.Error(err)
				return models.Persona{}, err
			}
		} else {
			persona.Acceso = "No autorizado"
			if valorCupo == 0 {
				persona.Causa = "No hay cupo disponible en el espacio"
			} else if comorbilidades {
				persona.Causa = "Presenta comorbilidades"
			} else if sintomas {
				persona.Causa = "Presenta sintomas"
			}
		}
	} else if tipoQr == "out" {
		aforo, err := ConsultarAforo(idEspacio)
		if err != nil {
			logs.Error(err)
			return models.Persona{}, err
		}
		if valor, err := strconv.Atoi(cupo.Valor); valor < aforo && err == nil {
			_, err1 := ActualizarCupo(cupo, tipoQr)
			if err1 != nil {
				logs.Error(err1)
				return models.Persona{}, err1
			}
		} else if err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{"funcion": "/ActualizarAforo", "err": err, "status": "502"}
			return models.Persona{}, outputError
		}
	}
	return
}

func ConsultarAforo(id string) (aforo int, outputError map[string]interface{}) {
	var respuesta_peticion_aforo []models.EspacioFisicoCampo
	if response, err := getJsonTest(beego.AppConfig.String("UrlCrudOikos")+"espacio_fisico_campo/?query=CampoId.Id:5,EspacioFisicoId.Id:"+id, &respuesta_peticion_aforo); (err == nil) && (response == 200) {
		if len(respuesta_peticion_aforo) != 0 {
			aforoStr := respuesta_peticion_aforo[0].Valor
			aforo, err = strconv.Atoi(aforoStr)
			if err != nil {
				logs.Error(err)
				outputError = map[string]interface{}{"funcion": "/ConsultarAforo", "err": err, "status": "502"}
				return 0, outputError
			}
		} else {
			outputError = map[string]interface{}{"funcion": "/ConsultarAforo", "err": "No hay aforo registrado para el espacio", "status": "502"}
			return 0, outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ConsultarAforo", "err": err, "status": "502"}
		return 0, outputError
	}
	return
}

func ConsultarCupo(id string, cupo interface{}) (outputError map[string]interface{}) {
	var respuesta_peticion_cupo []map[string]interface{}
	if response, err := getJsonTest(beego.AppConfig.String("UrlCrudOikos")+"espacio_fisico_campo/?query=CampoId.Id:4,EspacioFisicoId.Id:"+id, &respuesta_peticion_cupo); (err == nil) && (response == 200) {
		if len(respuesta_peticion_cupo) != 0 {
			res, err1 := json.Marshal(respuesta_peticion_cupo[0])
			if err1 != nil {
				logs.Error(err1)
				outputError = map[string]interface{}{"funcion": "/ConsultarCupo", "err": err1, "status": "502"}
				return outputError
			}
			json.Unmarshal(res, &cupo)
		} else {
			outputError = map[string]interface{}{"funcion": "/ConsultarCupo", "err": "El campo cupo no tiene registros", "status": "502"}
			return outputError
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ConsultarCupo", "err": err, "status": "502"}
		return outputError
	}
	return
}

func ActualizarCupo(cupo models.EspacioFisicoCampo, tipoScan string) (valor int, outputError map[string]interface{}) {
	valor, err := strconv.Atoi(cupo.Valor)
	if err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ActualizarCupo", "err": err, "status": "502"}
		return 0, outputError
	}
	if tipoScan == "in" {
		valor = valor - 1
	} else if tipoScan == "out" {
		valor = valor + 1
	}
	cupo.Valor = strconv.Itoa(valor)
	if err := putJson(beego.AppConfig.String("UrlCrudOikos")+"espacio_fisico_campo", strconv.Itoa(cupo.Id), cupo); err != nil {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ActualizarCupo", "err": err, "status": "502"}
		return 0, outputError
	}
	return
}

func ConsultarSintomas(idQr string) (sintomas bool, outputError map[string]interface{}) {
	var respuesta_peticion_sintomas map[string]interface{}
	var sintoma []models.Sintomas
	if response, err := getJsonTest(beego.AppConfig.String("UrlCrudSintomas")+"sintomas?limit=1&order=desc&sortby=fecha_creacion&query=terceroId:"+idQr, &respuesta_peticion_sintomas); (err == nil) && (response == 200) {
		LimpiezaRespuestaRefactor(respuesta_peticion_sintomas, &sintoma)
		sintomasRegistrados := sintoma[0].InfoSalud
		sintomas = (sintomasRegistrados.Agotamiento ||
			sintomasRegistrados.CongestionNasal ||
			sintomasRegistrados.ContactoCovid ||
			sintomasRegistrados.DificultadRespiratoria ||
			sintomasRegistrados.EstadoEmbarazo ||
			sintomasRegistrados.Fiebre ||
			sintomasRegistrados.MalestarGeneral)
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ConsultarSintomas", "err": err, "status": "502"}
		return false, outputError
	}
	return
}

func ConsultarVacuna(idQr string) (vacuna bool, outputError map[string]interface{}) {
	var respuesta_peticion_vacuna []models.InfoComplementariaTercero
	if response, err := getJsonTest(beego.AppConfig.String("UrlCrudTerceros")+"info_complementaria_tercero/?query=tercero_id:"+idQr+",InfoComplementariaId:305", &respuesta_peticion_vacuna); (err == nil) && (response == 200) {
		if len(respuesta_peticion_vacuna) != 0 {
			layout := "2006-01-02T15:04:05.000Z"
			var dato map[string]interface{}
			json.Unmarshal([]byte(respuesta_peticion_vacuna[0].Dato), &dato)
			var fecha = dato["dato"]
			data, _ := json.Marshal(fecha)
			if string(data) != "\"\"" {
				str := fmt.Sprintf("%v", fecha)
				t, err := time.Parse(layout, str)

				if err != nil {
					outputError = map[string]interface{}{"funcion": "/ConsultarVacuna", "err": err, "status": "502"}
					return false, outputError
				}
				duracion := time.Since(t)
				dias := int(duracion.Hours() / 24)
				if dias > 14 {
					vacuna = true
				}
			}
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ConsultarVacuna", "err": err, "status": "502"}
		return false, outputError
	}
	return
}

func ConsultarComorbilidades(idQr string) (comorbilidad bool, outputError map[string]interface{}) {
	var respuesta_peticion_comorbilidades []models.InfoComplementariaTercero
	var dato map[string]interface{}
	if response, err := getJsonTest(beego.AppConfig.String("UrlCrudTerceros")+"info_complementaria_tercero/?limit=-1&query=tercero_id:"+idQr+",InfoComplementariaId.GrupoInfoComplementariaId.Id:47", &respuesta_peticion_comorbilidades); (err == nil) && (response == 200) {
		for _, info := range respuesta_peticion_comorbilidades {
			json.Unmarshal([]byte(info.Dato), &dato)
			if dato["dato"] == true {
				comorbilidad = true
			}
		}
	} else {
		logs.Error(err)
		outputError = map[string]interface{}{"funcion": "/ConsultarComorbilidades", "err": err, "status": "502"}
		return false, outputError
	}
	return
}
