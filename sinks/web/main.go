package web

import (
    "encoding/json"
    "embed"
    "fmt"
    "io/fs"
    "net/http"
    "strconv"
    "time"

    "github.com/G-PST/inertia"
)

const NoNewData = 204
const BadRequest = 400
const ServerError = 500

//go:embed app
var assets embed.FS

type WebDataSink struct {

    bind string

    Metadata inertia.SystemMetadata
    States SnapshotRing

}

func New(bind string, bufferlength int) *WebDataSink {

    buffer := NewSnapshotRing(bufferlength)

    wv := &WebDataSink { bind: bind, States: buffer }
    return wv

}

func (wv *WebDataSink) Init(meta inertia.SystemMetadata) error {

    appdir, err := fs.Sub(assets, "app")
    if err != nil { return err }

    wv.Metadata = meta
    http.Handle("/", http.FileServer(http.FS(appdir)))
    http.HandleFunc("/metadata", serveMetadata(wv))
    http.HandleFunc("/inertia", serveInertiaData(wv))

    go http.ListenAndServe(wv.bind, nil)

    return nil

}

func (wv *WebDataSink) Update(state inertia.Snapshot) {
    wv.States.Push(state)
}

func serveMetadata(wv *WebDataSink) http.HandlerFunc {

    return func(w http.ResponseWriter, r *http.Request) {

        meta, err := jsonify_meta(wv.Metadata)
        if err != nil {
            w.WriteHeader(ServerError)
            return
        }

        w.Write(meta)
        return

    }

}

func serveInertiaData(wv *WebDataSink) http.HandlerFunc {

    // TODO: Return ALL states newer than latest, not just one
    return func(w http.ResponseWriter, r *http.Request) {

        err := r.ParseForm()
        if err != nil {
            w.WriteHeader(ServerError)
            return
        }

        latest, err := parseTime(r.FormValue("last"))
        if err != nil {
            w.WriteHeader(BadRequest)
            fmt.Fprintln(w, "Invalid last timestamp")
            return
        }

        state, err := wv.States.FirstAfter(latest)
        if err != nil {
            w.WriteHeader(NoNewData)
            fmt.Fprintln(w, "No new data")
            return
        }

        response, err := jsonify(state)
        if err != nil {
            w.WriteHeader(ServerError)
            return
        }

        w.Write(response)
        return

    }

}

func parseTime(timestamp string) (time.Time, error) {

    if timestamp == "" {
        return time.Time {}, nil
    }

    t, err := strconv.Atoi(timestamp)
    if err != nil {
        return time.Time {}, err
    }

    return time.UnixMilli(int64(t) + 1), nil

}

// TODO: Just define appropriate methods in inertia/internal?
func jsonify_meta(meta inertia.SystemMetadata) ([]byte, error) {

    response := map[string]any {
        "regions": meta.Regions,
        "categories": meta.Categories,
    }

    return json.Marshal(response)

}

func jsonify(report inertia.Snapshot) ([]byte, error) {

    response := map[string]any {
        "time": report.Time.UnixMilli(),
        "requirement": report.Requirement,
        "total": report.Total,
        "categories": report.Categories,
        "regions": report.Regions,
    }

    return json.Marshal(response)

}
