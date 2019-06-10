# build stage
FROM golang AS build-env
COPY . /src
RUN cd /src && CGO_ENABLED=0 GOOS=linux go build -o gpsv2

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/gpsv2 /app/
CMD ["./gpsv2"]