# CODEACT-PROJECT
Proyecto de prueba en *golang* del framework ***codeact***, el cual se encarga de resolver solicitudes mediante consulta
 IA (en este caso Claude), el cual analiza, verifica, genera y ejecuta el código necesario para resolver la solicitud
realizada.

## Requisitos
- Go 1.23+

## Instalación
```bash
          
git clone https://github.com/alejoespin/codeact-project.git
cd codeact-project
go mod tidy
```

## Configuración
Los archivos descritos a continuación deben estar en la carpeta ***/agent*** al nivel del archivo ejecutable del proyecto.
```
    - /basePath/
    - /basePath/codeact-project (ejecutable)
    -> /basePath/agent/base_prompt.md
    -> /basePath/agent/context.md
    -> /basePath/agent/configs.env 
```


### Contexto (**agent/context.md**)
El proyecto cuenta con un archivo de configuración el cual describe las librerías, funciones de las mismas o estructuras
necesarias para interactuar con componentes y así establecer contexto al momento de realizar la solicitud a la IA.

Actualmente se establecen algunos valores para el proceso de lectura, interpretación y almacenamiento en base de datos.

```markdown
- Ruta de facturas: ~/./files
- Lib manejo de archivos: os
- - ReadDir, Name
- Lib lectura de PDF: 	github.com/ledongthuc/pdf
- - Open, GetPlainText 
- Consulta de divisas URL: https://open.er-api.com/v6/latest/USD
- Lib para SQL: gorm.io/driver/mysql, gorm.io/gorm
  - Open, Save
- Datos base de datos:
  - path: localhost:3306
  - password: my_passwd
  - database: my_database
  - table: facturas
  - schema: facturas (
    id              VARCHAR(36)     NOT NULL,
    expiration_date DATE            NOT NULL,
    total_cop       DECIMAL(18, 2)  NOT NULL,
    PRIMARY KEY (id)
    );
```

### BASE_PROMPT  (**agent/base_prompt.md**)
En el archivo se encuentra el prompt utilizado en las peticiones al cliente de IA, cuenta con las restricciones y 
condiciones de cada solcitud.

### Credenciales / Valores por defecto  (**configs.env**)
En el archivo de credenciales y valores por defecto agregamos las ***api-keys*** necesarias para conectar
con el cliente IA, el valor maxímo de loop a realizar y la configuración de creación de archivo
temporal de auditoria de código generado.

```text
LOOP-MAX=10
ANTHROPIC_KEY=api-key
AUDIT-RESPONSE=true
```

## Uso
Inicialmente se deben realizar la configuraciones correspondientes (**agent/context.md , agent/configs.env**) una vez 
configurados se debe generar el ejecutable con el siguiente comando.

Para el proceso se genera un archivo ***tmp.go*** con el código a ejecutar el cual es eliminado una vez se realiza la
ejecución del mismo.

```
cd codeact-project
go mod tidy
go build
```

 
Una vez ejecutado con el comando ***./codeact-project*** se vera el siguiente mensaje y el programa esperará la solicitud
```terminaloutput
-> Request:
```
una vez ingresada se generará un mensaje con el texto y la respuesta del proceso.
```terminaloutput
-> Response:
```

```
Ejemplo:
-> Request
¿Cuántos pesos colombianos son 150 dólares hoy?
Response >  549514.42 pesos colombianos son 150 dólares hoy.
```