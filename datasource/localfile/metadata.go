package localfile

import (
    "github.com/G-PST/inertia/internal"

    "bufio"
    "errors"
    "io"
    "strings"
)

/*
Load in unit and category metadata
*/
func parseMetadata(f io.Reader) (map[string]internal.UnitMetadata, error) {

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

        metadata[name] = internal.UnitMetadata {
            Name: name,
            Category: fields[cols.category],
            Region: fields[cols.region],
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
