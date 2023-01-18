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
}

func New(freq time.Duration) *DummyDataSource {
    return &DummyDataSource { freq: freq }
}

func (d *DummyDataSource) Query() (internal.SystemState, error) {

    now := time.Now()
    next_result := d.lastTime.Add(d.freq)

    if now.Before(next_result) {
        return internal.SystemState{}, errors.New("No new data")
    }

    d.lastTime = now

    units := []internal.UnitState {

        { internal.UnitMetadata { "U1", "C1", "R1"},
          randBool(), 10, 100 },

        { internal.UnitMetadata { "U2", "C1", "R2" },
          randBool(), 5, 50 },

        { internal.UnitMetadata { "U3", "C2", "R1" },
          randBool(), 10, 100 },

        { internal.UnitMetadata { "U4", "C2", "R2" },
          randBool(), 1, 100 },

    }
    
    return internal.SystemState { now, units }, nil

}

func randBool() bool {
    return rand.Float64() < 0.9
}
