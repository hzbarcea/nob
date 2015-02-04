package main

import (
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
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

func (service NobService) call(method string, url string, reqBody io.Reader, httpHeaders http.Header) string {
    fmt.Println("request:", method, url)

    req, err := http.NewRequest(method, url, reqBody)
    if err != nil {
        panic(err)
    }

    if service.Username != "" {
        fmt.Println("request: auth HTTP Basic as: ", service.Username)
        req.SetBasicAuth(service.Username, service.Password)
    }

    for header, values := range httpHeaders {
        for it := range values {
            req.Header.Set(header, values[it])
        }
    }

    fmt.Println("request: Headers:", req.Header)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response: Status:", resp.Status)
    fmt.Println("response: Headers:", resp.Header)
    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    fmt.Println("response: Body:", string(respBody))

    return string(respBody)
}

func (service NobService) ListBrokers(filter string) string {
    url := service.Url + "/brokers"
    if filter != "" {
        url += "?filter="
        url += filter // hopefully the http object will encode it
    }

    return service.call("GET", url, nil, jsonCall)
}

func (service NobService) CreateBroker() string {
    url := service.Url + "/brokers?create"

    return service.call("POST", url, nil, noHeaders)
}

func (service NobService) BrokerInfo(id string) string {
    url := service.Url + "/broker/" + id

    return service.call("GET", url, nil, jsonCall)
}
