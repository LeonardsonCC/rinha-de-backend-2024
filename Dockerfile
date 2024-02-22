FROM golang:1.22-alpine AS build

RUN apk add --no-cache git

WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -pgo=auto -ldflags="-w -s" -o ./app ./cmd/api


# try use scratch
FROM scratch

COPY --from=build /src/app /app

EXPOSE 8888

CMD ["/app"]
