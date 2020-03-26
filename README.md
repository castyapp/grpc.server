# Casty gRPC Server
* This is a gRPC server project written in go!

* **What is gRPC and why we're using it?** According to [gRPC official website](https://grpc.io/): <br/> gRPC is a modern open source high performance RPC framework that can run in any environment. It can efficiently connect services in and across data centers with pluggable support for load balancing, tracing, health checking and authentication. It is also applicable in last mile of distributed computing to connect devices, mobile applications and browsers to backend services.

## Prerequisites

* First, ensure that you have installed Go 1.11 or higher since we need the support for Go modules via go mod. [Go modules via `go mod`](https://github.com/golang/go/wiki/Modules)

* mongodb **This project uses mongodb as database connection!**  [Mongodb official website](https://www.mongodb.com/)

## Clone the project
```bash
$ git clone https://github.com/CastyLab/grpc.server.git
```

## Configuraition
Make a copy of `.env.example` for your own configuration. save it as `.env` file.
```bash
$ cp .env.example .env
```

## Environments
### Mongodb configuration
Put your mongodb connection here
```env
DB_HOST=localhost
DB_PORT=27017
DB_USER=root
DB_PASS=secret
DB_NAME=casty
```

### JWT configuration
We use JWT for our authentication method
```env
JWT_EXPIRE_TIME=60 # in minutes
JWT_REFRESH_TOKEN_VALID_TIME=7 # in days
JWT_PRIVATE_KEY_PATH=./jwt/keys/app.key
JWT_PUBLIC_KEY_PATH=./jwt/keys/app.key.pub
```

to generate jwt [public/private] keys you can use `ssh-keygen`
```bash
$ ssh-keygen -t rsa -N '' -b 4096 -m PEM -f ./jwt/keys/app.key &&\
    openssl rsa -in ./jwt/keys/app.key -pubout -outform PEM &&\
    -out ./jwt/keys/app.key.pub;
``` 

### Other environments
```env
# Storage path is used for upload avatars, banners etc...
# This environment is useful for shared volumes between containers
STORAGE_PATH=

# Sentry DSN path *optional
SENTRY_DSN=
```

You're ready to Go!

## Run project with go compiler
you can simply run the project with following command
* this command with install dependencies and after that will run the project
* this project uses go mod file, You can run this project out of the $GOPAH file!
```bash
$ go run server.go
```

or if you're considering building the project use
```bash
$ go build -o server .
```

### or you can [Build/Run] docker image
```bash
$ docker build . --tag=casty.grpc

$ docker run -dp --restart=always 55283:55283 casty.grpc
```

## Contributing
Thank you for considering contributing to this project!