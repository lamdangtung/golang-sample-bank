FROM golang:1.21-alpine3.18
WORKDIR /app
COPY . .
RUN go build -o main

EXPOSE 8080
CMD [ "./main" ]