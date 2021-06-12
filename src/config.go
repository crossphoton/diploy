package src

import "os"

var LOG_PATH = os.Getenv("DIPLOY_LOG_PATH")
var DB_URL = LOG_PATH + "/diploy.db"

var MODES = map[string]bool{
	"run":    true,
	"build":  true,
	"update": true,
}
