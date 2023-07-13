package uc

import (
    "time"
    "github.com/G-PST/inertia"
)

type UnitState struct {
    inertia.UnitMetadata
    Committed bool
    H float64 // s
    Rating float64 // MVA
}

type SystemState struct {
    Time time.Time
    Requirement float64
    Units []UnitState
}

// TODO: Always report all categories, even when SystemState doesn't have
// any for that category (committed or otherwise)
func (st SystemState) Inertia() (inertia.Snapshot, error) {

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

    return inertia.Snapshot{
        st.Time,
        total_inertia, st.Requirement,
        category_inertia,
    }, nil

}
