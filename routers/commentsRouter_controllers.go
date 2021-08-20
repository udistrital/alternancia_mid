package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/alternancia_mid/controllers:AccesoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/alternancia_mid/controllers:AccesoController"],
        beego.ControllerComments{
            Method: "GetAutorizacion",
            Router: `/:idQr/:idScan/:tipo`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/alternancia_mid/controllers:AccesoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/alternancia_mid/controllers:AccesoController"],
        beego.ControllerComments{
            Method: "GetIngreso",
            Router: `/idEspacio/tipoQR`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
