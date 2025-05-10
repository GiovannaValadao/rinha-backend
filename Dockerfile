FROM golang:1.24 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rinha-backend

FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/rinha-backend .
EXPOSE 8080
CMD ["./rinha-backend"]
