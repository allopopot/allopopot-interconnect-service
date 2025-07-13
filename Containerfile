FROM docker.io/fedora:latest

WORKDIR /app

COPY ./allopopot-interconnect-service ./allopopot-interconnect-service

EXPOSE 4000

CMD [ "./allopopot-interconnect-service" ]