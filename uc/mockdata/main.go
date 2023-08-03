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

    regions []inertia.Region
    categories []inertia.UnitCategory
    units []inertia.UnitMetadata

}

func New(freq time.Duration) *MockDataSource {

    regions := []inertia.Region {
        { "Region A" }, { "Region B" },
    }

    categories := []inertia.UnitCategory {
        { "C1", "#00FF00", 1 }, { "C2", "#0000FF", 2 },
    }

    units := []inertia.UnitMetadata {
        { "U1", &categories[0], &regions[0], 100, 10 },
        { "U2", &categories[0], &regions[1], 50, 5 },
        { "U3", &categories[1], &regions[0], 100, 10 },
        { "U4", &categories[1], &regions[1], 100, 1 },
    }

    return &MockDataSource {
        freq: freq,
        regions: regions,
        categories: categories,
        units: units,
    }

}

func (d *MockDataSource) Metadata() inertia.SystemMetadata {

    return inertia.SystemMetadata { d.regions, d.categories }

}

func (d *MockDataSource) Query() (inertia.Snapshot, error) {

    now := time.Now()
    next_result := d.lastTime.Add(d.freq)

    if now.Before(next_result) {
        return inertia.Snapshot {}, errors.New("No new data")
    }

    d.lastTime = now

    units := []uc.UnitState {
        { d.units[0], randBool() },
        { d.units[1], randBool() },
        { d.units[2], randBool() },
        { d.units[3], randBool() },
    }

    state := uc.SystemState { now, 1500, units }

    return state.Inertia()

}

func randBool() bool {
    return rand.Float64() < 0.9
}
