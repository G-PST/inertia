package main

import (
    "time"

    "github.com/G-PST/inertia"
    "github.com/G-PST/inertia/sinks/web"
    d "github.com/G-PST/inertia/uc/mockdata"
)

func main() {

    datasource := d.New(10 * time.Second)
    sinks := []inertia.DataSink { web.New(":8181") }

    inertia.Run(datasource, sinks, 500 * time.Millisecond, time.Second)

}
