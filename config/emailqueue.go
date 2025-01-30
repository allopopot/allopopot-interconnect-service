package config

import "allopopot-interconnect-service/utility"

var AMQP_HOST = utility.ParseEnv("AMQP_HOST", true, "192.168.137.201")
var AMQP_PORT = utility.ParseEnv("AMQP_PORT", true, "5672")
var AMQP_USERNAME = utility.ParseEnv("AMQP_USERNAME", true, "admin")
var AMQP_PASSWORD = utility.ParseEnv("AMQP_PASSWORD", true, "admin")

var AMQP_EXCHANGE_NAME = utility.ParseEnv("AMQP_EXCHANGE_NAME", true, "EMAIL-SERVICE-EXCHANGE")
