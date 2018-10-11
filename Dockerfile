FROM scratch
COPY bin/appdocker /app
# COPY --from=0 /go/src/github.com/tehcyx/cloud-build-poc/bin/app /app
ENTRYPOINT [ "/app" ]