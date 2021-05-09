package main

import (
	"github.com/onrik/logrus/filename"
	"github.com/seoyhaein/caleb-middle/command"
	"os"

	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
)

// 이거 왜 했는지 파악해야함.
func init() {
	log.AddHook(filename.NewHook())
}

const version = "0.01"

func main() {
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}
	coloredui := &cli.ColoredUi{
		Ui:         ui,
		ErrorColor: cli.UiColorRed,
	}
	c := cli.NewCLI("caleb-middle", version)
	c.Args = os.Args[1:]
	baseCmd := command.BaseCommand{Ui: coloredui}
	c.Commands = map[string]cli.CommandFactory{
		// 여기에 함수들을 넣어줘야 함. 형식은 string func()(cli.Command,error) 임
		// cli.Command 는 인터페이스 임. 포인터로 보냄.
		"server": func() (cli.Command, error) {
			return &command.ServerCommand{BaseCommand: &baseCmd, Version: version}, nil
		},

		"client": func() (cli.Command, error) {
			return &command.ClientCommand{BaseCommand: &baseCmd, Version: version}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Error(err)
	}
	os.Exit(exitStatus)
}
