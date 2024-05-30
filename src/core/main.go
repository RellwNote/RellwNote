package main

import (
	"flag"
	"fmt"
	"github.com/RellwNote/RellwNote/build"
	"github.com/RellwNote/RellwNote/config"
	"github.com/RellwNote/RellwNote/log"
	"github.com/RellwNote/RellwNote/tempServer"
	"os"
)

// 命令行中全部可执行的命令
var actions = map[string]func(){
	"server": func() {
		tempServer.Start()
	},
	"build": func() {
		err := build.Build()
		if err != nil {
			log.Error.Printf("%v\n", err.Error())
		}
	},
	"help": func() {
		if len(os.Args) <= 2 || os.Args[2] == "help" {
			fmt.Printf("Use 'help [command]' to get the help of the command\n")
			return
		}

		if createFunc, has := newFlagSetFunc[os.Args[2]]; has {
			createFunc().Usage()
			printGeneralFlag()
		} else {
			fmt.Printf("Unknow command %s", os.Args[2])
		}
	},
}

// 在这里定义命令的参数
var newFlagSetFunc = map[string]func() *flag.FlagSet{
	"server": func() *flag.FlagSet {
		res := flag.NewFlagSet("server", flag.ExitOnError)
		res.IntVar(&config.ServerPort, "port", config.ServerPort, "server port")
		res.StringVar(&config.ServerHost, "host", config.ServerHost, "server host")
		res.Float64Var(&config.ServerDebugDelay, "debug_delay", config.ServerDebugDelay, "add server delay for debug")
		return res
	},
	"build": func() *flag.FlagSet {
		res := flag.NewFlagSet("server", flag.ExitOnError)
		res.StringVar(&config.BuildOutput, "output", config.BuildOutput, "build result output directory")
		return res
	},
}

// 设置通用命令行参数
func applyGeneralFlag(flag *flag.FlagSet) {
	flag.StringVar(&config.LibraryPath, "library", config.LibraryPath, "library path")
}

func printGeneralFlag() {
	f := flag.NewFlagSet("general", 1)
	applyGeneralFlag(f)
	f.Usage()
}

func printBaseHelp() {
	fmt.Printf("Can use command:\n")
	for k := range actions {
		fmt.Printf("\t%s\n", k)
	}
}

func main() {
	if len(os.Args) == 1 {
		printBaseHelp()
		return
	}

	command := os.Args[1]
	if flagCreateFunc, hasCreateFunc := newFlagSetFunc[command]; hasCreateFunc {
		flagSet := flagCreateFunc()
		applyGeneralFlag(flagSet)
		err := flagSet.Parse(os.Args[2:])
		if err != nil {
			println(err.Error())
			return
		}
	}

	if action, hasAction := actions[command]; hasAction {
		action()
	} else {
		printBaseHelp()
		return
	}
}
