package inertia

import (
    "github.com/G-PST/inertia/internal"
)

type DataSource interface {
    // Query the datasource, returning the oldest unseen data if it's
    // available, or an error otherwise (e.g. if no new data is available).
    // This allows for making repeated queries until there's
    // nothing new to report, at which point the user can wait some amount of
    // time before trying again.
    Query() (internal.SystemState, error)
}

type Visualizer interface {
    // Pass in a new SystemState to be added to the visualization
    Update(internal.SystemState)
}
