package output

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

func JSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func Table(headers []string, rows [][]string) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	upperHeaders := make([]string, len(headers))
	for i, header := range headers {
		upperHeaders[i] = strings.ToUpper(header)
	}
	fmt.Fprintln(writer, strings.Join(upperHeaders, "\t"))
	for _, row := range rows {
		clean := make([]string, len(row))
		for i, cell := range row {
			clean[i] = strings.NewReplacer("\t", " ", "\r", " ", "\n", " ").Replace(cell)
		}
		fmt.Fprintln(writer, strings.Join(clean, "\t"))
	}
	_ = writer.Flush()
}

func KeyValue(title string, values map[string]string) {
	if title != "" {
		fmt.Println(title)
	}
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s: %s\n", k, values[k])
	}
}
