# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

RUN  mkdir -p /go/src \
  && mkdir -p /go/bin \
  && mkdir -p /go/pkg
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH   

RUN mkdir -p $GOPATH/src/PowerPlantMonitor \
 && mkdir -p $GOPATH/src/PowerPlantMonitor/distributed \
 && mkdir -p $GOPATH/src/PowerPlantMonitor/distributed/coordinator 

# Copy the local package files to the container's workspace.
ADD . $GOPATH/src/PowerPlantMonitor/distributed/coordinator

WORKDIR $GOPATH/src/PowerPlantMonitor/distributed/coordinator 

# Environment Variables
# ENV AUTOSTART_INGESTION_JOB false
# ENV INGESTOR_DOWNLOAD_FOLDER /tmp
# ENV INGESTION_TASK_FREQUENCY_M 10
# ENV INGESTION_TASK_FREQUENCY_H 4
# ENV HANA_USER <user>
# ENV HANA_USER_PW <password>
# ENV HANA_INSTANCE <instance>
# ENV HANA_DB_SCHEMA <schema>
# ENV HANA_EVENT_TABLE <Events Table>
# ENV HANA_USER_TABLE <Contributers Table>
# ENV HANA_STATUS_TABLE <Status Table>

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get -d -v
RUN go build -o coordinator ./exec/main.go

# Run the outyet command by default when the container starts.
# CMD ["./coordinator"]

# Document that the service listens on port 8080.
EXPOSE 8080