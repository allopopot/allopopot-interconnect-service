package config

import "fmt"

const DATABASE_HOST string = "10.0.0.2"
const DATABASE_USER string = "admin"
const DATABASE_PASSWORD string = "qwertyuiop"
const DATABASE_NAME string = "postgres"
const DATABASE_PORT string = "5432"
const DATABASE_SSL_MODE string = "disable"

var DSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", DATABASE_HOST, DATABASE_USER, DATABASE_PASSWORD, DATABASE_NAME, DATABASE_PORT, DATABASE_SSL_MODE)
