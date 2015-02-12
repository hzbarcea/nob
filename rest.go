package main

import (
    "io"
    "io/ioutil"
    "net/http"
)

type RestService struct {
    Url                string
    Username, Password string
}

var NoHeaders = http.Header{ }
var JsonCall = http.Header{
    "Content-Type": { "application/json" },
    "Accept": { "application/json"},
}

func (service RestService) Call(method string, url string, reqBody io.Reader, httpHeaders http.Header) string {
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