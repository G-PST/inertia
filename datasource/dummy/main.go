package dummy

import (
    "errors"
    "math/rand"
    "time"

    "github.com/G-PST/inertia/internal"
)

type DummyDataSource struct {

    freq time.Duration
    lastTime time.Time

    regions []internal.Region
    categories []internal.UnitCategory
    units []internal.UnitMetadata

}

func New(freq time.Duration) *DummyDataSource {

    regions := []internal.Region {
        { "Region A" }, { "Region B" },
    }

    categories := []internal.UnitCategory {
        { "C1", "#00FF00" }, { "C2", "#0000FF" },
    }

    units := []internal.UnitMetadata {
        { "U1", &categories[0], &regions[0] },
        { "U2", &categories[0], &regions[1] },
        { "U3", &categories[1], &regions[0] },
        { "U4", &categories[1], &regions[1] },
    }

    return &DummyDataSource {
        freq: freq,
        regions: regions,
        categories: categories,
        units: units,
    }

}

func (d *DummyDataSource) Metadata() internal.SystemMetadata {

    return internal.SystemMetadata { d.regions, d.categories }

}

func (d *DummyDataSource) Query() (internal.SystemState, error) {

    now := time.Now()
    next_result := d.lastTime.Add(d.freq)

    if now.Before(next_result) {
        return internal.SystemState{}, errors.New("No new data")
    }

    d.lastTime = now

    units := []internal.UnitState {
        { d.units[0], randBool(), 10, 100 },
        { d.units[1], randBool(), 5, 50 },
        { d.units[2], randBool(), 10, 100 },
        { d.units[3], randBool(), 1, 100 },
    }
    
    return internal.SystemState { now, 1500, units }, nil

}

func randBool() bool {
    return rand.Float64() < 0.9
}
