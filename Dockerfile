FROM golang:latest

ENV REDIS_ADDR=myredis.5ldsry.0001.use1.cache.amazonaws.com:6379
ENV PORT=8080
 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 

RUN go build

EXPOSE $PORT

CMD ["/app/rediclient"]
