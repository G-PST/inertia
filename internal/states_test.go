package internal

import (
    "testing"
    "time"
)

func TestInertiaCalculation(t *testing.T) {

    units := []UnitState {
        { "U1", "C1", true, 10, 100 },
        { "U2", "C1", true, 5, 50 },
        { "U3", "C2", false, 10, 100 },
        { "U4", "C2", true, 1, 100 },
    }
    state := SystemState { time.Now(), units }
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
