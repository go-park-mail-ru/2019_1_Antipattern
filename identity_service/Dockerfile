FROM golang
WORKDIR /opt/identity_service
COPY ./identity_service/ ./identity_service/
COPY ./identity_struct ./identity_struct
COPY ./auth ./auth

RUN go get github.com/google/uuid
RUN go get github.com/dgrijalva/jwt-go
RUN go get google.golang.org/grpc
RUN go build ./identity_service/server.go 
EXPOSE 8081
ENTRYPOINT ./server
