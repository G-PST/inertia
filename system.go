package inertia

import (
    "time"
)

type UnitCategory struct {
    Name string `json:"name"`
    Color string `json:"color"`
    Order int `json:"order"`
}

type Region struct {
    Name string `json:"name"`
}

type SystemMetadata struct {
    Regions []Region `json:"regions"`
    Categories []UnitCategory `json:"categories"`
}

type UnitMetadata struct {
    Name string
    Category *UnitCategory
    Region *Region
}

// The SystemState interface represents an abstract collection of power system
// state data necessary for a particular inertia estimation method. This
// interface only needs to be implemented when adding new estimation methods,
// not when adding new DataSources for an existing estimation method.
type SystemState interface {

    // Types implementing SystemState should have an Inertia method that
    // converts internal state data into inertia levels for reporting.
    Inertia() (Snapshot, error)

}

// Snapshot is the common data structure used for reporting inertia
// levels at a point in time.
type Snapshot struct {
    Time time.Time
    Total float64
    Requirement float64
    Categories map[string]float64
}
