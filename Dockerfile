FROM golang:1.14

LABEL maintainer="Alireza Josheghani <josheghani.dev@gmail.com>"

ARG DEBIAN_FRONTEND=noninteractive

# Update and install curl
RUN apt-get update &&\
    apt-get -y install openssh-client &&\
    apt-get -y install nano ffmpeg

RUN ffprobe -version

# Creating work directory
RUN mkdir /code

# Adding project to work directory
ADD . /code

RUN mkdir /config

# Removing old JWT keys
RUN rm -rf /config/jwt.key /config/jwt.pub

# Generate jwt keys
RUN cd /config && ssh-keygen -t rsa -N '' -b 4096 -m PEM -f jwt.key &&\
    openssl rsa -in jwt.key -pubout -outform PEM -out jwt.pub;

# Choosing work directory
WORKDIR /code

# build project
RUN go build -o casty.gRPC.server .

EXPOSE 55283

ENTRYPOINT ["./casty.gRPC.server"]
CMD ["--port", "55283"]