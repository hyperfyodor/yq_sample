package main

import (
	"context"
	"flag"
	"fmt"
	consumer "github.com/hyperfyodor/yq_sample/internal/app/consumer"
	config "github.com/hyperfyodor/yq_sample/internal/config/consumer"
	"github.com/hyperfyodor/yq_sample/pkg"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	versionFlag := flag.Bool("version", false, "print version information")
	configExplainFlag := flag.Bool("cfg_explain", false, "print explanation of configuration options")
	flag.Parse()

	if *versionFlag {
		fmt.Println(pkg.Version)
		return
	}

	if *configExplainFlag {
		var cfg config.Config

		help, err := cleanenv.GetDescription(&cfg, nil)

		if err != nil {
			fmt.Println("failed to get help")
			return
		}

		fmt.Println(help)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())

	app := consumer.MustLoad(ctx)

	go app.Start()
	go app.StartMetrics()
	go app.StartProfiling()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	app.Stop()
	cancel()
	log.Println("Consumer finished!")

}
