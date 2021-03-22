package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:    "dlive-dl",
		Usage:   "DLive Downloader",
		Version: "1.0.0",
		Commands: []*cli.Command{
			{
				Name:    "download",
				Aliases: []string{"dl"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "url",
						Aliases:  []string{"u"},
						Required: true,
					},
				},
				Action: cmdDownload,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalln("🤬", err)
	}
}
