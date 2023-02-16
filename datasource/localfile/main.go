package localfile

import (
    "github.com/G-PST/inertia/internal"

    "errors"
    "io"
    "os"
    "path/filepath"
    "time"
)

type LocalFile struct {

    dir string
    unitfileFormat string
    lastTime time.Time

    system SystemMetadata
    regions map[string]*internal.Region
    categories map[string]*internal.UnitCategory
    units map[string]internal.UnitMetadata

}

func New(statedir string, unitfileformat string, metadir string) (*LocalFile, error) {

    system, err := loadSystem(metadir)
    if err != nil { return nil, err }

    regions, err := loadRegions(metadir)
    if err != nil { return nil, err }

    categories, err := loadCategories(metadir)
    if err != nil { return nil, err }

    units, err := loadUnits(metadir, regions, categories)
    if err != nil { return nil, err }

    lf := &LocalFile {
        dir: statedir,
        unitfileFormat: unitfileformat,
        system: system,
        regions: regions,
        categories: categories,
        units: units,
    }

    return lf, nil

}

func (lf *LocalFile) Metadata() internal.SystemMetadata {

    regions := make([]internal.Region, 0, len(lf.regions))

    for _, region := range lf.regions {
        regions = append(regions, *region)
    }

    categories := make([]internal.UnitCategory, 0, len(lf.regions))

    for _, category := range lf.categories {
        categories = append(categories, *category)
    }

    return internal.SystemMetadata { regions, categories }

}

func (lf *LocalFile) Query() (internal.SystemState, error) {

    file, filetime, err := lf.NextFile()
    if err != nil {
        return internal.SystemState {}, err
    }

    units, err := parseUnitStates(file, lf.units)
    if err != nil {
        return internal.SystemState {}, err
    }

    lf.lastTime = filetime

    return internal.SystemState {filetime, lf.system.requirement, units}, nil

}

func (lf LocalFile) NextFile() (io.Reader, time.Time, error) {

    var nextfile string
    var nextfiletime time.Time

    files, err := os.ReadDir(lf.dir)
    if err != nil {
        return nil, time.Time{}, err
    }

    for _, file := range files {

        filename := file.Name()

        filetime, err := time.ParseInLocation(
            lf.unitfileFormat, filename, time.Local)
        if err != nil {
            continue
        }

        later_than_last := lf.lastTime.Before(filetime)
        before_next := filetime.Before(nextfiletime)

        if later_than_last && (before_next || nextfile == "") {
            nextfile = filename
            nextfiletime = filetime
        }

    }

    if nextfile == "" {
        return nil, time.Time{}, errors.New("No new SystemStates available")
    }

    f, err := os.Open(filepath.Join(lf.dir, nextfile))
    if err != nil {
        return nil, nextfiletime, err 
    }

    return f, nextfiletime, nil

}
