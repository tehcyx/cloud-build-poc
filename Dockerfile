FROM scratch
COPY bin/appdocker /app
ENTRYPOINT [ "/app" ]