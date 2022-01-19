package constant

import "os"

var JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")