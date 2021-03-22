package main

import (
	surveyCore "github.com/AlecAivazis/survey/v2/core"
	"github.com/darmiel/dlive-dl/internal/cmd"
	"github.com/mgutz/ansi"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	surveyCore.TemplateFuncsWithColor["color"] = func(style string) string {
		switch style {
		case "white":
			return ansi.ColorCode("default")
		default:
			return ansi.ColorCode(style)
		}
	}

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
						Name:    "url",
						Aliases: []string{"u"},
					},
					&cli.StringFlag{
						Name:    "out-file",
						Aliases: []string{"f"},
					},
					&cli.StringFlag{
						Name:    "variant",
						Aliases: []string{"v"},
					},
					&cli.BoolFlag{
						Name:    "no-progress-bar",
						Aliases: []string{"P"},
					},
				},
				Action: cmd.Download,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalln("🤬", err)
	}
}
