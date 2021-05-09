package command

import (
	"context"
	"flag"
	"fmt"

	"github.com/seoyhaein/caleb-middle/conf"
	"github.com/seoyhaein/caleb-middle/mesos"
)

// 여기가 실제적인 mesos-go 의 App 역활을 하는 곳이다.

type ServerCommand struct {
	*BaseCommand
	Version  string
	confPath string
	//testLogging bool
}

// cli 사용하기 위해서 3개의 인터페이스 함수를 구현해줘야 함.
func (c *ServerCommand) Run(args []string) int {
	fs := c.Flags()
	fs.Parse(args) // Run 의 args 가 들어오네
	// mesos 환경설정 파일 읽어와서 세팅해줌.
	conf, err := conf.New(c.confPath) // 수정 검토. readability 관점

	// error 처리 해야함.
	err = mesos.Run(context.Background(), conf)

	if err != nil {
		fmt.Println("왜")
		return 1 // 이 부분 살펴보자
	}

	return 0
}

func (c *ServerCommand) Help() string {
	return "nil"
}

func (c *ServerCommand) Synopsis() string {
	return "Start a caleb-middle server!"
}

// flagset 만들어주기
// https://blog.rapid7.com/2016/08/04/build-a-simple-cli-tool-with-golang/ 읽어보기

func (c *ServerCommand) Flags() *flagSet {
	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.Usage = func() { c.Printf(c.Help()) }
	fs.StringVar(&c.confPath, "config", "config.json", "Path to configuration file")
	//fs.BoolVar(&c.testLogging, "testlogging", false, "log sample error and exit")
	return &flagSet{fs}
}
