package main

import (
  "os"
  "github.com/codegangsta/cli"
)

func main() {
  app := cli.NewApp()
  app.Name = "nob"
  app.Usage = "ActiveMQ network of brokers manager"
  app.Version = "1.0"
  app.Author = "Hadrian Zbarcea"
  app.Email = "hadrian@apache.org"
  app.Commands = []cli.Command{
    {
      Name:  "login",
      ShortName:   "l",
      Usage: "establish nob session",
      Description: "This is how we describe hello the function",
      Action: func(c *cli.Context) {
        println("login")
      },
    }, {
      Name:        "broker",
      ShortName:   "b",
      Usage:       "usage of the broker command",
      Description: "Describe the usage of the broker command",
      Subcommands: []cli.Command{
        {
          Name:        "list",
          ShortName:   "l",
          Usage:       "print a list of brokers",
          Description: "Description of how the list of brokers is printed",
          Flags: []cli.Flag{
            cli.StringFlag{
              Name:  "filter",
              Value: "",
              Usage: "Filter the brokers to print",
            },
          },
          Action: func(c *cli.Context) {
            println("Listing brokers: , ", c.String("filter"))
          },
        }, {
          Name:        "create",
          ShortName:   "c",
          Usage:       "create a new broker",
          Description: "Description of how to create a new broker",
          Flags: []cli.Flag{
            cli.StringFlag{
              Name:  "name",
              Value: "",
              Usage: "Name of the new broker",
            },
          },
          Action: func(c *cli.Context) {
            println("Creating broker: , ", c.String("name"))
          },
        }, {
          Name:      "info",
          ShortName: "i",
          Usage:     "prints broker metadata",
          Flags: []cli.Flag{
            cli.StringFlag{
              Name:  "id",
              Value: "",
              Usage: "id of the broker",
            },
          },
          Action: func(c *cli.Context) {
            println("Broker info: , ", c.String("id"))
          },
        },
      },
    },
  }

  app.Run(os.Args)
}
