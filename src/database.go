package src

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var database *gorm.DB
var err error

type process struct {
	ID   uint
	Name string
	Pid  int
}

func openDatabase() {
	database, err = gorm.Open(sqlite.Open(DB_URL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "couldn't open database file")
		panic(err)
	}
	database.AutoMigrate(&Config{})
	database.AutoMigrate(&process{})
}

func getConfig(c *Config) error {
	if database == nil {
		openDatabase()
	}

	return database.First(c).Error
}

func saveConfig(c *Config) error {
	if database == nil {
		openDatabase()
	}

	return database.Create(c).Error
}

func processEntry(name string, pid int) error {
	a := process{Name: name, Pid: pid}

	return database.Save(&a).Error
}

func getProcesses(name string) (a []process, err error) {
	err = database.Where("name = ?", name).Find(&a).Error
	if err != nil {
		return nil, err
	}
	return
}

func deleteProcesses(name string) error {
	return database.Where("name = ?", name).Delete(process{}).Error
}

func deleteConfig(c *Config) error {
	return database.Delete(c).Error
}
