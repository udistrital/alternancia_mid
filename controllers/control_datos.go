package controllers

import (
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/alternancia_mid/helpers"
)

// Control_datosController operations for Control_datos
type Control_datosController struct {
	beego.Controller
}

// URLMapping ...
func (c *Control_datosController) URLMapping() {
	c.Mapping("GetTraza", c.GetTraza)
	c.Mapping("GetIngresos", c.GetIngresos)
	c.Mapping("IngresosSedes", c.IngresosSedes)
}

// GetTraza ...
// @Title GetTraza
// @Description Devuelve la traza del estudiante
// @Param	id		path 	int	true		"Id del estudiante del que se quiere obtener la traza"
// @Success 200 {object} []models.TrazaEstudiante
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Control_datosController) GetTraza() {
	idEstudiante := c.GetString(":id")

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/GetTraza/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	_, err := strconv.Atoi(idEstudiante)
	if err != nil {
		panic(map[string]interface{}{"funcion": "GetTraza", "err": "Error parametro de ingreso \"id\"", "status": "400", "log": err})
	}

	if respuesta, err := helpers.ConsultarTraza(idEstudiante); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta}
	} else {
		panic(err)
	}

	c.ServeJSON()
}

// GetIngresos ...
// @Title GetIngresos
// @Description Devuelve la cantidad de ingresos que tuvo un espacio
// @Param	idEspacio		path 	int		true		"Id del espacio"
// @Param	fecha			path 	int		true		"Fecha de la que se quieren consultar los ingresos (AAAAMMDD)"
// @Success 200 {object} map[string]interface{}
// @Failure 403 :id is empty
// @router /:idEspacio/:fecha [get]
func (c *Control_datosController) GetIngresos() {
	idEspacio := c.GetString(":idEspacio")
	fecha := c.GetString(":fecha")
	fecha = fecha[:4] + "-" + fecha[4:6] + "-" + fecha[6:]

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/GetTraza/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	_, err := strconv.Atoi(idEspacio)
	if err != nil {
		panic(map[string]interface{}{"funcion": "GetTraza", "err": "Error parametro de ingreso \"idEspacio\"", "status": "400", "log": err})
	}
	_, err = time.Parse("2006-01-02", fecha)
	if err != nil {
		panic(map[string]interface{}{"funcion": "GetTraza", "err": "Error parametro de ingreso \"fecha\"", "status": "400", "log": err})
	}

	if respuesta, err := helpers.ContarIngresos(fecha, idEspacio); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta}
	} else {
		panic(err)
	}

	c.ServeJSON()
}

// Ingresos por sede ...
// @Title Ingresos por sede
// @Description Devuelve la cantidad de ingresos de cada sede en la fecha indicada
// @Param	dia			path 	string		true		"Día en formato DD"
// @Param	mes			path 	string		true		"Mes en formato MM"
// @Param	anio		path 	string		true		"Año en formato AAAA"
// @Success 200 {object} map[string]interface{}
// @Failure 403 fecha is empty
// @router /:dia/:mes/:anio [get]
func (c *Control_datosController) IngresosSedes() {
	dia := c.GetString(":dia")
	mes := c.GetString(":mes")
	año := c.GetString(":anio")
	fecha := año + "-" + mes + "-" + dia

	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["message"] = (beego.AppConfig.String("appname") + "/IngresosSedes/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("404")
			}
		}
	}()

	_, err := strconv.Atoi(dia)
	if err != nil {
		panic(map[string]interface{}{"funcion": "GetTraza", "err": "Error parametro de ingreso \"dia\"", "status": "400", "log": err})
	}
	_, err = strconv.Atoi(mes)
	if err != nil {
		panic(map[string]interface{}{"funcion": "GetTraza", "err": "Error parametro de ingreso \"mes\"", "status": "400", "log": err})
	}
	_, err = strconv.Atoi(año)
	if err != nil {
		panic(map[string]interface{}{"funcion": "GetTraza", "err": "Error parametro de ingreso \"año\"", "status": "400", "log": err})
	}
	_, err = time.Parse("2006-01-02", fecha)
	if err != nil {
		panic(map[string]interface{}{"funcion": "GetTraza", "err": "La fecha ingresada es invalida", "status": "400", "log": err})
	}

	if respuesta, err := helpers.IngresosPorSedes(fecha); err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": respuesta}
	} else {
		panic(err)
	}

	c.ServeJSON()
}
