FROM golang:1.22

ENV GOPATH=/
COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh


RUN go mod download
RUN go build -o lts-migrate ./cmd/migrate/main.go
RUN go build -o lts ./cmd/lts/main.go
CMD ["./lts-migrate"]
CMD ["./lts"]