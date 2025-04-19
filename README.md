# AMDP Task

The project currently implements the 5<sub>th</sub> to 9<sub>th</sub> steps of the [task.pdf](./task.pdf), as well as the 15<sub>th</sub>.

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
KAFKA_BROKERS=kafka:9093
KAFKA_TOPIC_OBJECTS_INCOMING=objects.incoming
KAFKA_TOPIC_OBJECTS_OUTGOING=objects.outgoing
KAFKA_TOPIC_OBJECTS_ERRORS=objects.errors

# Objects Service
LOG_ENVIRONMENT=Development
LOG_FILE_ENC_TYPE=Console
LOG_FILE_LVLF=Gte
LOG_FILE_LVL=INFO
LOG_FILE_PATH=/amdp-task/services/objects/logs/app.log
LOG_KAFKA_ENC_TYPE=JSON
LOG_KAFKA_LVLF=Gte
LOG_KAFKA_LVL=ERROR
LOG_KAFKA_TOPIC=objects.errors
HTTP_SERVER_ADDR=:3000
POSTGRES_DSN=postgres://postgres:postgres_1234@postgres:5432/amdp_objects?sslmode=disable
MESSAGE_POLL_INTERVAL=1s
MESSAGE_POLL_SIZE=50
SHUTDOWN_TIMEOUT=10s

```

The `objects` service will restart a few times until the `kafka` service is full up. Once it's up,
you can send requests to `PUT localhost:3000/api/v1/objects/` and`GET localhost:3000/api/v1/objects/{objectid}/meta/`. The former expects a `form-data` body with a single `file` field, the latter expects a route path parameter which is the `uuid` of the uploaded object. You can use the [postman_collection.json](.postman_collection.json) as reference for sending the HTTP requests.

The log files will be stored at `/var/amdp-task/services/objects/logs` and also the
`objects.errors` topics. Use `tail` and `kcat` to view them.

```sh
tail -f /var/amdp-task/services/objects/logs/app.log
```

```sh
kcat -C -b localhost:9092 -t objects.errors
```

The result of the `PUT` request will first be stored in an outbox table called `messages`, and then
periodically queried and published to `objects.outgoing` of the `kafka` broker. You can use the following
command to observe the messages.

```sh
kcat -C -b localhost:9092 -t objects.outgoing
```

You can also visit `MinIO` control panel at `localhost:9000`. The username and password are set by the
environment variables `MINIO_ROOT_USER` and `MINIO_ROOT_PASSWORD`.
