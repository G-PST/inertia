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
    System *SystemMetadata
}

type UnitAggregation struct {
    Units int `json:"units"`
    TotalRating float64 `json:"total_rating"`
    TotalInertia float64 `json:"total_inertia"`
}

func (agg *UnitAggregation) AddUnit(h float64, rating float64) {
    agg.Units += 1
    agg.TotalRating += rating
    agg.TotalInertia += h * rating
}

// Snapshot is the common data structure used for reporting inertia
// levels at a point in time.
type Snapshot struct {
    Time time.Time
    Requirement float64
    Total UnitAggregation
    Categories map[string]*UnitAggregation
    Regions map[string]*UnitAggregation
}

func (st SystemState) Inertia() (Snapshot, error) {

    total_inertia := UnitAggregation {}
    category_inertias := make(map[string]*UnitAggregation)
    region_inertias := make(map[string]*UnitAggregation)

    for region, _ := range st.System.Regions {
        region_inertias[region] = &UnitAggregation {}
    }

    for category, _ := range st.System.Categories {
        category_inertias[category] = &UnitAggregation {}
    }

    for _, unit := range st.Units {

        if !unit.Committed { continue }

        total_inertia.AddUnit(unit.H, unit.Rating)

        category := unit.Category.Name
        category_inertias[category].AddUnit(unit.H, unit.Rating)

        region := unit.Region.Name
        region_inertias[region].AddUnit(unit.H, unit.Rating)

    }

    return Snapshot {
        st.Time,
        st.Requirement,
        total_inertia,
        category_inertias,
        region_inertias,
    }, nil

}
