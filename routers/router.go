// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/udistrital/alternancia_mid/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/acceso",
			beego.NSInclude(
				&controllers.AccesoController{},
			),
		),
		beego.NSNamespace("/control_datos",
			beego.NSInclude(
				&controllers.Control_datosController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
