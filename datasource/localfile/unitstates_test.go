package localfile

import (
    "github.com/G-PST/inertia/internal"

    "os"
    "reflect"
    "testing"
)

func TestUnitStateParsing(t *testing.T) {

    statefile, err := os.Open("test/states/sample_unitstates.csv")
    if err != nil {
        t.Fatalf("Error loading unit states sample: %v", err)
    }

    unitstates, err := parseUnitStates(statefile, metadata_ref)
    if err != nil {
        t.Fatalf("Error parsing unit states sample: %v", err)
    }

    if !reflect.DeepEqual(unitstates, units_ref) {
        t.Fatalf("Parsed unit states sample did not match expected value")
    }

}

var units_ref = []internal.UnitState{

    { UnitMetadata: metadata_ref["29EDD18_GEN"],
      Committed: true,
      H: 6.0599999,
      Rating: 104 },

    { UnitMetadata: metadata_ref["29EDD2_G3"],
      Committed: false,
      H: 2.8499999,
      Rating: 76 },

    { UnitMetadata: metadata_ref["29EDD3_G5"],
      Committed: true,
      H: 7.4099998,
      Rating: 146 },

    { UnitMetadata: metadata_ref["HOOVER_HVA4"],
      Committed: false,
      H: 5.9892001,
      Rating: 130 },

    { UnitMetadata: metadata_ref["RIMROCKW_D"],
      Committed: true,
      H: 0.01,
      Rating: 20.5 },

}
