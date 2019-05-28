# build stage
FROM golang
WORKDIR /src
ADD . .

RUN go build -o /bin/gpsv2
WORKDIR /
RUN rm -r /src

CMD ["/bin/gpsv2"]