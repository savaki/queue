package main

import (
	"github.com/codegangsta/cli"
	"os"
)

var globalFlags []cli.Flag = []cli.Flag{
	cli.BoolFlag{"verbose", "increase output"},
	cli.StringFlag{"provider, p", "sqs", "the queue provider"},
	cli.StringFlag{"region, r", "us-west-1", "the region the queue [AWS-only]"},
	cli.StringFlag{"name, n", "", "the name of the queue"},
}

func main() {
	app := cli.NewApp()
	app.Name = "queue"
	app.Author = "Matt Ho"
	app.Version = "0.0.1"
	app.Usage = "CLI interface to enqueue and dequeue messages from a queue"
	app.Commands = []cli.Command{
		{
			Name:   "write",
			Usage:  "write a single message into the queue",
			Action: WriteCommand,
			Flags: append(globalFlags, []cli.Flag{
				cli.StringFlag{"file, f", "", "the file to write to the queue"},
			}...),
		},
		{
			Name:   "directory",
			Usage:  "write an entire directory's content into a queue",
			Action: DirectoryCommand,
			Flags: append(globalFlags, []cli.Flag{
				cli.StringFlag{"dir, d", "", "the directory to write to the queue"},
			}...),
		},
	}

	app.Run(os.Args)
}
