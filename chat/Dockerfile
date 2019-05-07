FROM golang
WORKDIR /go/src/github.com/go-park-mail-ru/2019_1_DeathPacito/


RUN go get github.com/gorilla/websocket
RUN go get github.com/dgrijalva/jwt-go
RUN go get go.mongodb.org/mongo-driver/mongo
RUN go get go.mongodb.org/mongo-driver/bson

COPY ./chat/ ./chat/
COPY ./identity_struct ./identity_struct
COPY ./auth ./auth

RUN go build ./chat/server.go

EXPOSE 2000
ENTRYPOINT ./server