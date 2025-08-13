FROM golang:1.24 AS build

WORKDIR /workspace
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -trimpath -o webhook -ldflags '-w -extldflags "-static"' .

FROM gcr.io/distroless/static

COPY --from=build /workspace/webhook /usr/local/bin/webhook

ENTRYPOINT ["webhook"]
