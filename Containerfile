FROM docker.io/golang

WORKDIR /app

COPY ./ ./

EXPOSE 4000

RUN go build -o allopopot-interconnect-service

CMD [ "./allopopot-interconnect-service" ]