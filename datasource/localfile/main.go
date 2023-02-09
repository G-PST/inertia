package localfile

import (
    "github.com/G-PST/inertia/internal"

    "errors"
    "io"
    "os"
    "path/filepath"
    "time"
    "fmt"
)

type LocalFile struct {
    dir string
    unitfileFormat string
    lastTime time.Time
    metadata map[string]internal.UnitMetadata
}

func New(dir string, unitfileformat string, metafile string) (*LocalFile, error) {
    
    metadatafile, err := os.Open(metafile)
    if err != nil { return nil, err }

    metadata, err := parseMetadata(metadatafile)
    if err != nil {
        return nil, err
    }

    lf := &LocalFile {
        dir: dir,
        unitfileFormat: unitfileformat,
        metadata: metadata,
    }

    return lf, nil

}

func (lf *LocalFile) Query() (internal.SystemState, error) {

    // Extract filename and timestamp
    file, filetime, err := lf.NextFile()
    if err != nil {
        return internal.SystemState {}, err
    }
    fmt.Println(filetime)

    units, err := parseUnitStates(file, lf.metadata)
    if err != nil {
        return internal.SystemState {}, err
    }

    lf.lastTime = filetime

    return internal.SystemState {filetime, units}, nil

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
