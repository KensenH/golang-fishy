FROM golang:1.18 AS build

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

WORKDIR /work
COPY . /work

RUN go build -o bin/fish .


FROM alpine:3.18.3 AS run

COPY --from=build /work/bin/fish /usr/local/bin/
COPY fish /home/fish

RUN mkdir /home/public-keys

RUN apk --no-cache add tzdata

WORKDIR /home

CMD ["fish"]