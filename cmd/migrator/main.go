package main

import (
	"flag"
	"fmt"
	app "github.com/hyperfyodor/yq_sample/internal/app/migrator"
	config "github.com/hyperfyodor/yq_sample/internal/config/migrator"
	"github.com/hyperfyodor/yq_sample/pkg"
	"github.com/ilyakaznacheev/cleanenv"
)

func main() {
	versionFlag := flag.Bool("version", false, "print version information")
	configExplainFlag := flag.Bool("cfg_explain", false, "print explanation of configuration options")
	stepsFlag := flag.Int("steps", 0, "migration steps to apply")
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

	migrator := app.MustLoadMigratorApp()

	if *stepsFlag != 0 {
		if err := migrator.Steps(*stepsFlag); err != nil {
			panic(err)
		}

		return
	}

	if err := migrator.Up(); err != nil {
		panic(err)
	}
}
