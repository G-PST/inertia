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
        { "U1", &categories[0], &regions[0] },
        { "U2", &categories[0], &regions[1] },
        { "U3", &categories[1], &regions[0] },
        { "U4", &categories[1], &regions[1] },
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

func (d *MockDataSource) Query() (inertia.SystemState, error) {

    now := time.Now()
    next_result := d.lastTime.Add(d.freq)

    if now.Before(next_result) {
        return uc.SystemState{}, errors.New("No new data")
    }

    d.lastTime = now

    units := []uc.UnitState {
        { d.units[0], randBool(), 10, 100 },
        { d.units[1], randBool(), 5, 50 },
        { d.units[2], randBool(), 10, 100 },
        { d.units[3], randBool(), 1, 100 },
    }
    
    return uc.SystemState { now, 1500, units }, nil

}

func randBool() bool {
    return rand.Float64() < 0.9
}
