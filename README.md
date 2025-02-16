# AlloPopoT Interconnect Service

A REST API server to for inter-service commmunication.
Written in Golang.

## Build the Project

To build the project run:

    go build

## Environment Variables

| Name                  | Required | Default            |
| --------------------- | -------- | ------------------ |
| SERVER_PORT           | false    | 4000               |
| MONGODB_URI           | true     | -                  |
| MONGODB_DATABASE_NAME | false    | allopopot-services |
| AMQP_HOST             | false    | -                  |
| AMQP_PORT             | false    | 5672               |
| AMQP_USERNAME         | false    | -                  |
| AMQP_PASSWORD         | false    | -                  |
| JWT_SECRET            | false    | <AUTO_GENERATED>   |

> Please note: Email AMQP service will be disabled if wrong or no credentials are passed  in environment variable.