package main

import (
	//"fmt"

	log "github.com/golang/glog"
	"github.com/ipochi/watchcluster/config"
	"github.com/ipochi/watchcluster/watchcluster"
)

func main() {

	log.Info("Starting watchcluster controller ... ")

	//	Load up a config for watchcluster
	watchclusterConfig, err := config.New()
	if err != nil {
		log.Fatalf("Error loading config ...  %s", err.Error())
	}
	watchcluster.Start(watchclusterConfig)
}
