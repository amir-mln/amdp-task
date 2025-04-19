CREATE DATABASE "amdp_objects";

\c amdp_objects;

CREATE SCHEMA IF NOT EXISTS "app";

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM "pg_type"
        WHERE "typname" = 'object_state'
        AND "typnamespace" = ( SELECT "oid" FROM "pg_namespace" WHERE "nspname" = 'app' )
    ) THEN CREATE TYPE "app"."object_state" AS ENUM ( 'Initial', 'Completed', 'Failed' );
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS "app"."objects" (
    "id" BIGSERIAL PRIMARY KEY,               
    "oid" UUID NOT NULL,                   
    "user_id" BIGINT NOT NULL,             
    "name" VARCHAR(256) NOT NULL,          
    "mime" VARCHAR(64) NOT NULL,                    
    "size" BIGINT NOT NULL,                
    "hash" VARCHAR(256) NOT NULL,                   
    "state" "app"."object_state" NOT NULL,                   
    "created_at" TIMESTAMPTZ NOT NULL,  
    CONSTRAINT uq_objects_oid UNIQUE("oid"),
    CONSTRAINT uq_objects_uid_name_mime UNIQUE("user_id", "name", "mime")   
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM "pg_type"
        WHERE "typname" = 'message_type'
        AND "typnamespace" = ( SELECT "oid" FROM "pg_namespace" WHERE "nspname" = 'app' )
    ) THEN CREATE TYPE "app"."message_type" AS ENUM ( 'Event', 'Command', 'Query' );
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS "app"."messages" (
    "id" UUID PRIMARY KEY,
    "user_id" BIGINT NULL,
    "entity" VARCHAR(64) NULL,
    "entity_id" BIGINT NULL,
    "title" VARCHAR(64) NOT NULL,
    "type" "app"."message_type" NOT NULL,
    "payload" JSONB NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "publish_at" TIMESTAMPTZ NULL
);
