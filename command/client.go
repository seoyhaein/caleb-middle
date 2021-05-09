package command

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/mesos/mesos-go/api/v1/lib/backoff"
)

// ClientCommand implements interactive client.
type ClientCommand struct {
	*BaseCommand
	addr    string
	auth    string
	Version string
	//jobs                 []*model.Job
	jobsMut              sync.Mutex
	group                string
	project              string
	serverCallsThrottler <-chan struct{}
	token                string
	username             string
	password             string
	//apiClient            *apiclient.Client
}

func (c *ClientCommand) Run(args []string) int {

	fs := c.Flags()
	fs.Parse(args)
	fmt.Printf("Caleb-middle %s\n", c.Version)
	fmt.Println("Please use `exit` or `Ctrl-D` to exit.")

	// login 기능 omit
	// 자세히 파악하기
	c.serverCallsThrottler = backoff.Notifier(time.Second*10, time.Second*10, context.Background().Done())

	return 0
}

func (c *ClientCommand) Help() string {
	help := `
Usage: caleb-middle client [options]

  Start interactive client.

` + c.Flags().help()
	return strings.TrimSpace(help)
}

func (c *ClientCommand) Synopsis() string {

	return "Start interactive client"
}

func (c *ClientCommand) Flags() *flagSet {
	fs := flag.NewFlagSet("client", flag.ContinueOnError)
	fs.Usage = func() { c.Printf(c.Help()) }
	fs.StringVar(&c.addr, "addr", "", "Address of Caleb-middle server (with protocol e.g. \"https://example.com\")")
	return &flagSet{fs}
}
