package config

import "allopopot-interconnect-service/utility"

var MONGODB_URI string = utility.ParseEnv("MONGODB_URI", false, "")
var MONGODB_DATABASE_NAME string = utility.ParseEnv("MONGODB_DATABASE_NAME", true, "allopopot-services")
