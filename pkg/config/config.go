package config

import (
	"dbcheck/pkg/runvalue"
	"io/ioutil"
	"path"
	"time"

	"gopkg.in/yaml.v2"
)

type BaseInfo struct {
	DBInfo       dbInfo                 `yaml:"DBInfo"`
	ResultOutput ResultOutputFileEntity `yaml:"ResultOutput"`
	Repair       bool                   `yaml:"Repair"` // 是否尝试修复
}

type dbInfo struct {
	DriverName     string        `yaml:"driverName"`
	Username       string        `yaml:"username"`
	Password       string        `yaml:"password"`
	Host           string        `yaml:"host"`
	Port           string        `yaml:"port"`
	Database       string        `yaml:"database"`
	Charset        string        `yaml:"charset"`
	DBConnIdleTime time.Duration `yaml:"dbConnIdleTime"`
	MaxIdleConn    int           `yaml:"maxIdleConns"`
}

type ResultOutputFileEntity struct {
	OutputWay           string `yaml:"outputWay"`
	OutputPath          string `yaml:"outputPath"`
	OutputFile          string `yaml:"outputFile"`
	InspectionPersonnel string `yaml:"inspectionPersonnel"`
	InspectionLevel     string `yaml:"inspectionLevel"`
}

var C BaseInfo

func init() {
	confPath := path.Join(runvalue.RootPath, "config.yaml")

	yamlFile, err := ioutil.ReadFile(confPath)
	if err != nil {
		panic(err)
	}
	// 将读取到的yaml文件解析为响应的struct
	err = yaml.Unmarshal(yamlFile, &C)
	if err != nil {
		panic(err)
	}
}
