package controllers

import (
	"strconv"

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
