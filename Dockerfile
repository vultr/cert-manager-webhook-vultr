FROM ewr.vultrcr.com/dhi.io/golang:1.27-alpine AS build_deps

WORKDIR /workspace

COPY go.mod .
COPY go.sum .

RUN go mod download

FROM build_deps AS build

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -o /tmp/webhook -ldflags '-w -extldflags "-static"' .

FROM ewr.vultrcr.com/dhi.io/static:20230311

COPY --from=build /tmp/webhook /usr/local/bin/webhook

ENTRYPOINT ["webhook"]
