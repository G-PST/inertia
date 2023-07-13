package web

import (
    "encoding/json"
    "embed"
    "errors"
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

type WebVisualizer struct {

    bind string

    Metadata inertia.SystemMetadata
    States []inertia.Snapshot

}

func New(bind string) *WebVisualizer {

    wv := &WebVisualizer { bind: bind }
    return wv

}

func (wv *WebVisualizer) Init(meta inertia.SystemMetadata) error {

    appdir, err := fs.Sub(assets, "app")
    if err != nil { return err }

    wv.Metadata = meta
    http.Handle("/", http.FileServer(http.FS(appdir)))
    http.HandleFunc("/metadata", serveMetadata(wv))
    http.HandleFunc("/inertia", serveInertiaData(wv))

    go http.ListenAndServe(wv.bind, nil)

    return nil

}

func (wv *WebVisualizer) Update(state inertia.Snapshot) {
    wv.States = append(wv.States, state) 
}

func serveMetadata(wv *WebVisualizer) http.HandlerFunc {

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

func serveInertiaData(wv *WebVisualizer) http.HandlerFunc {

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

        state, err := getNewer(wv.States, latest)
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

func getNewer(states []inertia.Snapshot, latest time.Time) (inertia.Snapshot, error) {

    for _, state := range states {

        if state.Time.After(latest) {
            return state, nil
        }

    }

    return inertia.Snapshot {}, errors.New("No newer states avaialable")

}

// TODO: Just define appropriate methods in inertia/internal?
func jsonify_meta(meta inertia.SystemMetadata) ([]byte, error) {

    regions := map[string]inertia.Region {}
    categories := map[string]inertia.UnitCategory {}

    for _, region := range meta.Regions {
        regions[region.Name] = region
    }

    for _, category := range meta.Categories {
        categories[category.Name] = category
    }

    response := map[string]any {
        "regions": regions,
        "categories": categories,
    }

    return json.Marshal(response)

}

func jsonify(report inertia.Snapshot) ([]byte, error) {

    response := map[string]any {
        "time": report.Time.UnixMilli(),
        "total": report.Total,
        "requirement": report.Requirement,
        "inertia": report.Categories,
    }

    return json.Marshal(response)

}
