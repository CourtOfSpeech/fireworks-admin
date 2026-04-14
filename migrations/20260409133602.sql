-- Create "tenants" table
CREATE TABLE "tenants" (
  "id" uuid NOT NULL,
  "status" smallint NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "certificate_no" character varying NOT NULL,
  "name" character varying NOT NULL,
  "type" smallint NOT NULL DEFAULT 1,
  "contact_name" character varying NOT NULL,
  "email" character varying NOT NULL,
  "phone" character varying NOT NULL,
  "expired_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "uk_certificate_no" to table: "tenants"
CREATE UNIQUE INDEX "uk_certificate_no" ON "tenants" ("certificate_no") WHERE (deleted_at IS NULL);
-- Create index "uk_email" to table: "tenants"
CREATE UNIQUE INDEX "uk_email" ON "tenants" ("email") WHERE (deleted_at IS NULL);
-- Create index "uk_phone" to table: "tenants"
CREATE UNIQUE INDEX "uk_phone" ON "tenants" ("phone") WHERE (deleted_at IS NULL);
