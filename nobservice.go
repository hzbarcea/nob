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

func call(service NobService, method string, url string, reqBody io.Reader, httpHeaders http.Header) string {
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

func ListBrokers(service NobService, filter string) string {
    url := service.Url + "/brokers"
    if filter != "" {
        url += "?filter="
        url += filter // hopefully the http object will encode it
    }

    return call(service, "GET", url, nil, jsonCall)
}

func CreateBroker(service NobService) string {
    url := service.Url + "/brokers?create"

    return call(service, "POST", url, nil, noHeaders)
}

func BrokerInfo(service NobService, id string) string {
    url := service.Url + "/broker/" + id

    return call(service, "GET", url, nil, jsonCall)
}
