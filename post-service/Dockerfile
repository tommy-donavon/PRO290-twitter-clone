FROM golang:1.15-alpine as dev

WORKDIR /work

FROM golang:1.15-alpine as build
WORKDIR /app
COPY . /app/
RUN go mod download
RUN go build -o app

FROM alpine as runtime
COPY --from=build /app/app /
CMD ./app