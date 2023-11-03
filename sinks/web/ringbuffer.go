package web

import (
    "errors"
    "time"

    "github.com/G-PST/inertia"
)

// Important assumptions:
// Snapshots are always added in chronological order
// Elements cannot be removed, only added (such that the ring is always full
// if ring.first > 0)
type SnapshotRing struct {
    buffer []inertia.Snapshot
    first int
    last int
    length int
}

func NewSnapshotRing(length int) SnapshotRing {

    buffer := make([]inertia.Snapshot, length, length)
    return SnapshotRing { buffer, -1, -1, length }

}

func (ring *SnapshotRing) Push(state inertia.Snapshot) {

    isEmpty := ring.first < 0
    isFull := ring.first > 0 || (ring.last + 1) == ring.length

    if (isFull || isEmpty) {

        ring.first += 1
        ring.first %= ring.length

    }

    ring.last += 1
    ring.last %= ring.length

    ring.buffer[ring.last] = state

}

func (ring SnapshotRing) FirstAfter(t time.Time) (inertia.Snapshot, error) {

    empty_err := errors.New("No newer states available")

    i := ring.first

    if i < 0 {
        return inertia.Snapshot {}, empty_err
    }

    for {

        state := ring.buffer[i]

        if state.Time.After(t) {
            return state, nil
        }

        i += 1

        if (i == ring.length) {
            i = 0
        }

        if (i == ring.first) {
            break
        }

    }

    return inertia.Snapshot {}, empty_err

}
