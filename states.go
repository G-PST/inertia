package inertia

import (
    "time"
)

type UnitState struct {
    name string
    category string
    committed bool
    h float64 // s
    rating float64 // MVA
}

type SystemState struct {
    time time.Time
    units []UnitState
}

type InertiaReport struct {
    Total float64
    Categories map[string]float64
}

func (st SystemState) Inertia() InertiaReport {

    total_inertia := 0.0
    category_inertia := make(map[string]float64)

    for _, unit := range st.units {

        if !unit.committed { continue }

        unit_inertia := unit.h * unit.rating
        unitcategory_inertia := category_inertia[unit.category]

        total_inertia += unit_inertia
        category_inertia[unit.category] = unitcategory_inertia + unit_inertia

    }

    return InertiaReport { total_inertia, category_inertia }

}
