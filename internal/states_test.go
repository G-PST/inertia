package internal

import (
    "testing"
    "time"
)

func TestInertiaCalculation(t *testing.T) {

    categories := []UnitCategory {
        { "C1", "#00FF00", 1 }, { "C2", "#0000FF", 2 },
    }

    regions := []Region { { "Region A" }, { "Region B" } }

    units := []UnitState {
        { UnitMetadata {"U1", &categories[0], &regions[0] }, true, 10, 100 },
        { UnitMetadata {"U2", &categories[0], &regions[1] }, true, 5, 50 },
        { UnitMetadata {"U3", &categories[1], &regions[0] }, false, 10, 100 },
        { UnitMetadata {"U4", &categories[1], &regions[1] }, true, 1, 100 },
    }

    state := SystemState { time.Now(), 1000, units }
    inertia := state.Inertia()

    if i := inertia.Categories["C1"]; i != 1250 {
        t.Errorf("C1 inertia should be 1250 MW s; got %f", i)
    }

    if i := inertia.Categories["C2"]; i != 100 {
        t.Errorf("C2 inertia should be 100 MW s; got %f", i)
    }

    if inertia.Total != 1350 {
        t.Errorf("Total inertia should be 1350 MW s; got %f", inertia.Total)
    }

}
