FROM golang:alpine AS build 
ARG ARCH_VALUE
ENV DOCKER_ARCH=$ARCH_VALUE
WORKDIR /build
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${DOCKER_ARCH} go build -o tftp

FROM scratch
COPY --from=build /build/tftp tftp
ENTRYPOINT ["/tftp"]
