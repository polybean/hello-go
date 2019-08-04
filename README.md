## Build the application image

```sh
docker-compose build app
```

## Start the stack locally

```sh
docker-compose up -d
```

## Testing using cURL

```sh
curl -i http://localhost:8080/greetings

curl -i -X POST \
  -H "Content-type: application/json" \
  -d '{"message":"Hola"}' \
  http://localhost:8080/greetings
```

## Push image to [Docker Hub](https://hub.docker.com/)

```sh
docker login
docker tag hello-go polybean/hello-go
docker push polybean/hello-go
```
