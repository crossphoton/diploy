package src

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Command struct {
	Command string `yaml:"command" gorm:"column:command"`
	Type    string `yaml:"type" gorm:"column:type"`
}

func (c *Command) start(name, logfile, workdir string) (err error) {
	command := strings.Split(c.Command, " ")
	_, err = StartProcess(command[0], command[1:], name, logfile, workdir, (c.Type == "script"))
	return
}

type Config struct {
	Name    string  `yaml:"name" gorm:"primaryKey"`
	Update  Command `yaml:"update" gorm:"embedded;embeddedPrefix:update_"`
	Build   Command `yaml:"build" gorm:"embedded;embeddedPrefix:build_"`
	Run     Command `yaml:"run" gorm:"embedded;embeddedPrefix:run_"`
	Workdir string  `yaml:"-" gorm:"column:workdir"` // Not set by user though
}

func (c *Config) Start(mode string) error {
	switch mode {
	case "build":
		return c.Build.start(c.Name, LOG_PATH+"/"+c.Name+"/build", c.Workdir)
	case "run":
		return c.Run.start(c.Name, LOG_PATH+"/"+c.Name+"/run", c.Workdir)
	case "update":
		return c.Update.start(c.Name, LOG_PATH+"/"+c.Name+"/update", c.Workdir)
	default:
		return fmt.Errorf("%s: not a valid option", mode)
	}
}

func (c *Config) Restart() {
	c.Stop()
	c.Start("run")
}

func (c *Config) Stop() {
	fmt.Println("Stopping: ", c.Name)
	processes, err := getProcesses(c.Name)
	if err != nil {
		fmt.Println("couldn't retrieve any processes for given configuration")
		return
	}

	for _, process := range processes {
		p, err := os.FindProcess(process.Pid)
		if err != nil {
			continue
		}

		if p.Signal(os.Interrupt) != nil {
			fmt.Println("couldn't kill proces with id: ", process.Pid)
		}
	}
	deleteProcesses(c.Name)
}

func (c *Config) Delete() error {
	return deleteConfig(c)
}

// StartProcess starts a process with given
func StartProcess(command string, args []string, name, logFile, workdir string, script bool) (pid int, err error) {
	if script {
		scriptPath := []string{workdir + "/" + command}
		args = append(scriptPath, args...)
		command = "$SHELL"
	}
	process := exec.Command(command, args...)
	process.Dir = workdir
	filename := fmt.Sprintf("%s-%d", logFile, time.Now().Unix())
	process.Stdout, err = os.Create(filename)
	process.Stderr, err = os.Create(filename + "-errors")
	if err != nil {
		fmt.Println("Couldn't create log file: ", err)
	}

	err = process.Start()
	if err != nil {
		fmt.Println("couldn't start the process: ", err)
		return
	}
	pid = process.Process.Pid
	processEntry(name, pid)

	go process.Wait() // IDK if this is a good idea or not

	println("Exiting start")

	return
}

func AddFromFile(file []byte, workdir string) error {
	config, err := Parser(file)
	if err != nil {
		return err
	}
	config.Workdir = workdir

	if config.Update.Command == "" {
		config.Update.Command = "git pull"
		config.Update.Type = "command"
	}

	err = saveConfig(&config)
	if err != nil {
		return err
	}

	// For logs
	os.Mkdir(LOG_PATH+"/"+config.Name, 0700)

	return nil
}

func Parser(file []byte) (Config, error) {
	var b Config
	err := yaml.Unmarshal(file, &b)
	if err != nil {
		return Config{}, err
	}
	return b, nil
}

func SearchConfig(configName string) (Config, error) {
	c := Config{Name: configName}
	err := getConfig(&c)
	return c, err
}
