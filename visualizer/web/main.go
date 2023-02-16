package web

import (
    "encoding/json"
    "errors"
    "net/http"
    "fmt"
    "time"

    "github.com/G-PST/inertia/internal"
)

const NoNewData = 204
const BadRequest = 400
const ServerError = 500
const ISO8601 = "2006-01-02T15:04:05"

type WebVisualizer struct {

    States []internal.SystemState

}

func New(bind string) *WebVisualizer {

    wv := &WebVisualizer {}

    // http.HandleFunc("/", _) // TODO: Serve dashboard
    // http.HandleFunc("/static", _) // TODO: Serve static assets
    http.HandleFunc("/inertia", serveInertiaData(wv))

    go http.ListenAndServe(bind, nil)

    return wv
}

func (wv *WebVisualizer) Update(state internal.SystemState) {
    wv.States = append(wv.States, state) 
}

func serveInertiaData(wv *WebVisualizer) http.HandlerFunc {

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
            return
        }

        response, err := jsonify(state.Inertia())
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

    t, err := time.Parse(ISO8601, timestamp)
    if err != nil {
        return time.Time {}, err
    }

    return t, nil

}

func getNewer(states []internal.SystemState, latest time.Time) (internal.SystemState, error) {

    for _, state := range states {

        if state.Time.After(latest) {
            return state, nil
        }

    }

    return internal.SystemState {}, errors.New("No newer states avaialable")

}

func jsonify(report internal.InertiaReport) ([]byte, error) {

    response := map[string]any {
        "time": report.Time.Format(ISO8601),
        "total": report.Total,
        "requirement": report.Requirement,
        "inertia": report.Categories,
    }

    return json.Marshal(response)

}
