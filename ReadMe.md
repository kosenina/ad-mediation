# How to run app

You can run this application without or using Docker.

## Without Docker

To run the application on your PC you need to have installed go SDK and MongoDB.

Once all the requirements are installed execute the following commands:

```bash
go mod init github.com/kosenina/ad-mediation
go build
./ad-mediation
```

The first command specified to use Go modules for dependency management, where the second one builds the project and the last one runs the builded executable.

Now the application is up and running.
HTTP server is exposing one API endpoint: `http://localhost:8080/api/v1/adNetworkList`, where only HTTP GET and PUT actions are supported.

### Retrieve the list of ad networks

To get the list of ad networks you need to create `HTTP GET` request to the specified `URL`.
Use `curl` or any other utility to create this kind of request.
At first, the list does not exists and the API will return `JSON` object saying that ad network list does not exists.
To retrieve the ad network list we need to create one.

### Update or insert the list of ad networks

To update or insert one need to create a `HTTP PUT` request to the specified `URL` where request body consist of properly formated list of ad networks.
Example of request body is the following `JSON` object:

```json
{
  "items": [
      {
          "name": "AdX",
          "rank": 0
      },
      {
          "name": "AdMob",
          "rank": 1
      },
      {
          "name": "Ad Unity",
          "rank": 2
      }
  ]
}
```

If provided request body is not in the right format or is not valid the API will properly return `BadRequest`. Valid body must be in the right format, where ad network name is not empty and its rank are integer value (greater or equal to 0). Also, all the ranks must be in a sequence (0, 1, 2, ...).

## Using Docker

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
