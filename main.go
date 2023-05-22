package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/EwanValentine/eze-rpc/generator"
)

func main() {
	app := &cli.App{
		Name:  "eze",
		Usage: "I dunno, look at the code and figure it out",
		Action: func(*cli.Context) error {
			fmt.Println("Eze RPC, apparently.")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "generate",
				Usage: "Generate code from a DSL",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "input",
						Value: "example/example.eze",
						Usage: "input file to read from: --input=example/example.eze",
					},
					&cli.StringFlag{
						Name:  "output",
						Value: "users",
						Usage: "output directory to write to: --output=.",
					},
				},
				Action: func(c *cli.Context) error {
					input := c.String("input")
					output := c.String("output")

					content, err := os.ReadFile(input)
					if err != nil {
						return errors.New("error reading input file, should be like: 'eze generate --input=example/example.eze --output=example/'")
					}

					dsl := generator.ParseDSL(string(content))

					outputCode := generator.GenerateCode(dsl)

					if err := os.MkdirAll(fmt.Sprintf("%s/%s/", output, dsl.Package), 0755); err != nil {
						return err
					}

					if err := os.WriteFile(fmt.Sprintf("%s/%s/%s", output, dsl.Package, "eze_rpc.go"), []byte(outputCode), 0644); err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
