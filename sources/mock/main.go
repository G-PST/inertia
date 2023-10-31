package mockdata

import (
    "errors"
    "math/rand"
    "time"

    "github.com/G-PST/inertia"
    "github.com/G-PST/inertia/uc"
)

type MockDataSource struct {

    freq time.Duration
    lastTime time.Time

    system inertia.SystemMetadata

}

func New(freq time.Duration) *MockDataSource {

    regions := map[string]*inertia.Region {
        "Region A": &inertia.Region { "Region A" },
        "Region B": &inertia.Region { "Region B" },
    }

    categories := map[string]*inertia.UnitCategory {
        "C1": &inertia.UnitCategory { "C1", "#00FF00", 1 },
        "C2": &inertia.UnitCategory { "C2", "#0000FF", 2 },
    }

    units := map[string]inertia.UnitMetadata {
        "U1": inertia.UnitMetadata {
            "U1", categories["C1"], regions["Region A"], 100, 10 },
        "U2": inertia.UnitMetadata {
            "U2", categories["C1"], regions["Region B"], 50, 5 },
        "U3": inertia.UnitMetadata {
            "U3", categories["C2"], regions["Region A"], 100, 10 },
        "U4": inertia.UnitMetadata {
            "U4", categories["C2"], regions["Region B"], 100, 1 },
    }

    return &MockDataSource {
        freq: freq,
        system: inertia.SystemMetadata { regions, categories, units },
    }

}

func (d *MockDataSource) Metadata() inertia.SystemMetadata {

    return d.system

}

func (d *MockDataSource) Query() (inertia.Snapshot, error) {

    now := time.Now()
    next_result := d.lastTime.Add(d.freq)

    if now.Before(next_result) {
        return inertia.Snapshot {}, errors.New("No new data")
    }

    d.lastTime = now

    units := []uc.UnitState {
        { d.system.Units["U1"], randBool() },
        { d.system.Units["U2"], randBool() },
        { d.system.Units["U3"], randBool() },
        { d.system.Units["U4"], randBool() },
    }

    state := uc.SystemState { now, 1500, units }

    return state.Inertia()

}

func randBool() bool {
    return rand.Float64() < 0.9
}
