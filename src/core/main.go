package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"rellwnote/core/build"
	"rellwnote/core/config"
	"rellwnote/core/log"
	"rellwnote/core/server"
)

// actions 定义命令行中全部可执行的命令
var actions = map[string]func(){
	"server": func() {
		server.Start()
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

// newFlagSetFunc 中定义了不同命令的参数
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

// preGeneralAction 会在任何子命令执行之前被调用。
// 在这里面处理那些在不同命令之间通用的数据。
func preGeneralAction() {
	if len(config.LibraryName) == 0 {
		config.LibraryName = filepath.Base(config.LibraryPath)
	}
}

// applyGeneralFlag 用于设置通用命令行参数
func applyGeneralFlag(flag *flag.FlagSet) {
	flag.StringVar(&config.LibraryPath, "library", config.LibraryPath, "library path")
	flag.StringVar(&config.LibraryName, "library-name", config.LibraryName, "library name, show in tab title")
	flag.StringVar(&config.Theme, "theme", config.Theme, "theme")
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
		preGeneralAction()
		action()
	} else {
		printBaseHelp()
		return
	}
}
