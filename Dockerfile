FROM golang:1.11
WORKDIR /go/src/github.com/tehcyx/cloud-build-poc
RUN make setup
COPY . /app
RUN make

FROM scratch
# COPY bin/appdocker /app
COPY --from=0 /go/src/github.com/tehcyx/cloud-build-poc/bin/app /app
ENTRYPOINT [ "/app" ]