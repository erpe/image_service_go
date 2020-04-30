package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"strconv"
)

// a Config struct to pass into our app
type Config struct {
	DB         db
	Server     server
	Storage    storage
	S3         s3
	Localstore localstore
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

type localstore struct {
	Directory string
	Assethost string
}

func (srv *server) ToString() string {
	str := srv.Name + ":" + strconv.Itoa(srv.Port)
	return str
}

func (dbcfg *db) ToString() string {
	str := dbcfg.Host + " db: " + dbcfg.Name + " - user: " + dbcfg.Username
	return str
}

func (s *storage) IsS3() bool {
	res := s.Type == "s3"
	return res
}

func (s *storage) IsLocal() bool {
	res := s.Type == "local"
	return res
}
func GetConfig() *Config {
	var config Config

	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		log.Fatal(err)
	}
	return &config
}
