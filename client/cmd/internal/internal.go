package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/poonman/entry-task/client/app"
	"github.com/poonman/entry-task/client/domain"
	"github.com/poonman/entry-task/client/infra/config"
	"github.com/poonman/entry-task/client/infra/gateway"
	"github.com/poonman/entry-task/dora/misc/helper"
	"github.com/spf13/pflag"
	"go.uber.org/dig"
	"os"
	"strings"
)

var (
	Arg Argument
)

type Argument struct {
	Commands    []string `json:"commands"`
	Method      string   `json:"method"`
	Username    string   `json:"user"`
	Password    string   `json:"password"`
	Concurrency int      `json:"concurrency"`
	Requests    int      `json:"requests"`
	Key         string   `json:"key"`
	Value       string   `json:"value"`
}

func (a *Argument) String() string {
	b, err := json.Marshal(a)
	if err != nil {
		return ""
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return ""
	}
	return out.String()
}

func CommandLine() *Argument {
	bmrFlag := pflag.NewFlagSet("etclient benchmark read", pflag.ExitOnError)
	bmrFlag.IntVarP(&Arg.Concurrency, "concurrency", "c", 1, "concurrency")
	bmrFlag.IntVarP(&Arg.Requests, "requests", "r", 1, "requests")
	bmrFlag.StringVarP(&Arg.Username, "user", "u", "1", "username")
	bmrFlag.StringVarP(&Arg.Password, "password", "p", "1", "password")
	bmrFlag.StringVarP(&Arg.Key, "key", "k", "a", "key must be in the range [a:z]")

	bmwFlag := pflag.NewFlagSet("etclient benchmark write", pflag.ExitOnError)
	bmwFlag.IntVarP(&Arg.Concurrency, "concurrency", "c", 1, "concurrency")
	bmwFlag.IntVarP(&Arg.Requests, "requests", "r", 1, "requests")
	bmwFlag.StringVarP(&Arg.Username, "user", "u", "1", "username")
	bmwFlag.StringVarP(&Arg.Password, "password", "p", "1", "password")
	bmwFlag.StringVarP(&Arg.Key, "key", "k", "a", "key must be in the range [a:z]")
	bmwFlag.StringVarP(&Arg.Value, "value", "v", "a", "value must be in the range [a:z]")

	loginFlag := pflag.NewFlagSet("etclient login", pflag.ExitOnError)
	loginFlag.StringVarP(&Arg.Username, "user", "u", "1", "username")
	loginFlag.StringVarP(&Arg.Password, "password", "p", "1", "password")

	readFlag := pflag.NewFlagSet("etclient read", pflag.ExitOnError)
	readFlag.StringVarP(&Arg.Username, "user", "u", "1", "username")
	readFlag.StringVarP(&Arg.Password, "password", "p", "1", "password")
	readFlag.StringVarP(&Arg.Key, "key", "k", "a", "key must be in the range [a:z]")

	writeFlag := pflag.NewFlagSet("etclient write", pflag.ExitOnError)
	writeFlag.StringVarP(&Arg.Username, "user", "u", "1", "username")
	writeFlag.StringVarP(&Arg.Password, "password", "p", "1", "password")
	writeFlag.StringVarP(&Arg.Key, "key", "k", "a", "key must be in the range [a:z]")
	writeFlag.StringVarP(&Arg.Value, "value", "v", "a", "value must be in the range [a:z]")

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	ComposedCommands := os.Args[1]

	Arg.Commands = strings.Split(ComposedCommands, ",")
	cmdMap := make(map[string]struct{})
	for _, cmd := range Arg.Commands {
		cmdMap[cmd] = struct{}{}
	}

	for _, cmd := range Arg.Commands {
		switch cmd {
		case "benchmark":

			Arg.Method = os.Args[2]
			if Arg.Method != "read" && Arg.Method != "write" {
				benchmarkUsage()
			}

			if Arg.Method == "read" {
				err := bmrFlag.Parse(os.Args[3:])
				if err != nil {
					bmrFlag.Usage()
					os.Exit(1)
				}

				if len(Arg.Key) != 1 || (Arg.Key < "a" || Arg.Key > "z") {
					bmrFlag.Usage()
					os.Exit(1)
				}
			} else if Arg.Method == "write" {
				err := bmwFlag.Parse(os.Args[3:])
				if err != nil {
					bmwFlag.Usage()
					os.Exit(1)
				}

				if len(Arg.Key) != 1 || (Arg.Key < "a" || Arg.Key > "z") {
					bmwFlag.Usage()
					os.Exit(1)
				}

				if len(Arg.Value) != 1 || (Arg.Value < "a" || Arg.Value > "z") {
					bmwFlag.Usage()
					os.Exit(1)
				}
			} else {
				benchmarkUsage()
			}

		case "login":
			if _, ok := cmdMap["read"]; ok {
				continue
			}

			if _, ok := cmdMap["write"]; ok {
				continue
			}

			if _, ok := cmdMap["benchmark"]; ok {
				continue
			}
			err := loginFlag.Parse(os.Args[2:])
			if err != nil {
				loginFlag.Usage()
				os.Exit(1)
			}

		case "read":

			if _, ok := cmdMap["write"]; ok {
				continue
			}

			if _, ok := cmdMap["benchmark"]; ok {
				continue
			}

			err := readFlag.Parse(os.Args[2:])
			if err != nil {
				readFlag.Usage()
				os.Exit(1)
			}

			if len(Arg.Key) != 1 || (Arg.Key < "a" || Arg.Key > "z") {
				readFlag.Usage()
				os.Exit(1)
			}

		case "write":

			if _, ok := cmdMap["benchmark"]; ok {
				continue
			}
			err := writeFlag.Parse(os.Args[2:])
			if err != nil {
				writeFlag.Usage()
				os.Exit(1)
			}

			if len(Arg.Key) != 1 || (Arg.Key < "a" || Arg.Key > "z") {
				writeFlag.Usage()
				os.Exit(1)
			}

			if len(Arg.Value) != 1 || (Arg.Value < "a" || Arg.Value > "z") {
				writeFlag.Usage()
				os.Exit(1)
			}

		default:
			usage()
		}
	}

	return &Arg
}

func usage() {
	fmt.Println("Usage: etclient COMMAND [OPTIONS]")
	fmt.Println()
	fmt.Println("Support composed command separated by ','. eg: './etclient login,read -u 1 -p 1 -k a' ")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("\tbenchmark")
	fmt.Println("\tlogin")
	fmt.Println("\tread")
	fmt.Println("\twrite")
	fmt.Println()
	fmt.Println("Run 'etclient COMMAND --help' for more information on a Commands.")
	os.Exit(0)
}

func benchmarkUsage() {
	fmt.Println("Usage: etclient benchmark METHOD [OPTIONS]")
	fmt.Println()
	fmt.Println("Methods:")
	fmt.Println("\tread")
	fmt.Println("\twrite")
	fmt.Println()
	fmt.Println("Run 'etclient benchmark METHOD --help' for more information on a Commands.")
	os.Exit(0)
}

func BuildContainer() *dig.Container {
	c := dig.New()

	helper.MustContainerProvide(c, config.NewConfig)
	helper.MustContainerProvide(c, gateway.NewKvGateway)
	helper.MustContainerProvide(c, app.NewService)
	helper.MustContainerProvide(c, domain.NewService)

	return c
}
