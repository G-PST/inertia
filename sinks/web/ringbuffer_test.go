package web

import (
    "testing"
    "time"

    "github.com/G-PST/inertia"
)

func TestRingBuffer(t *testing.T) {

    ring := NewSnapshotRing(2)
    t0 := time.Time {}

    _, err := ring.FirstAfter(t0)
    if err == nil {
        t.Errorf("Query should report no result when ring is empty")
    }

    state1 := inertia.Snapshot { Time: time.Now() }
    ring.Push(state1)

    state_result, err := ring.FirstAfter(t0)
    if state_result.Time != state1.Time {
        t.Errorf("Result should match only contained state")
    }

    _, err = ring.FirstAfter(state1.Time)
    if err == nil {
        t.Errorf("Query should report no result when states are older")
    }

    state2 := inertia.Snapshot { Time: state1.Time.Add(5000000000) }
    ring.Push(state2)

    state_result, err = ring.FirstAfter(t0)
    if state_result.Time != state1.Time {
        t.Errorf("Result should match oldest contained state")
    }

    state_result, err = ring.FirstAfter(state1.Time)
    if state_result.Time != state2.Time {
        t.Errorf("Result should match newest contained state")
    }

    _, err = ring.FirstAfter(state2.Time)
    if err == nil {
        t.Errorf("Query should report no result when states are older")
    }

    state3 := inertia.Snapshot { Time: state2.Time.Add(5000000000) }
    ring.Push(state3)

    state_result, err = ring.FirstAfter(t0)
    if state_result.Time != state2.Time {
        t.Errorf("Result should match oldest contained state")
    }

    state_result, err = ring.FirstAfter(state1.Time)
    if state_result.Time != state2.Time {
        t.Errorf("Result should match oldest contained state")
    }

    state_result, err = ring.FirstAfter(state2.Time)
    if state_result.Time != state3.Time {
        t.Errorf("Result should match newest contained state")
    }

    _, err = ring.FirstAfter(state3.Time)
    if err == nil {
        t.Errorf("Query should report no result when states are older")
    }

    state4 := inertia.Snapshot { Time: state3.Time.Add(5000000000) }
    ring.Push(state4)

    state_result, err = ring.FirstAfter(t0)
    if state_result.Time != state3.Time {
        t.Errorf("Result should match oldest contained state")
    }

    state_result, err = ring.FirstAfter(state1.Time)
    if state_result.Time != state3.Time {
        t.Errorf("Result should match oldest contained state")
    }

    state_result, err = ring.FirstAfter(state2.Time)
    if state_result.Time != state3.Time {
        t.Errorf("Result should match newest contained state")
    }

    state_result, err = ring.FirstAfter(state3.Time)
    if state_result.Time != state4.Time {
        t.Errorf("Result should match newest contained state")
    }

    _, err = ring.FirstAfter(state4.Time)
    if err == nil {
        t.Errorf("Query should report no result when states are older")
    }

}
