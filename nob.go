package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type NobService struct {
	Url                string
	Username, Password string
}

func listBrokers(service NobService, filter string) string {
	url := service.Url + "/brokers"
	if filter != "" {
		url += "?filter="
		url += filter // hopefully the http object will encode it
	}

	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(service.Username, service.Password)

	//	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return string(body)
}

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
			Value: os.Getenv("HOME") + "/.nobrc",
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
						listBrokers(service, c.String("filter"))
					},
				}, {
					Name:        "create",
					ShortName:   "c",
					Usage:       "create a new broker",
					Description: "Description of how to create a new broker",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "name",
							//							Value: "",
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
