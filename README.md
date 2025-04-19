# AMDP Task

## Getting Started

Copy the following `.env` file in the root folder and then run `docker compose up`. You can modify
the variables to your preference.

```sh
# PostgreSQL
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres_1234

# MinIO
MINIO_USE_SSL=false
MINIO_ACCESS_KEY=edf02c8c-fb8f-4219-8956-78900215958b
MINIO_SECRET_KEY=f02ca9f1-7a82-432b-9184-c578a1f02184
MINIO_ROOT_USER=minio_root
MINIO_ROOT_PASSWORD=adminm_1234
MINIO_BUCKET_NAME=default
MINIO_ENDPOINT=minio:9000

# Kafka
KAFKA_BROKERS=kafka:9092
KAFKA_TOPIC_OBJECTS_INCOMING=objects.incoming
KAFKA_TOPIC_OBJECTS_OUTGOING=objects.outgoing
KAFKA_TOPIC_OBJECTS_ERRORS=objects.errors

# Objects Service
LOG_ENVIRONMENT=Development
LOG_FILE_ENC_TYPE=Console
LOG_FILE_LVLF=Gte
LOG_FILE_LVL=INFO
LOG_FILE_PATH=/amdp-task/services/objects/logs/data.log
LOG_KAFKA_ENC_TYPE=JSON
LOG_KAFKA_LVLF=Gte
LOG_KAFKA_LVL=ERROR
LOG_KAFKA_TOPIC=objects.errors
HTTP_SERVER_ADDR=:3000
POSTGRES_DSN=postgres://postgres:postgres_1234@postgres:5432/amdp_objects?sslmode=disable

```
