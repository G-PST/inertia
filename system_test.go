package inertia

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
        { UnitMetadata {"U1", &categories[0], &regions[0], 10, 100 }, true },
        { UnitMetadata {"U2", &categories[0], &regions[1], 5, 50 }, true },
        { UnitMetadata {"U3", &categories[1], &regions[0], 10, 100 }, false },
        { UnitMetadata {"U4", &categories[1], &regions[1], 1, 100 }, true },
    }

    state := SystemState { time.Now(), 1000, units }
    inertia, _ := state.Inertia()

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
