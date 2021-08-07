FROM golang:1.16-alpine AS build
WORKDIR /go/src
COPY . ./
ENV CGO_ENABLED=0
RUN go build \
    -o /out/service \
    ./cmd/face-search/main.go

FROM keyvanfatehi/chrome-xvfb
WORKDIR /app
COPY --from=build /out/service /app/service
CMD [ "Xvfb", ":99", "-screen", "0", "640x480x8", "-nolisten", "tcp", "&", "google-chrome"]
ENTRYPOINT ["/app/service"]
