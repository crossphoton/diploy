package src

var LOG_PATH string

var MODES = map[string]bool{
	"run":    true,
	"build":  true,
	"update": true,
}

type CommandType string

const (
	command CommandType = "command"
	script  CommandType = "script"
)
