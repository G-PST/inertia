package localfile

import (
    "github.com/G-PST/inertia/internal"

    "os"
    "reflect"
    "testing"
)


func TestMetdataParsing(t *testing.T) {

    metadatafile, err := os.Open("test/sample_metadata.csv")
    if err != nil {
        t.Fatalf("Error loading metadata sample: %v", err)
    }

    metadata, err := parseMetadata(metadatafile)
    if err != nil {
        t.Fatalf("Error parsing metadata sample: %v", err)
    }

    if !reflect.DeepEqual(metadata, metadata_ref) {
        t.Fatalf("Parsed metadata sample did not match expected value")
    }

}

var metadata_ref = map[string]internal.UnitMetadata{

    "29EDD18_GEN": {
        Name: "29EDD18_GEN",
        Category: "NGAS",
        Region: "AESO" },

    "29EDD2_G3": {
        Name: "29EDD2_G3",
        Category: "NGAS",
        Region: "AESO" },

    "29EDD3_G5": {
        Name: "29EDD3_G5",
        Category: "NGAS",
        Region: "AESO" },

    "HOOVER_HVA4": {
        Name: "HOOVER_HVA4",
        Category: "HYDR",
        Region: "WAPALC" },

    "RIMROCKW_D": {
        Name: "RIMROCKW_D",
        Category: "WIND",
        Region: "WWA" },

}
