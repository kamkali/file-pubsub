FROM golang:1.19-alpine as build

RUN apk add --no-cache git

WORKDIR /app

COPY internal/app/testdata .
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 go build -o ./app_build ./cmd/file-pubsub

FROM scratch as runner

COPY --from=build ./app/.env ./

COPY --from=build ./app/app_build ./app

ENV CONSUMED_FOLDER_PATH="."
EXPOSE $PORT

ENTRYPOINT ["./app"]