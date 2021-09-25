package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/alternancia_mid/controllers:AccesoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/alternancia_mid/controllers:AccesoController"],
        beego.ControllerComments{
            Method: "GetIngreso",
            Router: `/:idPersona/:idEspacio/:tipoQR`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/alternancia_mid/controllers:AccesoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/alternancia_mid/controllers:AccesoController"],
        beego.ControllerComments{
            Method: "GetAutorizacion",
            Router: `/:idQr/:idScan`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/alternancia_mid/controllers:Control_datosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/alternancia_mid/controllers:Control_datosController"],
        beego.ControllerComments{
            Method: "IngresosSedes",
            Router: `/:dia/:mes/:anio`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/alternancia_mid/controllers:Control_datosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/alternancia_mid/controllers:Control_datosController"],
        beego.ControllerComments{
            Method: "GetTraza",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/alternancia_mid/controllers:Control_datosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/alternancia_mid/controllers:Control_datosController"],
        beego.ControllerComments{
            Method: "GetIngresos",
            Router: `/:idEspacio/:fecha`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
