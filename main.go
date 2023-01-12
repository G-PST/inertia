package inertia

import (
    "time"
)

func Run(source DataSource, vizs []Visualizer, freq time.Duration) {

    for {

        state, err := source.Query()
        if err != nil {
            time.Sleep(freq)
            continue
        }

        for _, viz := range vizs {
            viz.Update(state)
        }

    }
}
