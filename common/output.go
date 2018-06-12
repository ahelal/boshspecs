package common

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/savaki/jq"
)

//BoshJSONOutput represents a bosh outpit with --env json
type BoshJSONOutput struct {
	Tables []struct {
		Content string `json:"Content"`
		Header  struct {
			Az           string `json:"az"`
			Instance     string `json:"instance"`
			Ips          string `json:"ips"`
			ProcessState string `json:"process_state"`
		} `json:"Header"`
		Rows []struct {
			Az           string `json:"az"`
			Instance     string `json:"instance"`
			Ips          string `json:"ips"`
			ProcessState string `json:"process_state"`
		} `json:"Rows"`
		Notes interface{} `json:"Notes"`
	} `json:"Tables"`
	Blocks interface{} `json:"Blocks"`
	Lines  []string    `json:"Lines"`
}

//ParseJSON query with JQ compatible language a JSON string
func ParseJSON(jsonString string, query string) (string, error) {
	op, err := jq.Parse(query)
	if err != nil {
		return "", err
	}
	jsonData := []byte(jsonString)
	value, err := op.Apply(jsonData)
	if err != nil {
		return "", err
	}

	return string(value), nil
}

//InfoPrint print an info message
func InfoPrint(info string) {
	fmt.Println("I>", info)
}

//Info print info with a prefix
func Info(prefix string, info string) {
	fmt.Println(prefix, info)
}

// MatrixPrint Prints a matrix
func MatrixPrint(headers []string, items [][]string) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, strings.Join(headers, "\t"))
	for _, item := range items {
		line := strings.Join(item, "\t")
		fmt.Fprintf(w, line)
		fmt.Fprintf(w, "\n")
	}
	fmt.Fprintln(w)
	w.Flush()
}
