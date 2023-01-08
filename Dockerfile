FROM golang:1.12.0-alpine:latest 

RUN MKDIR /app 

ADD . /app 

WORKDIR /app 

RUN env GOOS=linux GOARCH=amd46 go build -o main .

CMD [ "app/main" ]