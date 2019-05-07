FROM golang
WORKDIR /go/src/github.com/go-park-mail-ru/2019_1_DeathPacito/
ADD . ./

RUN go get github.com/google/uuid
RUN go get github.com/dgrijalva/jwt-go
RUN go get google.golang.org/grpc
RUN go build server.go 
EXPOSE 8081
ENTRYPOINT ./server