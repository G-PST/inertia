package internal

import (
    "time"
)

type UnitCategory struct {
    Name string
    Color string
}

type Region struct {
    Name string
}

type UnitMetadata struct {
    Name string
    Category *UnitCategory
    Region *Region
}

type UnitState struct {
    UnitMetadata
    Committed bool
    H float64 // s
    Rating float64 // MVA
}

type SystemState struct {
    Time time.Time
    Requirement float64
    Units []UnitState
}

type SystemMetadata struct {
    Regions []Region
    Categories []UnitCategory
}

type InertiaReport struct {
    Time time.Time
    Total float64
    Requirement float64
    Categories map[string]float64
}

func (st SystemState) Inertia() InertiaReport {

    total_inertia := 0.0
    category_inertia := make(map[string]float64)

    for _, unit := range st.Units {

        if !unit.Committed { continue }

        unit_inertia := unit.H * unit.Rating

        total_inertia += unit_inertia
        category_inertia[unit.Category.Name] += unit_inertia

    }

    return InertiaReport {
        st.Time,
        total_inertia, st.Requirement,
        category_inertia,
    }

}
