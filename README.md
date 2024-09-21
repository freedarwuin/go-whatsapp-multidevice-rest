# Implementación de WhatsApp para múltiples dispositivos en API REST

Este repositorio contiene ejemplos de implementación. [go.mau.fi/whatsmeow](https://go.mau.fi/whatsmeow/) paquete con compatibilidad con múltiples sesiones y cuentas. Este ejemplo utiliza un [labstack/echo](https://github.com/labstack/echo) versión 4.x.

## Características

- Compatibilidad con múltiples sesiones y cuentas
- Compatibilidad con múltiples dispositivos
- Autenticación de WhatsApp (código QR y cierre de sesión)
- Envío de mensajes de WhatsApp
- Envío de archivos multimedia (documentos, imágenes, audio, videos, stickers) a través de mensajes de WhatsApp
- Envío de ubicación a través de mensajes de WhatsApp
- Envío de contactos a través de mensajes de WhatsApp
- Envío de enlaces a través de mensajes de WhatsApp
- Y mucho más...

## Primeros pasos

Estas instrucciones le permitirán obtener una copia del proyecto en funcionamiento en su máquina local para fines de desarrollo y prueba.
Consulte la sección de implementación para obtener notas sobre cómo implementar el proyecto en un sistema en vivo.

### Requisitos previos

Paquetes de requisitos previos:
* Go (lenguaje de programación Go)
* Swag (convertidor de anotaciones Go a documentación Swagger)
* GoReleaser (compilación automática de binarios Go)
* Make (ejecución automatizada mediante Makefile)

Paquetes opcionales:
* Docker (contenedorización de aplicaciones)

### Implementación

#### **Uso de contenedores**

1) Instale Docker CE según el [manual documentation](https://docs.docker.com/desktop/)

2) Ejecute el siguiente comando en su Terminal o PowerShell
```sh
docker run -d \
  -p 3000:3000 \
  --name go-whatsapp-multidevice \
  --rm dimaskiddo/go-whatsapp-multidevice-rest:latest
```

3) Ahora debería ser accesible en su máquina accediendo `localhost:3000/api/v1/whatsapp` or `127.0.0.1:3000/api/v1/whatsapp`

4) Intente utilizar documentos API integrados a los que se pueda acceder en `localhost:3000/api/v1/whatsapp/docs/` or `127.0.0.1:3000/api/v1/whatsapp/docs/`

#### **Uso de binarios prediseñados** 
1) Descargue los binarios prediseñados desde [release page](https://github.com/dimaskiddo/go-whatsapp-multidevice-rest/releases)

2) Extraiga el archivo comprimido

3) Copie el `.env.default` archivar como archivo `.env` 

4) Ejecutar el binario pre-compilado
```sh
# MacOS / Linux
chmod 755 go-whatsapp-multidevice-rest
./go-whatsapp-multidevice-rest

# Windows
# Puede hacer doble clic en él o usar PowerShell
.\go-whatsapp-multidevice-rest.exe
```

5) Ahora debería ser accesible en su máquina accediendo `localhost:3000/api/v1/whatsapp` or `127.0.0.1:3000/api/v1/whatsapp`

6) Intente utilizar documentos API integrados a los que se pueda acceder en `localhost:3000/api/v1/whatsapp/docs/` or `127.0.0.1:3000/api/v1/whatsapp/docs/`

#### **Crear desde el código fuente**

A continuación se muestran las instrucciones para ejecutar este código fuente:

1) Cree un directorio de Go Workspace y expórtelo como el directorio GOPATH extendido
```sh
cd <your_go_workspace_directory>
export GOPATH=$GOPATH:"`pwd`"
```

2) En el directorio Go Workspace crea un directorio de origen
```sh
mkdir -p src/github.com/dimaskiddo/go-whatsapp-multidevice-rest
```

3) Muévete al directorio creado y extrae el código base
```sh
cd src/github.com/dimaskiddo/go-whatsapp-multidevice-rest
git clone -b master https://github.com/dimaskiddo/go-whatsapp-multidevice-rest.git .
```

