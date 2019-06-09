# build stage
FROM golang AS build-env
ADD . /src
RUN cd /src && GOOS=linux go build

FROM alpine
WORKDIR /app
COPY --from=build-env /src/gps2.0 /app
ENTRYPOINT ./gps2.0