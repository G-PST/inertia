package web

import (
    "testing"
    "time"

    "github.com/G-PST/inertia"
    d "github.com/G-PST/inertia/datasource/dummy"
)

func TestRun(t *testing.T) {

    datasource := d.New(10 * time.Second)
    vizs := []inertia.Visualizer { New(":8181") }

    inertia.Run(datasource, vizs, 500 * time.Millisecond, time.Second)

}
