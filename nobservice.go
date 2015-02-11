package main

import (
    "encoding/json"
    "io"
    "io/ioutil"
    "net/http"
    "time"
)

type NobService struct {
    Url                string
    Username, Password string
}

var noHeaders = http.Header{ }
var jsonCall = http.Header{
    "Content-Type": { "application/json" },
    "Accept": { "application/json"},
}

// Broker resource representation in activemq-nob-api
type Broker struct {
    Id                string
    Name              string
    Status            string
    LastModifiedXbean time.Time
}

// used in parsing weirdly embedded json responses :D
type brokerList struct {
    BrokerList []Broker `json:"broker"`
}
type listBrokersResponse struct {
    Brokers brokerList
}
type brokerInfoResponse struct {
    Broker Broker
}


func (service NobService) call(method string, url string, reqBody io.Reader, httpHeaders http.Header) string {
    TraceLog.Println("request:", method, url)

    req, err := http.NewRequest(method, url, reqBody)
    if err != nil {
        panic(err)
    }

    if service.Username != "" {
        TraceLog.Println("request: auth HTTP Basic as: ", service.Username)
        req.SetBasicAuth(service.Username, service.Password)
    }

    for header, values := range httpHeaders {
        for it := range values {
            req.Header.Set(header, values[it])
        }
    }

    TraceLog.Println("request: Headers:", req.Header)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    TraceLog.Println("response: Status:", resp.Status)
    TraceLog.Println("response: Headers:", resp.Header)
    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    TraceLog.Println("response: Body:", string(respBody))

    return string(respBody)
}

func (service NobService) ListBrokers(filter string) []Broker {
    url := service.Url + "/brokers"
    if filter != "" {
        url += "?filter="
        url += filter // hopefully the http object will encode it
    }

    jsonString := service.call("GET", url, nil, jsonCall)

    var brokerList listBrokersResponse
    err := json.Unmarshal([]byte(jsonString), &brokerList)
    if err != nil {
        panic(err)
    }

    return brokerList.Brokers.BrokerList

}

func (service NobService) CreateBroker() string {
    url := service.Url + "/brokers?create"

    brokerId := service.call("POST", url, nil, noHeaders)
    return brokerId
}

func (service NobService) BrokerInfo(id string) Broker {
    url := service.Url + "/broker/" + id

    jsonString := service.call("GET", url, nil, jsonCall)

    var brokerInfo brokerInfoResponse
    err := json.Unmarshal([]byte(jsonString), &brokerInfo)
    if err != nil {
        panic(err)
    }

    return brokerInfo.Broker
}
