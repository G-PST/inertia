package main

import (
    "time"

    "github.com/G-PST/inertia"
    "github.com/G-PST/inertia/viz/web"
    d "github.com/G-PST/inertia/uc/mockdata"
)

func main() {

    datasource := d.New(10 * time.Second)
    vizs := []inertia.Visualizer { web.New(":8181") }

    inertia.Run(datasource, vizs, 500 * time.Millisecond, time.Second)

}
