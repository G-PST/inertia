package inertia

import (
    "log"
    "time"
)

func Run(source DataSource, vizs []Visualizer,
         success_freq time.Duration, fail_freq time.Duration) {

    metadata := source.Metadata()

    for _, viz := range vizs {
        err := viz.Init(metadata)
        if err != nil {
            log.Fatal("Visualizer initialization error: ", err)
        }
    }

    for {

        state, err := source.Query()
        if err != nil {
            time.Sleep(fail_freq)
            continue
        }

        for _, viz := range vizs {
            viz.Update(state)
        }

        time.Sleep(success_freq)

    }

}
