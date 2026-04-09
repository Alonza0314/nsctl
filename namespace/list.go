package namespace

import (
	"fmt"
	"os"
	"strings"
)

func List() error {
	nsList, err := getNsList()
	if err != nil {
		return err
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

func getNsList() ([]string, error) {
	nsFile, err := os.ReadDir("/var/run/netns")
	if err != nil {
		return nil, fmt.Errorf("failed to get list ns: %v", err)
	}

	nsList := make([]string, 0)
	for _, ns := range nsFile {
		if strings.HasPrefix(ns.Name(), NS_PREFIX) {
			nsList = append(nsList, ns.Name()[len(NS_PREFIX):])
		}
	}

	return nsList, nil
}
