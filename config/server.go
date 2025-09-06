package config

import (
	"allopopot-interconnect-service/utility"
	"fmt"
)

var SERVER_HOST string = "0.0.0.0"
var SERVER_PORT string = utility.ParseEnv("SERVER_PORT", true, "4000")

var SERVER_URI = fmt.Sprintf("%s:%s", SERVER_HOST, SERVER_PORT)

var CORS_ORIGINS = []string{"http://localhost:5173", "http://localhost:4173"}
