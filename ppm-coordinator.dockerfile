FROM golang

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
RUN mkdir -p $GOPATH/src/github.com/snoopy1964/PowerPlantMonitor 

COPY ./distributed $GOPATH/src/github.com/snoopy1964/powerPlantMonitor/distributed
WORKDIR /go/src/github.com/snoopy1964/powerPlantMonitor

RUN go get github.com/lib/pq
RUN go get github.com/streadway/amqp
# RUN go get -d -v
# RUN CGO_ENABLED=0 go build -o datamanager ./distributed/datamanager/exec/main.go
RUN CGO_ENABLED=0 go build -o coordinator ./distributed/coordinator/exec/main.go

#--------------------------- 

FROM alpine:latest

# configuration for 
# RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
# RUN update-ca-certificates

WORKDIR /

# COPY --from=0 /go/src/github.com/snoopy1964/powerPlantMonitor/datamanager /datamanager
COPY --from=0 /go/src/github.com/snoopy1964/powerPlantMonitor/coordinator /ppm-coordinator

# COPY ./start_app.sh /start_app.sh

#-- Example for initialization files, if needed
# COPY --from=0 /go/src/elevations/web /web
# COPY --from=0 /go/src/elevations/config.ini /config.ini
# COPY --from=0 /go/src/elevations/urls.json  /urls.json

# Environment Variables

ENV PPM_RABBIT_HOST=rabbit1
ENV PPM_DB_CONNECTION=postgres1
ENV PPM_DB_CONNECT="postgres://user:password@postgres1/powerPlantMonitor?sslmode=disable"

CMD ["/ppm-coordinator"]

#-- if needed
# EXPOSE 8080