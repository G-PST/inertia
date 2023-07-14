package main

import (
    "time"

    "github.com/G-PST/inertia"
    "github.com/G-PST/inertia/viz/text"
    d "github.com/G-PST/inertia/uc/mockdata"
)

func main() {

    datasource := d.New(2 * time.Second)
    vizs := []inertia.Visualizer { text.New() }

    inertia.Run(datasource, vizs, 500 * time.Millisecond, time.Second)

}

