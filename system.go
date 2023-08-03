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
    Rating float64 // in MVA
    H float64 // in s
}

// Snapshot is the common data structure used for reporting inertia
// levels at a point in time.
type Snapshot struct {
    Time time.Time
    Total float64
    Requirement float64
    Categories map[string]float64
}
