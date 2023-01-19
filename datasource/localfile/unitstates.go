package localfile

import (
    "github.com/G-PST/inertia/internal"

    "bufio"
    "errors"
    "io"
    "strconv"
    "strings"
)

/*
Load in files from disk, where the filename structure is:
    UNInertia_yyyy-mm-dd-HH-MM.csv

*/
func parseUnitStates(f io.Reader, metadata map[string]internal.UnitMetadata) ([]internal.UnitState, error) {

    scanner := bufio.NewScanner(f)
    var line string

    cols := NewUnitStateColumnIndices()
    colnum := 0

    rep := strings.NewReplacer(",", "", "\\", "", "\"", "")

    for { 

        if !scanner.Scan() {
            return nil, errors.New("Ran out of column definition lines to parse")
        }

        line = scanner.Text()

        firstchar := line[0:1]
        if firstchar != "#" && firstchar != "," {
            break
        }

        colname := strings.TrimSpace(rep.Replace(line))

        switch colname {
        case "ID_ST":
            cols.station_name = colnum
        case "ID_UN":
            cols.unit_name = colnum
        case "OPEN_UN":
            cols.open = colnum
        case "REMOVE_UN":
            cols.remove = colnum
        case "H_UN":
            cols.h = colnum
        case "MVARATE_UN":
            cols.rating = colnum
        }

        colnum += 1

    }

    if cols.HasUnpopulated() {
        return nil, errors.New("Missing required columns")
    }

    units := []internal.UnitState {}

    for {

        fields := strings.Split(line, ",")

        station_name := strings.Trim(fields[cols.station_name], "'\" ")
        unit_name := strings.Trim(fields[cols.unit_name], "'\" ")
        name := station_name + "_" + unit_name

        committed := fields[cols.open] == "F" && fields[cols.remove] == "F"

        h, err := strconv.ParseFloat(fields[cols.h], 64)
        if err != nil {
            return nil, errors.New("Error parsing H")
        }

        rating, err := strconv.ParseFloat(fields[cols.rating], 64)
        if err != nil {
            return nil, errors.New("Error parsing Rating")
        }

        unitmeta, ok := metadata[name]
        if !ok {
            return nil, errors.New("Encountered unit with no metadata (" + name + ")")
        }

        units = append(units, internal.UnitState {
            UnitMetadata: unitmeta,
            Committed: committed,
            H: h,
            Rating: rating,
        })

        if !scanner.Scan() { break }

        line = scanner.Text()

    }

    return units, nil

}

type UnitStateColumnIndices struct {
    station_name int
    unit_name int
    open int
    remove int
    h int
    rating int
}

func NewUnitStateColumnIndices() UnitStateColumnIndices {
    return UnitStateColumnIndices { -1, -1, -1, -1, -1, -1 }
}

func (ci UnitStateColumnIndices) HasUnpopulated() bool {
    return ci.station_name < 0 || ci.unit_name < 0 ||
           ci.open < 0 || ci.remove < 0 || ci.h < 0 || ci.rating < 0
}
