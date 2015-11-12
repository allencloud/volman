package main

import (
	//"fmt"
	"github.com/codegangsta/cli"
	"os"
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
			Name:   "all",
			Usage:  "display global information of volumes.",
			Action: GetAll,
		},
		{
			Name:   "con",
			Usage:  "display specified information of volume.",
			Action: GetOne,
		},
	}

	app.Run(os.Args)
}
