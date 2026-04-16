package namespace

import "fmt"

const (
	NS_PREFIX = "nsctl-"
)

func GetNsName(ns string) string {
	return fmt.Sprintf("%s%s", NS_PREFIX, ns)
}

func GetNs(nsTarget string) (bool, error) {
	list, err := getNsList()
	if err != nil {
		return false, err
	}

	for _, ns := range list {
		if nsTarget == ns {
			return true, nil
		}
	}

	return false, nil
}
