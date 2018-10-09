FROM golang:1.11
WORKDIR /go/src/github.com/tehcyx/cloud-build-poc
RUN go get github.com/jinzhu/gorm
RUN go get gopkg.in/DATA-DOG/go-sqlmock.v1
RUN go get github.com/gorilla/mux
RUN go get github.com/urfave/negroni
RUN go get github.com/lib/pq
COPY . .
RUN go test ./...
RUN go test ./... -cover
RUN go build -i -o bin/app

FROM scratch
# COPY bin/appdocker /app
COPY --from=0 /go/src/github.com/tehcyx/cloud-build-poc/bin/app /app
ENTRYPOINT [ "/app" ]