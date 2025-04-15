CREATE DATABASE "amdp_objects";

\c amdp_objects;

CREATE SCHEMA IF NOT EXISTS "app";

CREATE TABLE IF NOT EXISTS "app"."objects" (
    "id" BIGSERIAL PRIMARY KEY,               
    "oid" UUID NOT NULL,                   
    "user_id" BIGINT NOT NULL,             
    "name" VARCHAR(255) NOT NULL,          
    "mime" VARCHAR(50) NOT NULL,                    
    "size" BIGINT NOT NULL,                
    "hash" VARCHAR(255) NOT NULL,                   
    "state" VARCHAR(50) NOT NULL,                   
    "created_at" TIMESTAMPTZ NOT NULL,  
    CONSTRAINT uq_objects_oid UNIQUE("oid"),
    CONSTRAINT uq_objects_uid_name_mime UNIQUE("user_id", "name", "mime")   
);

CREATE TABLE IF NOT EXISTS "app"."outbox" (
    "id" BIGSERIAL PRIMARY KEY,
    "uuid" UUID NOT NULL,
    "entity" VARCHAR(50) NOT NULL,
    "entity_id" BIGINT NOT NULL,
    "title" VARCHAR(50) NOT NULL,
    "target" VARCHAR(50) NOT NULL,
    "payload" JSONB NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "processed_at" TIMESTAMPTZ NULL
);
