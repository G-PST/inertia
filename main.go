// inertia is a Go package for real-time estimation of
// a power system's inertia levels. It defines software interfaces
// for ingesting and reporting data in real-time.
//
// Unit commitment ("H-constant")-based estimation logic and data interfaces
// are available in the inertia/uc package. PMU-based estimation
// methods are planned as future work.
// 
// System integrators can provide deployment-specfic data ingestion code
// (e.g., developed for use with a specific EMS or historian system) that
// conforms to the stated data interfaces for the desired estimation method.
// Once these input interfaces are
// implemented, ingested data can be automatically processed and reported out
// via the package's real-time visualization framework.
// 
// This package provides two off-the-shelf visualization
// modules in inertia/sink/text and inertia/sink/web, but
// custom implementations of
// the [Visualizer] interface can also be used. Multiple Visualizers can be
// associated with a single real-time data stream, allowing for reporting to
// multiple outputs at the same time, for example logging to a text file while
// also visualizing results in a web browser.
package inertia

import (
    "log"
    "time"
)

type DataSource interface {

    // Types implementing DataSource should have a Metadata method that
    // returns SystemMetadata information.
    Metadata() SystemMetadata

    // Query the DataSource, returning the oldest unseen data if it's
    // available, or an error otherwise (e.g. if no new data is available).
    // This allows for making repeated queries until there's
    // nothing new to report, at which point the user can wait some amount of
    // time before trying again.
    Query() (Snapshot, error)

}

type DataSink interface {

    // Pass in static system parameters to the visualization.
    // Should be called exactly once, before any Updates are provided.
    Init(SystemMetadata) error

    // Pass in a new SystemState to be added to the visualization.
    Update(Snapshot)

}

func Run(source DataSource, sinks []DataSink,
         success_freq time.Duration, fail_freq time.Duration) {

    metadata := source.Metadata()
    log.Printf(
        "Loaded metadata: %v regions, %v unit categories",
        len(metadata.Regions), len(metadata.Categories),
    )

    for _, sink := range sinks {
        err := sink.Init(metadata)
        if err != nil {
            log.Fatal("Visualizer initialization error: ", err)
        }
    }

    for {

        snapshot, err := source.Query()
        if err != nil {
            // TODO: Only report non-waiting errors
            log.Print("Query error: ", err)
            time.Sleep(fail_freq)
            continue
        }

        for _, sink := range sinks {
            sink.Update(snapshot)
        }

        time.Sleep(success_freq)

    }

}
