FROM golang:alpine AS build 
WORKDIR /build
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o tftp

FROM scratch
COPY --from=build /build/tftp tftp
ENTRYPOINT ["/tftp"]
