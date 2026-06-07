-- +goose Up
-- +goose StatementBegin
CREATE TABLE raw_orders (
   id             TEXT NOT NULL,
   lastName       TEXT NOT NULL,
   firstName      TEXT NOT NULL,
   phoneNumber    TEXT NOT NULL,
   email          TEXT NOT NULL,
   amountSizeS    INTEGER  NOT NULL,
   amountSizeM    INTEGER  NOT NULL,
   amountSizeL    INTEGER  NOT NULL,
   address        TEXT NOT NULL,
   city           TEXT NOT NULL,
   postalCode     TEXT NOT NULL,
   deliveryMethod TEXT NOT NULL,
   remarks        TEXT,
   orderMethod    TEXT,
   vendorGroup    TEXT,
   PRIMARY KEY(id)
);  
CREATE TABLE order_coordinates (
   order_id TEXT NOT NULL REFERENCES raw_orders(id),
   x        REAL NOT NULL,
   y        REAL NOT NULL,
   label    TEXT NOT NULL,
   PRIMARY KEY(order_id)
);
CREATE TABLE geocode_errors (
   order_id           TEXT PRIMARY KEY NOT NULL  REFERENCES raw_orders(id),
   message            TEXT NOT NULL,
   suggested_northing REAL,
   suggested_easting  REAL,
   suggested_label    TEXT 
);
CREATE TABLE delivery (
   id    INTEGER PRIMARY KEY AUTOINCREMENT,
   label TEXT    NOT NULL,
   until TEXT    NOT NULL,
   color TEXT    NOT NULL
);
CREATE TABLE delivery_orders (
   order_id     TEXT    REFERENCES raw_orders(id) NOT NULL,
   delivery_id  INTEGER NOT NULL    REFERENCES delivery(id),
   PRIMARY KEY(order_id)
);
CREATE TABLE avz(
  ADR_EASTING REAL NOT NULL,
  ADR_NORTHING REAL NOT NULL,
  ADR_NUMBER TEXT NOT NULL,
  COM_CANTON TEXT NOT NULL,
  COM_NAME TEXT NOT NULL,
  STN_LABEL TEXT NOT NULL,
  ZIP_LABEL TEXT NOT NULL,
   PRIMARY KEY(STN_LABEL,ADR_NUMBER,ZIP_LABEL)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE raw_orders;
DROP TABLE order_coordinates;
DROP TABLE geocode_errors;
DROP TABLE delivery;
DROP TABLE delivery_orders;
DROP TABLE avz;
-- +goose StatementEnd

