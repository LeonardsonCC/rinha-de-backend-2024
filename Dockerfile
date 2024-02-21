FROM golang:1.22-alpine AS build

RUN apk add --no-cache git

WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./app ./cmd/api


# try use scratch
FROM alpine:3.19

COPY --from=build /src/app /app

EXPOSE 8888

CMD ["/app"]
