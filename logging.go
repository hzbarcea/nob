package main

import (
    "io"
    "log"
)

var (
    TraceLog   *log.Logger
    InfoLog    *log.Logger
    ErrorLog   *log.Logger
)

func Init(traceHandle io.Writer, infoHandle io.Writer, errorHandle io.Writer) {

    TraceLog = log.New(traceHandle,
        "TRACE: ",
                log.Ldate|log.Ltime|log.Lshortfile)

    InfoLog = log.New(infoHandle,
        "INFO: ",
                log.Ldate|log.Ltime|log.Lshortfile)

    ErrorLog = log.New(errorHandle,
        "ERROR: ",
                log.Ldate|log.Ltime|log.Lshortfile)
}