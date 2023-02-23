package localfile

import (
    "github.com/G-PST/inertia/internal"

    "reflect"
    "testing"
)


func TestMetdataParsing(t *testing.T) {

    metadir := "test/toy/metadata"

    system, err := loadSystem(metadir)
    if err != nil {
        t.Fatalf("Error loading system metadata: %v", err)
    } else if system != system_ref {
        t.Fatalf("Loaded system metadata did not match expected value")
    }

    regions, err := loadRegions(metadir)

    if err != nil {
        t.Fatalf("Error loading regions metadata: %v", err)
    } else if !reflect.DeepEqual(regions, regions_ref) {
        t.Fatalf("Loaded regions metadata did not match expected value")
    }

    categories, err := loadCategories(metadir)

    if err != nil {
        t.Fatalf("Error loading categories metadata: %v", err)
    } else if !reflect.DeepEqual(categories, categories_ref) {
        t.Fatalf("Loaded categories metadata did not match expected value")
    }

    metadata, err := loadUnits(metadir, regions, categories)

    if err != nil {
        t.Fatalf("Error parsing metadata metadata: %v", err)
    } else if !reflect.DeepEqual(metadata, metadata_ref) {
        t.Fatalf("Loaded units metadata did not match expected value")
    }

}

var system_ref  = SystemMetadata { requirement: 1500 }

var regions_ref = map[string]*internal.Region {
    "AESO":   &internal.Region { "AESO" },
    "WAPALC": &internal.Region { "WAPALC" },
    "WWA":    &internal.Region { "WWA" },
}

var categories_ref = map[string]*internal.UnitCategory {
    "NGAS": &internal.UnitCategory { "NGAS", "#52216B", 1 },
    "HYDR": &internal.UnitCategory { "HYDR", "#187F94", 2 },
    "WIND": &internal.UnitCategory { "WIND", "#00B6EF", 3 },
}

var metadata_ref = map[string]internal.UnitMetadata{

    "29EDD18_GEN": {
        Name: "29EDD18_GEN",
        Category: categories_ref["NGAS"],
        Region: regions_ref["AESO"],
    },

    "29EDD2_G3": {
        Name: "29EDD2_G3",
        Category: categories_ref["NGAS"],
        Region: regions_ref["AESO"],
    },

    "29EDD3_G5": {
        Name: "29EDD3_G5",
        Category: categories_ref["NGAS"],
        Region: regions_ref["AESO"],
    },

    "HOOVER_HVA4": {
        Name: "HOOVER_HVA4",
        Category: categories_ref["HYDR"],
        Region: regions_ref["WAPALC"],
    },

    "RIMROCKW_D": {
        Name: "RIMROCKW_D",
        Category: categories_ref["WIND"],
        Region: regions_ref["WWA"],
    },

}
