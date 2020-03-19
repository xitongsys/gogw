package main

import ( 
	"flag"
	"sync"

	"gogw/logger"
	"gogw/config"
	"gogw/server"
	"gogw/client"
)

var cfgFile = flag.String("c", "cfg.json", "config file")
var role = flag.String("r", "server", "role: server/client")
var logLevel = flag.String("l", "info", "log level: info/debug")

func main(){
	logger.LEVEL = logger.INFO

	logger.Info("gogw start")
	flag.Parse()

	cfg, err := config.NewConfigFromFile(*cfgFile)
	if err != nil {
		logger.Error(err)
		return
	}

	if *logLevel == "debug" {
		logger.LEVEL = logger.DEBUG
	}

	if *role == "server" {
		server := server.NewServer(cfg.Server.ServerAddr, cfg.Server.TimeoutSecond)
		server.Start()
	}

	if *role == "client" {
		var wg sync.WaitGroup

		for _, cfg := range cfg.Clients {
			c := client.NewClient(
				cfg.ServerAddr, 
				cfg.SourceAddr, 
				cfg.ToPort, 
				cfg.Direction, 
				cfg.Protocol, 
				cfg.Description,
				cfg.Compress,
				cfg.HttpVersion,
			)

			wg.Add(1)
			go func(){
				defer wg.Done()
				c.Start()
			}()
		}

		wg.Wait()
	}
}