4) Ejecute el siguiente comando para extraer los paquetes del proveedor
```sh
make vendor
```

5) Vincular o copiar archivo de variables de entorno
```sh
ln -sf .env.development .env
# - OR -
cp .env.development .env
```

6) Hasta este paso ya puedes ejecutar este código usando este comando
```sh
make run
```

7) *(Opcional)* Utilice el siguiente comando para crear este código en una plataforma binaria específica
```sh
make build
```

8) *(Opcional)* Para realizar una distribución masiva de binarios puedes utilizar el siguiente comando
```sh
make release
```

9) Ahora debería ser accesible en su máquina accediendo `localhost:3000/api/v1/whatsapp` o `127.0.0.1:3000/api/v1/whatsapp`

10) Intente utilizar documentos API integrados a los que se pueda acceder en `localhost:3000/api/v1/whatsapp/docs/` o `127.0.0.1:3000/api/v1/whatsapp/docs/`

## Acceso a la API

Puede acceder a cualquier punto final mediante la variable de entorno **HTTP_BASE_URL**, que de forma predeterminada se encuentra en el archivo `.env`.

Se puede acceder a la documentación API integrada en `<HTTP_BASE_URL>/docs/` o por defecto esta en `localhost:3000/api/v1/whatsapp/docs/` o `127.0.0.1:3000/api/v1/whatsapp/docs/`

## Ejecución de las pruebas

Actualmente, la prueba aún no está lista :)

## Creado con

* [Go](https://golang.org/) - Lenguaje de programación Go
* [Swag](https://github.com/swaggo/swag) - Documentación del convertidor de anotaciones Go a Swagger
* [GoReleaser](https://github.com/goreleaser/goreleaser) - Creación automática de binarios
* [Make](https://www.gnu.org/software/make/) - Ejecución automatizada de GNU Make
* [Docker](https://www.docker.com/) - Contenerización de aplicaciones

## Autores

* **Dimas Restu Hidayanto** - *Trabajo Inicial* - [DimasKiddo](https://github.com/dimaskiddo)

Véase también la lista de [Contribuyentes](https://github.com/dimaskiddo/go-whatsapp-multidevice-rest/contributors) who participated in this project

## Anotación

Puede buscar más información sobre los parámetros del comando make en [Makefile](https://github.com/dimaskiddo/go-whatsapp-multidevice-rest/-/raw/master/Makefile)

## Licencia

Copyright (C) 2022 Dimas Restu Hidayanto

Se otorga permiso, sin cargo, a cualquier persona que obtenga una copia de este software y los archivos de documentación asociados (el "Software") para utilizar el Software sin restricciones, incluidos, entre otros, los derechos de uso, copia, modificación, fusión, publicación, distribución, sublicencia y/o venta de copias del Software, y para permitir que las personas a quienes se les proporciona el Software lo hagan, sujeto a las siguientes condiciones:

El aviso de copyright anterior y este aviso de permiso se incluirán en todas las copias o partes sustanciales del Software.

EL SOFTWARE SE PROPORCIONA "TAL CUAL", SIN GARANTÍA DE NINGÚN TIPO, EXPRESA O IMPLÍCITA, INCLUYENDO, ENTRE OTRAS, LAS GARANTÍAS DE COMERCIABILIDAD, IDONEIDAD PARA UN PROPÓSITO PARTICULAR Y NO INFRACCIÓN. EN NINGÚN CASO LOS AUTORES O TITULARES DE LOS DERECHOS DE AUTOR SERÁN RESPONSABLES DE CUALQUIER RECLAMO, DAÑO U OTRA RESPONSABILIDAD, YA SEA EN UNA ACCIÓN CONTRACTUAL, AGRAVIO O DE OTRO MODO, QUE SURJA DE, A PARTIR DE O EN CONEXIÓN CON EL SOFTWARE O EL USO U OTRAS OPERACIONES EN EL SOFTWARE.
