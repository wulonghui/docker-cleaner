package main

import (
	"github.com/codegangsta/cli"
)

var Commands = []cli.Command{
	commandVolume,
	commandImage,
	commandContainer,
}

var commandVolume = cli.Command{
	Name:      "volume",
	ShortName: "v",
	Usage:     "Removes orphaned volumes from Docker host",
	Action:    doVolume,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "dryrun, D",
			Usage: "dry run, do not delete anything",
		},
		cli.BoolFlag{
			Name:  "force, f",
			Usage: "delete volumes by force",
		},
	},
}

var commandImage = cli.Command{
	Name:      "image",
	ShortName: "i",
	Usage:     "Removes orphaned images from Docker host",
	Action:    doImage,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "dryrun, D",
			Usage: "dry run, do not delete anything",
		},
		cli.BoolFlag{
			Name:  "force, f",
			Usage: "delete images by force",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "specify a image name",
		},
		cli.StringFlag{
			Name:  "duration, d",
			Value: "0s",
			Usage: "delete images whose Created is passed for a specified duration, e.g. 10s,10m,1h10m",
		},
	},
}

var commandContainer = cli.Command{
	Name:      "container",
	ShortName: "c",
	Usage:     "Removes orphaned containers from Docker host",
	Action:    doContainer,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "dryrun, D",
			Usage: "dry run, do not delete anything",
		},
		cli.BoolFlag{
			Name:  "force, f",
			Usage: "delete containers by force",
		},
		cli.StringFlag{
			Name:  "duration, d",
			Value: "0s",
			Usage: "delete containers whose Fininished is passed for a specified duration, e.g. 10s,10m,1h10m",
		},
	},
}
