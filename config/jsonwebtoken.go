package config

import "allopopot-interconnect-service/utility"

const JWT_EXPIRY_MINUTES = 60
const JWT_ISSUER = "allopopot-identity-manager"

var JWT_SECRET = utility.ParseEnv("JWT_SECRET", true, utility.GenerateSecret(18))
