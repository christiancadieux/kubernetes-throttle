package utils

import (
	"os"
	"strings"
)

func MyNodeName() string {
	nodename := os.Getenv("MY_NODE_NAME")
	if nodename != "" {
		return nodename
	}
	hostname := os.Getenv("HOSTNAME")
	words := strings.Split(hostname, "--")
	rc := ""
	if len(words) > 1 {
		rc = strings.ReplaceAll(words[1], "-", ".")
	}
	return rc
}
