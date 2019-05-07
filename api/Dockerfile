FROM golang
WORKDIR /opt/api


RUN go get github.com/google/uuid
RUN go get github.com/gorilla/mux
RUN go get github.com/mailru/easyjson
RUN go get github.com/mailru/easyjson/jlexer
RUN go get github.com/mailru/easyjson/jwriter
RUN go get go.mongodb.org/mongo-driver/mongo
RUN go get go.mongodb.org/mongo-driver/bson
RUN go get github.com/dgrijalva/jwt-go
RUN go get google.golang.org/grpc

COPY ./api ./api
COPY ./identity_struct ./identity_struct
COPY ./auth ./auth

RUN go build ./api/server.go 
EXPOSE 8080
ENTRYPOINT ./server
