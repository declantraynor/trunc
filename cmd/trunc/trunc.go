package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/declantraynor/trunc/web"
	"github.com/urfave/cli/v2"
)

func shorten(c *cli.Context) error {
	server := c.String("server")
	url := c.Args().Get(0)

	client := web.NewClient(server)
	res, err := client.Shorten(url)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(strings.Repeat("-", len(res.ShortURL)))
	fmt.Println(res.ShortURL)
	fmt.Println(strings.Repeat("-", len(res.ShortURL)))
	fmt.Println()

	return nil
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "server",
				Value:   "http://localhost:5000",
				Usage:   "Address of the trunc server to use",
				Aliases: []string{"s"},
				EnvVars: []string{"TRUNC_SERVER"},
			},
		},
		Name:      "trunc",
		Usage:     "Shorten URLs with ease.",
		UsageText: "trunc [global options] URL",
		Action:    shorten,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
