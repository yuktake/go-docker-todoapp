-- Create "users" table
CREATE TABLE "public"."users" ("id" bigserial NOT NULL, "name" character varying NOT NULL, "password" character varying NOT NULL, "email" character varying NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"));
-- Create index "user_email_idx" to table: "users"
CREATE UNIQUE INDEX "user_email_idx" ON "public"."users" ("email");
