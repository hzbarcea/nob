package main

import (
    "encoding/json"
    "time"
)

type NobService struct {
    RestService
}

// Broker resource representation in activemq-nob-api
type Broker struct {
    Id                string
    Name              string
    Status            string
    LastModifiedXbean time.Time
}

// used in parsing weirdly embedded json responses
type brokerList struct {
    BrokerList []Broker `json:"broker"`
}
type listBrokersResponse struct {
    Brokers brokerList
}
type brokerInfoResponse struct {
    Broker Broker
}

//// NoB service API

func (service NobService) ListBrokers(filter string) []Broker {
    url := service.Url + "/brokers"
    if filter != "" {
        url += "?filter="
        url += filter // hopefully the http object will encode it
    }

    jsonString := service.Call("GET", url, nil, JsonCall)

    var brokerList listBrokersResponse
    err := json.Unmarshal([]byte(jsonString), &brokerList)
    if err != nil {
        panic(err)
    }

    return brokerList.Brokers.BrokerList
}

func (service NobService) CreateBroker() string {
    url := service.Url + "/brokers?create"

    brokerId := service.Call("POST", url, nil, NoHeaders)
    return brokerId
}

func (service NobService) BrokerInfo(id string) Broker {
    url := service.Url + "/broker/" + id

    jsonString := service.Call("GET", url, nil, JsonCall)

    var brokerInfo brokerInfoResponse
    err := json.Unmarshal([]byte(jsonString), &brokerInfo)
    if err != nil {
        panic(err)
    }

    return brokerInfo.Broker
}
