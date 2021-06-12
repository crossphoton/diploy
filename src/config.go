package src

var DB_URL = "./diploy.db"
var DB_PATH = "./"
var LOG_PATH = "./logs/diploy/"

var MODES = map[string]bool{
	"run":    true,
	"build":  true,
	"update": true,
}
