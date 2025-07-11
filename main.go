package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"
)

var (
	AppVersion = "development"
	hostFlag   = cli.StringSliceFlag{
		Name:       "host",
		Required:   true,
		HasBeenSet: true,
		Usage:      "Target Host",
	}
	portFlag = cli.StringFlag{
		Name:        "port",
		Required:    false,
		HasBeenSet:  true,
		Usage:       "Ports Range e.g 80 or 1-1024 or 80,22,23",
		DefaultText: `1-32000`,
		Value:       `1-32000`,
	}
	timeoutFlag = cli.StringFlag{
		Name:       "timeout",
		Required:   false,
		HasBeenSet: true,
		Usage:      `TCP Timeout in Millisecond`,
		Value:      `500`,
	}
	parallelismFlag = cli.IntFlag{
		Name:        "parallelism",
		Required:    false,
		HasBeenSet:  true,
		Usage:       `How many parallel job to run`,
		DefaultText: fmt.Sprintf(`%d`, runtime.NumCPU()),
		Value:       runtime.NumCPU(),
	}
	debugFlag = cli.BoolFlag{
		Name:     "debug",
		Required: false,
		Usage:    `Enable Debug Logs`,
		Value:    false,
	}
)

func main() {
	cliApp := &cli.App{
		Name:                 "port-scanner-go",
		Version:              AppVersion,
		EnableBashCompletion: true,
		Authors: []*cli.Author{
			{
				Name:  "M Amin Nasiri",
				Email: "Khodexenon@gmail.com",
			},
			{
				Name:  "Dmytro Horkhover",
				Email: "gd.mail.89@gmail.com",
			},
		},
		Flags: []cli.Flag{
			&hostFlag,
			&portFlag,
			&timeoutFlag,
			&parallelismFlag,
			&debugFlag,
		},
		Action: cliAction,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer cancel()

	if err := cliApp.RunContext(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}

func cliAction(c *cli.Context) error {

	debug := c.Bool(debugFlag.Name)
	tcpTimeout := time.Duration(c.Int(timeoutFlag.Name)) * time.Millisecond
	concurrency := c.Int(parallelismFlag.Name)
	hosts := c.StringSlice(hostFlag.Name)

	var allIPs []string
	for _, host := range hosts {
		ips, err := nslookup(host)
		if err != nil {
			return err
		}
		allIPs = append(allIPs, ips...)
	}

	if len(allIPs) == 0 {
		return fmt.Errorf(`nslookup returns empty list of IPs for hosts "%s"`, strings.Join(hosts, ", "))
	}

	portsList, err := getPortsList(c.String(portFlag.Name))
	if err != nil {
		return err
	}

	jobsCount := len(allIPs) * len(portsList)
	executor := NewJobsExecutor(c.Context, jobsCount, concurrency)
	defer executor.Shutdown()

	for _, ip := range allIPs {
		for _, port := range portsList {
			s := scanner{
				ip:      ip,
				port:    port,
				timeout: tcpTimeout,
				debug:   debug,
			}
			executor.Submit(s.scan)
		}
	}

	return nil
}
