FROM docker.io/golang

WORKDIR /app

COPY ./ /app

EXPOSE 4000

RUN ["go", "build"]

ENTRYPOINT [ "allopopot-interconnect-service" ]