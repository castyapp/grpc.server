FROM golang:1.14

LABEL maintainer="Alireza Josheghani <josheghani.dev@gmail.com>"

ARG DEBIAN_FRONTEND=noninteractive

# Creating work directory
WORKDIR /code

# Adding project to work directory
ADD . /code

RUN mkdir /config

# build project
RUN go build -o server .

EXPOSE 55283

ENTRYPOINT ["/code/server"]
CMD ["--port", "55283"]