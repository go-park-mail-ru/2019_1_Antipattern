FROM golang
WORKDIR /go/src/github.com/go-park-mail-ru/2019_1_DeathPacito/
ADD . ./

RUN go get github.com/gorilla/websocket
RUN go get github.com/dgrijalva/jwt-go
RUN go build server.go 

EXPOSE 2000
ENTRYPOINT ./server