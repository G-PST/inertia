package internal

import (
    "time"
)

type UnitMetadata struct {
    Name string
    Category string
    Region string
}

type UnitState struct {
    UnitMetadata
    Committed bool
    H float64 // s
    Rating float64 // MVA
}

type SystemState struct {
    Time time.Time
    Units []UnitState
}

type InertiaReport struct {
    Total float64
    Categories map[string]float64
}

func (st SystemState) Inertia() InertiaReport {

    total_inertia := 0.0
    category_inertia := make(map[string]float64)

    for _, unit := range st.Units {

        if !unit.Committed { continue }

        unit_inertia := unit.H * unit.Rating
        unitcategory_inertia := category_inertia[unit.Category]

        total_inertia += unit_inertia
        category_inertia[unit.Category] = unitcategory_inertia + unit_inertia

    }

    return InertiaReport { total_inertia, category_inertia }

}
