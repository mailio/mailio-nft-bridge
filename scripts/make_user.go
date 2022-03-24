package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/chryscloud/go-microkit-plugins/crypto"
	leveldb "github.com/ipfs/go-ds-leveldb"
	"github.com/mailio/mailio-nft-server/model"
	"github.com/mailio/mailio-nft-server/util"
	"gopkg.in/yaml.v3"
)

var Red = "\033[31m"
var Green = "\033[32m"
var Reset = "\033[0m"

type YamlConfig struct {
	DatastorePath string `yaml:"datastore_path"`
}

func main() {

	email := flag.String("email", "", "Admins email address")
	password := flag.String("password", "", "Admins password")
	config := flag.String("config", "", "Config file path")

	if len(os.Args) < 3 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()

	if *email == "" || *password == "" || *config == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// check if config exists
	if _, err := os.Stat(*config); errors.Is(err, os.ErrNotExist) {
		Pkg.Fatal("Config file not found", *config)
	}
	var dat []byte
	if d, err := os.ReadFile(*config); err != nil {
		Pkg.Fatal("Failed to read config file", err.Error())
	} else {
		dat = d
	}
	conf := YamlConfig{}
	uErr := yaml.Unmarshal(dat, &conf)
	if uErr != nil {
		Pkg.Fatal("Failed to parse config file", uErr.Error())
	}

	ds, err := leveldb.NewDatastore(conf.DatastorePath, &leveldb.Options{})
	if err != nil {
		Pkg.Fatal("Failed to open datastore", err.Error())
	}

	passHash, err := crypto.HashPassword(*password)
	if err != nil {
		Pkg.Fatal("Failed to hash password", err.Error())
	}

	u := model.User{
		ID:       "admin",
		Email:    *email,
		Password: passHash,
		Name:     "Admin",
		Created:  time.Now().UnixMilli(),
		Modified: time.Now().UnixMilli(),
	}
	key := util.CreateKey(model.UserTable, u.Email)
	if userBytes, mErr := json.Marshal(u); err != nil {
		Pkg.Fatal("Failed to marshal user: %s", mErr.Error())
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := ds.Put(ctx, key, userBytes); err != nil {
			Pkg.Fatal("Failed to save user: %s", err.Error())
		}
	}

	sayHello()

	fmt.Printf("%sSuccessfully created admin user%s\n", Green, Reset)
}

func sayHello() {
	helloText := `
███╗   ███╗ █████╗ ██╗██╗     ██╗ ██████╗     ███╗   ██╗███████╗████████╗    ██████╗ ██████╗ ██╗██████╗  ██████╗ ███████╗
████╗ ████║██╔══██╗██║██║     ██║██╔═══██╗    ████╗  ██║██╔════╝╚══██╔══╝    ██╔══██╗██╔══██╗██║██╔══██╗██╔════╝ ██╔════╝
██╔████╔██║███████║██║██║     ██║██║   ██║    ██╔██╗ ██║█████╗     ██║       ██████╔╝██████╔╝██║██║  ██║██║  ███╗█████╗  
██║╚██╔╝██║██╔══██║██║██║     ██║██║   ██║    ██║╚██╗██║██╔══╝     ██║       ██╔══██╗██╔══██╗██║██║  ██║██║   ██║██╔══╝  
██║ ╚═╝ ██║██║  ██║██║███████╗██║╚██████╔╝    ██║ ╚████║██║        ██║       ██████╔╝██║  ██║██║██████╔╝╚██████╔╝███████╗
╚═╝     ╚═╝╚═╝  ╚═╝╚═╝╚══════╝╚═╝ ╚═════╝     ╚═╝  ╚═══╝╚═╝        ╚═╝       ╚═════╝ ╚═╝  ╚═╝╚═╝╚═════╝  ╚═════╝ ╚══════╝
`
	fmt.Println(helloText)
}

type Pkg string

func (self Pkg) Fatal(s string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", self, fmt.Sprintf(s, a...))
	os.Exit(2)
}
