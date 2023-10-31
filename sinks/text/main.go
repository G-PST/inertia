package text

import (
    "io"
    "fmt"
    "os"
    "time"

    "github.com/G-PST/inertia"
)

// DataSink outputting inertia data in text format
type TextDataSink struct {
    outfile io.StringWriter
}

// Creates a new TextDataSink that reports results via standard output
func New() TextDataSink {
    return TextDataSink { os.Stdout }
}

func (tv TextDataSink) Init(state inertia.SystemMetadata) error {
    return nil
}

func (tv TextDataSink) Update(snapshot inertia.Snapshot) {

    timestamp := snapshot.Time.Format(time.RFC3339)

    text := fmt.Sprintf("%v: %v MWs\n", timestamp, snapshot.Total.TotalInertia)
    tv.outfile.WriteString(text)

    for category, inertia := range snapshot.Categories {
        text = fmt.Sprintf("\t%v\t%v MWs\n", category, inertia.TotalInertia)
        tv.outfile.WriteString(text)
    }

}
