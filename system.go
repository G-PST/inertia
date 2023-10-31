package inertia

import (
    "time"
)

// UnitCategory provides metadata for logical groupings of generating units
type UnitCategory struct {
    Name string `json:"name"`
    Color string `json:"color"`
    Order int `json:"order"`
}

// Region provides metadata on mutually-exclusive subsets of the full network
type Region struct {
    Name string `json:"name"`
}

// UnitMetadata provides parameters and classifications for specific
// generating units
type UnitMetadata struct {
    Name string
    Category *UnitCategory
    Region *Region
    Rating float64 // in MVA
    H float64 // in s
}

// SystemMetadata brings together metadata about different aspects of the
// system
type SystemMetadata struct {
    Regions map[string]*Region `json:"regions"`
    Categories map[string]*UnitCategory `json:"categories"`
    Units map[string]UnitMetadata
}

type UnitState struct {
    UnitMetadata
    Committed bool
}

type SystemState struct {
    Time time.Time
    Requirement float64
    Units []UnitState
}

// Snapshot is the common data structure used for reporting inertia
// levels at a point in time.
type Snapshot struct {
    Time time.Time
    Total float64
    Requirement float64
    Categories map[string]float64
}

// TODO: Always report all categories, even when SystemState doesn't have
// any for that category (committed or otherwise)
func (st SystemState) Inertia() (Snapshot, error) {

    total_inertia := 0.0
    category_inertia := make(map[string]float64)

    for _, unit := range st.Units {

        var unit_inertia float64

        if unit.Committed {
            unit_inertia = unit.H * unit.Rating
        }

        total_inertia += unit_inertia
        category_inertia[unit.Category.Name] += unit_inertia

    }

    return Snapshot{
        st.Time,
        total_inertia, st.Requirement,
        category_inertia,
    }, nil

}
