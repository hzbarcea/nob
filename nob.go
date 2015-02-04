package main

import (
    "github.com/codegangsta/cli"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "log"
    "os"
)

func createNobService(c *cli.Context) (NobService, error) {
    configFile := c.GlobalString("config")
    contents, err := ioutil.ReadFile(configFile)
    if err != nil {
        log.Fatalln("unable to read config file:", err)
        return NobService{}, err
    }

    var config map[string]NobService
    config = make(map[string]NobService)
    err = yaml.Unmarshal(contents, config)
    if err != nil {
        log.Fatalln("cannot parse ", configFile, ":", err)
        return NobService{}, err
    }

    nob := c.GlobalString("nob")
    service := config[nob]
    if service.Url == "" {
        log.Fatalln("no nob service '"+nob+"' defined in file ", configFile)
        return NobService{}, err
    }

    return service, nil
}

func main() {
    app := cli.NewApp()
    app.Name = "nob"
    app.Usage = "ActiveMQ network of brokers manager"
    app.Version = "1.0"
    app.Author = "Hadrian Zbarcea"
    app.Email = "hadrian@apache.org"
    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:  "config, f",
            Value: os.Getenv("HOME")+"/.nobrc",
            Usage: "nob client configuration file",
        },
        cli.StringFlag{
            Name:  "nob, n",
            Value: "",
            Usage: "nob service definition from the config file",
        },
    }
    app.Commands = []cli.Command{
        {
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
                        service, err := createNobService(c)
                        if err != nil {
                            println("could not fetch broker list:", err)
                            return
                        }
                        ListBrokers(service, c.String("filter"))
                    },
                }, {
                    Name:        "create",
                    ShortName:   "c",
                    Usage:       "create a new broker",
                    Description: "Description of how to create a new broker",
                    Action: func(c *cli.Context) {
                        service, err := createNobService(c)
                        if err != nil {
                            println("could not create broker:", err)
                            return
                        }
                        CreateBroker(service)
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
                        service, err := createNobService(c)
                        if err != nil {
                            println("could not get broker info:", err)
                            return
                        }
                        BrokerInfo(service, c.String("id"))
                    },
                },
            },
        },
    }

    app.Run(os.Args)
}
