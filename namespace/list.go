package namespace

import (
	"fmt"
	"os"
	"strings"
)

func List() error {
	nsFile, err := os.ReadDir("/var/run/netns")
	if err != nil {
		return fmt.Errorf("failed to list ns: %v", err)
	}

	nsList := make([]string, 0)
	for _, ns := range nsFile {
		if strings.HasPrefix(ns.Name(), NS_PREFIX) {
			nsList = append(nsList, ns.Name()[len(NS_PREFIX):])
		}
	}

	if len(nsList) != 0 {
		var strBuilder strings.Builder
		strBuilder.WriteString("Namespace List:\n")
		for _, ns := range nsList {
			fmt.Fprintf(&strBuilder, "- %s\n", ns)
		}

		fmt.Println(strBuilder.String())
	}

	fmt.Printf("Total: %d namespaces\n", len(nsList))
	return nil
}
