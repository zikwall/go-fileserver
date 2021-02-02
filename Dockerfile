FROM golang:alpine
RUN apk add ffmpeg
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]