FROM golang

COPY . /go/src/datamanager
WORKDIR /go/src/datamanager

RUN go get github.com/streadway/amqp
RUN go get -d -v
RUN CGO_ENABLED=0 go build -o datamanager ./exec/main.go

#--------------------------- 

FROM alpine:latest

# configuration for 
# RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
# RUN update-ca-certificates

WORKDIR /

COPY --from=0 /go/src/datamanager/datamanager /datamanager

#-- Example for initialization files, if needed
# COPY --from=0 /go/src/elevations/web /web
# COPY --from=0 /go/src/elevations/config.ini /config.ini
# COPY --from=0 /go/src/elevations/urls.json  /urls.json

# Environment Variables

ENV PPM_RABBIT_HOST=rabbit1
ENV PPM_DB_CONNECT="postgres://user:password@postgres1/powerPlantMonitor?sslmode=disable"

CMD ["/datamanager"]

#-- if needed
# EXPOSE 8080