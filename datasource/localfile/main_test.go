package localfile

import (
    "testing"
    "time"
)

func TestLocalFileLoadOrder(t *testing.T) {

    lf, err := New(
        "test/rtag/states", "UNInertia_2006-01-02-15-04.csv",
        "test/rtag/metadata")

    if err != nil {
        t.Fatalf("Error creating LocalFile: %v", err)
    }

    for i := 0; i < 5; i += 1 {

        state, err := lf.Query()

        if err != nil {
            t.Fatalf("Unexpected error querying SystemState: %v", err)
        }

        if state.Time != newTime(2022, time.December, 5, 13, 35 + 5*i) {
            t.Fatalf("SystemState loaded with unexpected time")
        }

    }

    _, err = lf.Query()
    if err == nil {
        t.Fatalf("SystemState loaded successfully while states exhausted")
    }

}

func newTime(year int, month time.Month, day int, hour int, minute int) time.Time {
    return time.Date(year, month, day, hour, minute, 0, 0, time.Local)
}
