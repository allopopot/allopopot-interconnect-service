FROM docker.io/fedora

WORKDIR /app

COPY ./allopopot-interconnect-service /app/allopopot-interconnect-service

CMD [ "/app/allopopot-interconnect-service" ]