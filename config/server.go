package config

import "fmt"

const SERVER_HOST string = "0.0.0.0"
const SERVER_PORT string = "5000"

var SERVER_URI = fmt.Sprintf("%s:%s", SERVER_HOST, SERVER_PORT)
