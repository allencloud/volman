package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "volman"
	app.Usage = "A tool for docker container's volume management."
	app.Version = "1.0.0"

	app.Author = "Allen Sun"
	app.Email = "allen.sun@daocloud.io"

	// Commands list of NewApp
	app.Commands = []cli.Command{
		{
			Name:  "all",
			Usage: "display global information of volumes.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config",
					Value: "/etc/stager/config.json",
					Usage: "Config file used to run programmer.",
				},
				cli.DurationFlag{
					Name:  "remove_staleimage_interval",
					Value: 30 * time.Minute,
					Usage: "Interval between two stale images removal",
				},
				cli.Float64Flag{
					Name:  "remove_staleimage_percent",
					Value: 0.85,
					Usage: "Percentage which means when to remove stale images.",
				},
				cli.StringFlag{
					Name:  "private_registry",
					Value: "10.10.15.97:5050",
					Usage: "Private Registry Address to push images",
				},
				cli.StringFlag{
					Name:  "dockerhub_private_registry",
					Value: "10.10.89.210:5000",
					Usage: "Private Registry Address to push images",
				},
				cli.BoolFlag{
					Name:  "registry_verify",
					Usage: " Whether to use an authorized docker registry.",
				},
				cli.StringFlag{
					Name:  "docker_endpoint",
					Value: "unix:///var/run/docker.sock",
					Usage: "Docker endpoint to construct a docker client",
				},
			},
			Action: GetAll,
		},
	}

	app.Run(os.Args)
}
