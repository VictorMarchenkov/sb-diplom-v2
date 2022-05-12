FROM golang:1.16 AS api

WORKDIR /buildapp
COPY . .
ADD go.mod go.sum /buildapp/

RUN CGO_ENABLED=0 GOOS=linux go build -o ./generator/generator ./generator/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bind/app/diploma ./cmd/app/main.go

EXPOSE 8282
EXPOSE 8383

CMD ["./start.sh"]

