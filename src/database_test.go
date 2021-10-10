package src

import (
	"os"
	"testing"
)

var conf = Config{Name: "test", Update: Command{Command: "node", Type: command}}

func TestConfigOperations(t *testing.T) {
	LOG_PATH = os.TempDir()
	openDatabase()

	err = saveConfig(&conf)
	if err != nil {
		t.Fatal("saveConfig - failed:", err)
	}

	var confGet = Config{Name: "test"}
	err = getConfig(&confGet)
	if err != nil && confGet != conf {
		t.Fatal("saveConfig - failed:", err)
	}
	confGet.Name = "random"

	err = deleteConfig(&confGet)
	if err != nil {
		t.Fatal("deleteConfig - failed:", err)
	}
}

// func TestProcessEntry(t *testing.T) {
// 	err = processEntry(conf.Name, 46)
// 	if err != nil {
// 		t.Fatal("processEntry - failed:", err)
// 	}
// 	processes, err := getProcesses(conf.Name)
// 	if err != nil || len(processes) != 1 {
// 		t.Fatal("getProcesses - failed")
// 	}

// 	err = deleteProcesses(conf.Name)
// 	if err != nil {
// 		t.Fatal("deleteProcesses - failed")
// 	}

// }
