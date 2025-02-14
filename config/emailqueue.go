package config

import "allopopot-interconnect-service/utility"

var AMQP_HOST = utility.ParseEnv("AMQP_HOST", true, "")
var AMQP_PORT = utility.ParseEnv("AMQP_PORT", true, "5672")
var AMQP_USERNAME = utility.ParseEnv("AMQP_USERNAME", true, "")
var AMQP_PASSWORD = utility.ParseEnv("AMQP_PASSWORD", true, "")

var AMQP_EXCHANGE_NAME = utility.ParseEnv("AMQP_EXCHANGE_NAME", true, "EMAIL-SERVICE-EXCHANGE")
