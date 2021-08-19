# alternancia_mid

API MID del Sistema de Alternancia

## Especificaciones Técnicas

### Tecnologías Implementadas y Versiones

- [Golang](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/golang.md)
- [BeeGo](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/beego.md)

### Variables de Entorno

```shell
ALTERNANCIA_MID_HTTP_PORT=[puerto de ejecución]
TERCEROS_SERVICE=[servicio api terceros]
SINTOMAS_SERVICE=[servicio api sintomas]
OIKOS_SERVICE=[servicio api oikos]
```

### Ejecución del Proyecto 

```shell
#1. Obtener el repositorio con Go
go get github.com/udistrital/alternancia_mid

#2. Moverse a la carpeta del repositorio
cd $GOPATH/src/github.com/udistrital/alternancia_mid

# 3. Moverse a la rama **develop**
git pull origin develop && git checkout develop

# 4. alimentar todas las variables de entorno que utiliza el proyecto.
ALTERNANCIA_MID_HTTP_PORT=8080 ALTERNANCIA_MID_SOME_VARIABLE=some_value bee run
```

### Ejecución Dockerfile

### Ejecución docker-compose

### Ejecución Pruebas

Pruebas unitarias

```shell
# En Proceso
```

## Estado CI 

## Licencia

This file is part of alternancia_mid.

alternancia_mid is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

alternancia_mid is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with novedades_crud. If not, see https://www.gnu.org/licenses/.
