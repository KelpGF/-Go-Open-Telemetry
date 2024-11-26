# Go Open Telemetry

## How to run?

```bash
docker-compose up
```

## How to test?

```bash
curl --request POST \
  --url http://localhost:8080/zip-code/validate \
  --header 'Content-Type: application/json' \
  --data '{
  "cep": "01153000"
}'
```

## How to see the metrics?

Access the URL: http://localhost:9411 click in "Run Query"
