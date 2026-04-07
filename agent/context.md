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
    client_id       BIGINT          NOT NULL,
    issuance_date   DATE            NOT NULL,
    expiration_date DATE            NOT NULL,
    total_cop       DECIMAL(18, 2)  NOT NULL,
    PRIMARY KEY (id)
    );