package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"strconv"
)

// a Config struct to pass into our app
type Config struct {
	DB      db
	Server  server
	Storage storage
	S3      s3
}

type db struct {
	Dialect  string
	Name     string
	Charset  string
	Username string
	Password string
	Host     string
}

type server struct {
	Name       string
	Port       int
	Debug      bool
	AdminToken string
	Cors       []string
}

type storage struct {
	Type string
}

type s3 struct {
	Region string
	Bucket string
	Host   string
}

func (srv *server) ToString() string {
	str := srv.Name + ":" + strconv.Itoa(srv.Port)
	return str
}

func (dbcfg *db) ToString() string {
	str := dbcfg.Host + " db: " + dbcfg.Name + " - user: " + dbcfg.Username
	return str
}

func GetConfig() *Config {
	var config Config

	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		log.Fatal(err)
	}
	return &config
}
