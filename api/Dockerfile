FROM golang
WORKDIR /go/src/github.com/go-park-mail-ru/2019_1_DeathPacito/
ADD . ./

RUN go get github.com/google/uuid
RUN go get github.com/gorilla/mux
RUN go get github.com/mailru/easyjson
RUN go get github.com/mailru/easyjson/jlexer
RUN go get github.com/mailru/easyjson/jwriter
RUN go get go.mongodb.org/mongo-driver/mongo
RUN go get go.mongodb.org/mongo-driver/bson
RUN go build server.go 
EXPOSE 8080
ENTRYPOINT ./server