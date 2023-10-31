package main

import (
    "time"

    "github.com/G-PST/inertia"
    "github.com/G-PST/inertia/sinks/text"
    d "github.com/G-PST/inertia/sources/mock"
)

func main() {

    datasource := d.New(2 * time.Second)
    sinks := []inertia.DataSink{ text.New() }

    inertia.Run(datasource, sinks, 500 * time.Millisecond, time.Second)

}

