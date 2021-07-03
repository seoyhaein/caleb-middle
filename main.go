package main

import (
	"context"
	"flag"
	"os"

	"github.com/onrik/logrus/filename"
	"github.com/seoyhaein/caleb-middle/command"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.AddHook(filename.NewHook())
}

// https://github.com/mitchellh/cli 살펴보자
const version = "0.01"

func main() {
	/*ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}
	coloredui := &cli.ColoredUi{
		Ui:         ui,
		ErrorColor: cli.UiColorRed,
	}
	c := cli.NewCLI("caleb-middle", version)
	*/
	// 신규 수정중 7/3
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.Parse(os.Args[1:])
	// 아래 코드 일단 살펴보아야 함. error pron
	if err := command.Run(context.Background(), fs.Args()); err != nil {
		// error prone
		os.Exit(1)
	}

	/*c.Args = os.Args[1:]
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
	}*/

	//exitStatus, err := c.Run()

	/*if err != nil {
		log.Error(err)
	}
	os.Exit(exitStatus)*/
}
