# Casty gRPC Server
* This is a gRPC server project written in go!

* **What is gRPC and why we're using it?** According to [gRPC official website](https://grpc.io/): <br/> gRPC is a modern open source high performance RPC framework that can run in any environment. It can efficiently connect services in and across data centers with pluggable support for load balancing, tracing, health checking and authentication. It is also applicable in last mile of distributed computing to connect devices, mobile applications and browsers to backend services.


## Prerequisites
* First, ensure that you have installed Go 1.15 or higher
* mongodb **This project uses mongodb!**  [Mongodb official website](https://www.mongodb.com/)

## Pull Docker Image
```bash
$ docker pull castyapp/grpc:latest
```

## Run docker container
```bash
$ docker run -p 55283:55283 castyapp/grpc
```

## Docker-Compose example
```yaml
version: '3'

services:
  grpc:
    image: castyapp/grpc:latest
    ports:
      - 55283:55283
    args: ['--config-file', '/config/config.hcl']
    volumes:
      - $PWD/config.hcl:/config/config.hcl
```

## Clone the project
```bash
$ git clone https://github.com/castyapp/grpc.server.git
```

## Configuraition
Make a copy of `example.config.hcl` for your own configuration. save it as `config.hcl` in your work directory.
```bash
$ cp example.config.hcl config.hcl
```

## Environments
### Mongodb configuration
Put your mongodb connection here
```hcl
db {
  name        = "casty"
  host        = "localhost"
  port        = 27017
  user        = "service"
  pass        = "super-secure-password"
  auth_source = ""
}
```

### Redis configuration
Put your redis connection here
```hcl
# Redis configurations
redis {
  # if you wish to use redis cluster, set this value to true
  # If cluster is true, sentinels is required
  # If cluster is false, addr is required
  cluster     = false
  master_name = "casty"
  addr        = "127.0.0.1:26379"
  sentinels   = [
    "127.0.0.1:26379"
  ]
  pass = "super-secure-password"
  sentinel_pass = "super-secure-sentinels-password"
}
```

### JWT configuration
We use JWT for our authentication method
```hcl
# JWT secrets
jwt {
  access_token {
    # make sure to use a strong secret key
    secret = "drWRU76y2Pc37TgjD5J8xcWg9e"
    # If you wish to change valid duration of a access_token, change this value
    expires_at {
      type  = "days" # can be [seconds|minutes|hours|days]
      value = 1
    }
  }
  refresh_token {
    # make sure to use a strong secret key
    secret = "S3pXMmmjWFYVPBSLeYdYCve5Ca"
    # If you wish to change valid duration of a refresh_token, change this value
    expires_at {
      type  = "days" # can be [seconds|minutes|hours|days]
      value = 7
    }
  }
}
```

You're ready to Go!

## Run project with go compiler
you can simply run the project with following command
```bash
$ go run server.go
```

or if you're considering building the project
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

## License
Casty is an open-source software licensed under the MIT license.
