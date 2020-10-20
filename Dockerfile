FROM golang:1.14

LABEL maintainer="Alireza Josheghani <josheghani.dev@gmail.com>"

# Creating work directory
WORKDIR /app

# Adding project to work directory
ADD . /app

# build project
RUN go build -o server .

EXPOSE 55283

ENTRYPOINT ["/app/server"]
CMD ["--port", "55283"]