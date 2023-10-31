package inertia

import (
    "testing"
    "time"
)

func TestInertiaCalculation(t *testing.T) {

    categories := map[string]*UnitCategory {
        "C1": { "C1", "#00FF00", 1 },
        "C2": { "C2", "#0000FF", 2 },
    }

    regions := map[string]*Region {
        "Region A": { "Region A" },
        "Region B": { "Region B" },
    }

    units := map[string]UnitMetadata {
        "U1": { "U1", categories["C1"], regions["Region A"], 100, 10 },
        "U2": { "U2", categories["C1"], regions["Region B"], 50, 5 },
        "U3": { "U3", categories["C2"], regions["Region A"], 100, 10 },
        "U4": { "U4", categories["C2"], regions["Region B"], 100, 1 },
    }

    system := SystemMetadata { regions, categories, units }

    unitstates := []UnitState {
        { units["U1"], true },
        { units["U2"], true },
        { units["U3"], false },
        { units["U4"], true },
    }

    state := SystemState { time.Now(), 1000, unitstates, &system }
    inertia, _ := state.Inertia()

    inertia.Total.Test(t, 3, 250, 1350)

    inertia.Categories["C1"].Test(t, 2, 150, 1250)
    inertia.Categories["C2"].Test(t, 1, 100, 100)

    inertia.Regions["Region A"].Test(t, 1, 100, 1000)
    inertia.Regions["Region B"].Test(t, 2, 150, 350)

    // Empty region

    state.Units[0].Committed = false
    inertia, _ = state.Inertia()

    inertia.Total.Test(t, 2, 150, 350)

    inertia.Categories["C1"].Test(t, 1, 50, 250)
    inertia.Categories["C2"].Test(t, 1, 100, 100)

    regionA, ok := inertia.Regions["Region A"]
    if !ok {
        t.Errorf("Results should exist even when no units are committed")
    }

    regionA.Test(t, 0, 0, 0)
    inertia.Regions["Region B"].Test(t, 2, 150, 350)

    // Empty category and region

    state.Units[3].Committed = false
    inertia, _ = state.Inertia()

    inertia.Total.Test(t, 1, 50, 250)

    c2, ok := inertia.Categories["C2"]
    if !ok {
        t.Errorf("Results should exist even when no units are committed")
    }

    inertia.Categories["C1"].Test(t, 1, 50, 250)
    c2.Test(t, 0, 0, 0)

    regionA, ok = inertia.Regions["Region A"]
    if !ok {
        t.Errorf("Results should exist even when no units are committed")
    }

    regionA.Test(t, 0, 0, 0)
    inertia.Regions["Region B"].Test(t, 1, 50, 250)

}

func (agg UnitAggregation) Test(t *testing.T, units int, rating float64, inertia float64) {

    if i := agg.Units; i != units {
        t.Errorf("Result should have %d committed units; got %d", units, i)
    }

    if i := agg.TotalRating; i != rating {
        t.Errorf("Result should have total rating of %f; got %f", rating, i)
    }

    if i := agg.TotalInertia; i != inertia {
        t.Errorf("Result should have total inertia of %f; got %f", inertia, i)
    }

}
