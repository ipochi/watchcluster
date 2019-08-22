package main

import (
	//"fmt"

	log "github.com/golang/glog"
	"github.com/ipochi/watchcluster/config"
	"github.com/ipochi/watchcluster/watchcluster"
)

func main() {

	log.Info("Starting watchcluster controller ... ")

	// Load up a config for watchcluster
	watchclusterConfig, err := config.New()
	if err != nil {
		log.Fatalf("Error loading config ...  %s", err.Error())
	}

	watchcluster.Start(watchclusterConfig)
	//Start(watchclusterConfig)

	// 	log.Logger.Info("Starting controller")
	// Config, err := config.New()
	// if err != nil {
	// 	log.Logger.Fatal(fmt.Sprintf("Error in loading configuration. Error:%s", err.Error()))
	// }

	// if Config.Communications.Slack.Enabled {
	// 	log.Logger.Info("Starting slack bot")
	// 	sb := bot.NewSlackBot()
	// 	go sb.Start()
	// }

	// if Config.Communications.Mattermost.Enabled {
	// 	log.Logger.Info("Starting mattermost bot")
	// 	mb := bot.NewMattermostBot()
	// 	go mb.Start()
	// }

	// if Config.Settings.UpgradeNotifier {
	// 	log.Logger.Info("Starting upgrade notifier")
	// 	go controller.UpgradeNotifier(Config)

	// }

	// 		// Init KubeClient, InformerMap and start controller
	// 		utils.InitKubeClient()
	// 	utils.InitInformerMap()
	// 	controller.RegisterInformers(Config)
}
