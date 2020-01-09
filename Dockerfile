FROM golang:1.13

LABEL maintainer="Alireza Josheghani <josheghani.dev@gmail.com>"

ARG GITLAB_ACCESS_TOKEN
ARG DEBIAN_FRONTEND=noninteractive

# Update and install curl
RUN apt-get update &&\
    apt-get -y install openssh-client

RUN git config --global url."https://oauth2:${GITLAB_ACCESS_TOKEN}@gitlab.com/".insteadOf "https://gitlab.com/"

# Creating work directory
RUN mkdir /code

# Adding project to work directory
ADD . /code

# Generate jwt keys
RUN cd /code/jwt/keys && ssh-keygen -t rsa -N '' -b 4096 -m PEM -f app.key &&\
    openssl rsa -in app.key -pubout -outform PEM -out app.key.pub;

# Choosing work directory
WORKDIR /code

# build project
RUN go build -o movie.night.gRPC.server .

CMD ["./movie.night.gRPC.server", "-port", "55283"]