package localfile

import (
    "github.com/G-PST/inertia/internal"

    "bufio"
    "errors"
    "encoding/csv"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
)

type SystemMetadata struct {
    requirement float64
}

// System metadata

func loadSystem(metadir string) (SystemMetadata, error) {

    systemfile, err := os.Open(filepath.Join(metadir, "system.csv"))
    if err != nil { return SystemMetadata {}, err }

    return parseSystem(systemfile)

}

func parseSystem(f io.Reader) (SystemMetadata, error) {

    systemfile := csv.NewReader(f)
    systemfile.FieldsPerRecord = 1

    fields, err := systemfile.Read()
    if err != nil { return SystemMetadata {}, err }

    requirement, err := strconv.ParseFloat(fields[0], 64)
    if err != nil { return SystemMetadata {}, err }

    return SystemMetadata { requirement }, nil

}

// Region metadata

func loadRegions(metadir string) (map[string]*internal.Region, error) {

    regionsfile, err := os.Open(filepath.Join(metadir, "regions.csv"))
    if err != nil { return nil, err }

    return parseRegions(regionsfile)

}

func parseRegions(f io.Reader) (map[string]*internal.Region, error) {

    regionfile := csv.NewReader(f)
    regionfile.FieldsPerRecord = 1

    regions := map[string]*internal.Region {}

    for {

        fields, err := regionfile.Read()

        if err == io.EOF {
            break
        } else if err != nil {
            return nil, err
        }

        regionname := fields[0]

        if _, ok := regions[regionname]; ok {
            return nil, errors.New("Duplicate region name")
        }

        regions[regionname] = &internal.Region { regionname }

    }

    return regions, nil

}

// Category metadata

func loadCategories(metadir string) (map[string]*internal.UnitCategory, error) {

    categoriesfile, err := os.Open(filepath.Join(metadir, "categories.csv"))
    if err != nil { return nil, err }

    return parseCategories(categoriesfile)

}

func parseCategories(f io.Reader) (map[string]*internal.UnitCategory, error) {

    regionfile := csv.NewReader(f)
    regionfile.FieldsPerRecord = 3

    categories := map[string]*internal.UnitCategory {}

    for {

        fields, err := regionfile.Read()

        if err == io.EOF {
            break
        } else if err != nil {
            return nil, err
        }

        categoryname := fields[0]
        color := fields[1]

        order,err := strconv.Atoi(fields[2])
        if err != nil { return nil, err }

        if _, ok := categories[categoryname]; ok {
            return nil, errors.New("Duplicate category name")
        }

        categories[categoryname] = &internal.UnitCategory {
            categoryname, color, order,
        }

    }

    return categories, nil

}

// Unit metadata

func loadUnits(
    metadir string,
    regions map[string]*internal.Region,
    categories map[string]*internal.UnitCategory,
) (map[string]internal.UnitMetadata, error) {

    metadatafile, err := os.Open(filepath.Join(metadir, "units.csv"))
    if err != nil { return nil, err }

    return parseUnits(metadatafile, regions, categories)

}

func parseUnits(
    f io.Reader,
    regions map[string]*internal.Region,
    categories map[string]*internal.UnitCategory,
) (map[string]internal.UnitMetadata, error) {

    scanner := bufio.NewScanner(f)
    cols := NewMetadataColumnIndices()

    if !scanner.Scan() {
        return nil, errors.New("File is empty")
    }
    colnames := strings.Split(scanner.Text(), ",")
    n_cols := len(colnames)

    for colnum, colname := range colnames {
        switch colname {
        case "ID_ST":
            cols.station_name = colnum
        case "ID_UN":
            cols.unit_name = colnum
        case "GTYPE":
            cols.category = colnum
        case "ID_CO":
            cols.region = colnum
        }
    }

    if cols.HasUnpopulated() {
        return nil, errors.New("Missing required columns")
    }

    metadata := make(map[string]internal.UnitMetadata)

    for scanner.Scan() {

        line := scanner.Text()
        fields := strings.Split(line, ",")
        if len(fields) != n_cols {
            return nil, errors.New("Malformed table structure")
        }

        name := fields[cols.station_name] + "_" + fields[cols.unit_name]
        category := fields[cols.category]
        region := fields[cols.region]

        metadata[name] = internal.UnitMetadata {
            Name: name,
            Category: categories[category],
            Region: regions[region],
        }

    }

    return metadata, nil

}

type MetadataColumnIndices struct {
    station_name int
    unit_name int
    region int
    category int
}

func NewMetadataColumnIndices() MetadataColumnIndices {
    return MetadataColumnIndices { -1, -1, -1 , -1 }
}

func (ci MetadataColumnIndices) HasUnpopulated() bool {
    return ci.region < 0 || ci.category < 0
}
