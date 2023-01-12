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
        { "U1", "C1", randBool(), 10, 100 },
        { "U2", "C1", randBool(), 5, 50 },
        { "U3", "C2", randBool(), 10, 100 },
        { "U4", "C2", randBool(), 1, 100 },
    }
    
    return internal.SystemState { now, units }, nil

}

func randBool() bool {
    return rand.Float64() < 0.9
}
