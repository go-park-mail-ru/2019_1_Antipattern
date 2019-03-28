FROM golang
WORKDIR /go/src/github.com/go-park-mail-ru/2019_1_DeathPacito/
ADD . ./
RUN CGO_ENABLED=0 GOOS=linux \
	go get github.com/google/uuid &&  \
	go get github.com/gorilla/mux && \
	go get github.com/mailru/easyjson && \
	go get github.com/mailru/easyjson/jlexer && \
	go get github.com/mailru/easyjson/jwriter && \
	go build server.go 
EXPOSE 8080
ENTRYPOINT ./server