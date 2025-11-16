CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE schema IF NOT EXISTS "users";

CREATE TABLE "users"."profiles" (
"id" uuid NOT NULL DEFAULT uuid_generate_v4(),
"created" timestamp not null,
"changed" timestamp not null,
"login" varchar(50) NOT NULL,
"fist_name" varchar(100) NOT NULL,
"middle_name" varchar(100) NOT NULL,
"last_name" varchar(100) NOT NULL,
"private" boolean DEFAULT true,
"primary_email" varchar(50) NOT NULL,
"email" varchar(50)[] DEFAULT array[]::varchar[],
"phone" varchar(50), "tags" varchar(50) DEFAULT array[]::varchar[],
"biography" text, "company" varchar(100),
"location" varchar(255),
"external_id" uuid,
"avatar" uuid,
"metadata" jsonb,
PRIMARY KEY ("id"), UNIQUE ("login"), UNIQUE ("primary_email"));


