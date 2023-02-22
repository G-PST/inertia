package text

import (
    "io"
    "fmt"
    "os"
    "time"

    "github.com/G-PST/inertia/internal"
)

// Visualizer outputting inertia data in text format
type TextVisualizer struct {
    outfile io.StringWriter
}

// Creates a new TextVisualizer that reports results via standard output
func New() TextVisualizer {
    return TextVisualizer { os.Stdout }
}

func (tv TextVisualizer) Init(state internal.SystemMetadata) error {
    return nil
}

func (tv TextVisualizer) Update(state internal.SystemState) {

    summary := state.Inertia()
    timestamp := state.Time.Format(time.RFC3339)

    text := fmt.Sprintf("%v: %v MWs\n", timestamp, summary.Total)
    tv.outfile.WriteString(text)

    for category, inertia := range summary.Categories {
        text = fmt.Sprintf("\t%v\t%v MWs\n", category, inertia)
        tv.outfile.WriteString(text)
    }

}
