*(Done.)*
</think>

# gow
gow es un framework web HTTP para Go

![gow logo](docs/logo.png)


> Proyectos de referencia e inspiración: gin/beego/mux

## Repositorio del proyecto

[https://github.com/zituocn/gow](https://github.com/zituocn/gow)

## Características

* Encapsulación de `Context`, agrupación de rutas y middleware similar a `gin`, fácil de aprender
* Coincidencia exacta de rutas mediante `regexp`, con soporte para ignorar mayúsculas y minúsculas
* Punto de entrada unificado para configuraciones (formato ini), pero permite implementar tu propio formato
* Soporte para renderizado de páginas HTML en el servidor, con funciones de plantilla HTML extensibles
* Opción libre de elegir SDKs encapsulados en lib, como mysql, redis, nsq, rpc, mem-cache, oauth, pay, etc.

## 1. Inicio rápido

```sh
# Crea un proyecto llamadohello
mkdir hello

cd hello

# Inicializa go mod
go mod init

# Instala gow

go get github.com/zituocn/gow
```

### 1.1 Crear main.go

```go
package main

import (
	"github.com/zituocn/gow"
)

func main() {
	r := gow.Default()

	r.GET("/", func(c *gow.Context) {
		c.JSON(gow.H{
			"code": 0,
			"msg":  "success",
		})
	})

	// default :8080
	r.Run()
}
```

También puedes escribirlo así:

```go
package main

import (
	"github.com/zituocn/gow"
)

func main() {
	r := gow.Default()
	r.GET("/", IndexHandler)
	// default :8080
	r.Run()
}

// IndexHandler response h
func IndexHandler(c *gow.Context) {
	h := map[string]interface{}{
		"project": "gow",
		"website": "https://github.com/zituocn/gow",
	}
	c.JSON(h)
}
```

Para métodos personalizados, utiliza `r.Handle`:

```go
package main

import "github.com/zituocn/gow"

func main() {
	r := gow.Default()
	r.Handle("GET,POST", "/", func(c *gow.Context) {
		c.String("index")
	})
}
```

### 1.2 Ejecutar

```sh
go run main.go
```

Resultado de la ejecución:

```sh
Listening and serving HTTP on http://127.0.0.1:8080
```

### 1.3 Acceso

*Acceso con curl*

```sh
curl -i http://127.0.0.1:8080
```

Resultado de la solicitud:

```sh
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 08 Jun 2021 08:51:25 GMT
Content-Length: 67

{
  "project": "gow",
  "website": "https://github.com/zituocn/gow"
}
```

Acceso desde el navegador:

```sh
Accede desde el navegador a: http://127.0.0.1:8080
```

---

## Algunos códigos de demostración

Ejecutables directamente:

* [github.com/zituocn/gow-demo](https://github.com/zituocn/gow-demo)

---

## 2. Más documentación

* [Detalle de rutas && Parámetros de ruta && Agrupación de rutas](https://github.com/zituocn/gow/blob/main/docs/route.md)
* [Uso de middleware](https://github.com/zituocn/gow/blob/main/docs/middleware.md)
* [Obtener valores de la solicitud](https://github.com/zituocn/gow/blob/main/docs/request.md)
* [Valores de respuesta && JSON / XML / JSONP / YAML](https://github.com/zituocn/gow/blob/main/docs/response.md)
* [Archivo de configuración unificado](https://github.com/zituocn/gow/blob/main/docs/config.md)
* [Crear un sitio web && Guía de uso de plantillas HTML](https://github.com/zituocn/gow/blob/main/docs/website.md)
* [Funciones de plantillas HTML](https://github.com/zituocn/gow/blob/main/docs/html.md)
* [Subida y descarga de archivos](https://github.com/zituocn/gow/blob/main/docs/upload.md)
* [Implementar proxy inverso (nuevo)](https://github.com/zituocn/gow/blob/main/docs/proxy.md)
* [Introducción a la librería lib: logx, mysql, config](https://github.com/zituocn/logx)

## 3. Agradecimientos

* [beego](https://github.com/beego/beego) -> Referencia para el diseño de plantillas HTML en la versión 1.x
* [gin](https://github.com/gin-gonic/gin) -> Referencia para el diseño de `engine` y `Context`
* [mux](https://github.com/gorilla/mux)   -> Referencia para el diseño de enrutamiento
* [gorm](https://github.com/go-gorm/gorm) -> Recomendado para el uso de ORM
* [gini](https://github.com/gkzy/gini)     -> Utilizado para manipular archivos de configuración en formato `ini`

## 4. Licencia

Licencia MIT. Consulta el archivo LICENSE para más detalles.
