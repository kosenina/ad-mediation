# Docker deployment

Use Go modules for dependency management.

```bash
go mod init github.com/kosenina/ad-mediation
```

Firstly build and run our application locally. Execute the following commands and create HTTP GET request to the http://localhost:8080/api/v1/adNetworkList, and you should get response saying that document does not exists.

```bash
go build
./ad-mediation
```

Build and tag docker image.

```bash
docker build -t ad-mediation-docker .
```

 List all the available images and make sure that `ad-mediation-docker` image is on the list.

```bash
docker image ls
```

Than run docker image.

```bash
docker run -d -p 8080:8080 ad-mediation-docker
```

List all the running containers and make sure that our image is up and running.

```bash
docker container ls
```