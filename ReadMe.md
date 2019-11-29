# Ad Mediation

This repository contains Ad Mediation _service_, written in `go`, which provides list of _ad networks_ through `REST API`.

The service has two functionalities:

1. Retrieve list of ad networks
2. Update list of ad networks

Idea is to design and build backend system to provide ad network list to the mobile devices. Data batch processing service will update the list accordingly and Ad Mediation service will provide the latest and newest data to the mobile devices.

Image below shows how will system operate:
![Overview](readme-files/adMediation-Overview.png?raw=true "Designed backend system")

From the picture we can tell that system must be designed in a scalable manner to handle all the request from the mobile devices. Final solution will be deployed to the Google Cloud, with that in mind, we need to use appropriate tools to enable easy and efficient deployment.

## System Design

I have designed the system using Domain Driven Design approach, where the solution consists of the following building block: storage, object cache, services and endpoint and finally main service to run the whole app.

Sketch of the described system is shown below:
![Design](readme-files/adMediation-Implementation.png?raw=true "Designed backend system")

### API endpoints

API has three endpoints:

1. GET /api/v1/adNetworkList
2. PUT /api/v1/adNetworkList
3. GET /swagger/index.html

First two endpoints provides the desired functionality of the service, while last one provides documentation page where developers can try the service functionality.

### Services and storage

Requests to the first two endpoint in previous section are processed by the adding or listing service. Those two services takes care that data is persisted and cached to enable fast data access and cost savings. 

Object cache provides in memory cache where all the objects are cached with time to live of 5 minutes. This means that cached objects will be removed from cache in 5 minutes from adding. If some document is updated we also delete the cached previously cached object to prevent serving stale data. Off course we need to be aware that this cache is per service and when we run the application in multiple nodes, some nodes will serve stale data. We need to be careful and configure TTL time, if we want to prevent serving stale data we need to implement distributed cache.

Persistent storage layer offers two implementations: MongoDB and Cloud Storage. I have implemented MongoDB storage for the sake of development and debugging the application. For the production environment is desirable to use Cloud Storage.

## Running and Deployment

You can run this application without or using Docker.

Also you need to choose what storage you want to use to persist data.
Using configuration you can switch between MongoDB or Google Cloud Storage.

### Without Docker

To run the application on your PC without Docker you need to have installed `go SDK` and `MongoDB`.

Once all the requirements are installed execute the following commands to run the service:

```bash
go mod init github.com/kosenina/ad-mediation
go build
./ad-mediation
```

The first command specify to use Go modules for dependency management, where the second one builds the project and the last one runs the builded executable.

Now the application is up and running and you can try it using Swagger on URL http://localhost:8080/swagger/index.html.

### Using Docker

If you choose to use docker, I recommend you to configure the application to use Cloud Storage, otherwise use docker compose, because it is more convenient as to defining MongoDB configuration.

To use docker, follow the instructions bellow:

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

### Using Docker compose

If you don't want to install dependencies (Go, MongoDB) than Docker compose is the right choice.
You need to have installed Docker and docker-compose, than you need to execute the following command to run the service:

```bash
docker-compose up -d --build
```

If you want to stop the service, execute:

```bash
docker-compose down
```

## Service Functionality

This backend system is REST API server which expose one endpoint: `http://localhost:8080/api/v1/adNetworkList`, where only HTTP GET and PUT actions are supported.

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
