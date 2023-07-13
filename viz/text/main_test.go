package text

import (
    "testing"
    "time"

    "github.com/G-PST/inertia"
    d "github.com/G-PST/inertia/uc/mockdata"
)

func TestRun(t *testing.T) {

    datasource := d.New(2 * time.Second)
    vizs := []inertia.Visualizer { New() }

    inertia.Run(datasource, vizs, 500 * time.Millisecond, time.Second)

}

