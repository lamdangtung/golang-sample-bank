FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz


FROM golang:1.21-alpine3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY db/migration ./migration
COPY app.env .
COPY start.sh .
RUN chmod +x ./start.sh
COPY wait-for.sh .
RUN chmod +x ./wait-for.sh
EXPOSE 8080
CMD [ "./main" ]
ENTRYPOINT [ "./start.sh"]