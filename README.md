## Build the application image

```sh
./build-image.sh
```

## Start the stack locally

```sh
docker-compose up -d
```

## Testing using cURL

```sh
curl -i http://localhost:8080/hello
curl -i http://localhost:8080/version
curl -i http://localhost:8080/greetings

curl -i -X POST \
  -H "Content-type: application/json" \
  -d '{"message":"Hola"}' \
  http://localhost:8080/greetings
```

## Push image to [Docker Hub](https://hub.docker.com/)

```sh
docker login
docker push polybean/hello-go:3.0
```

## Clean up

```sh
docker-compose down
```
