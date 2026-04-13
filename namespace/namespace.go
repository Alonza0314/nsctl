package namespace

import "fmt"

const (
	NS_PREFIX = "nsctl-"
)

func GetNsName(ns string) string {
	return fmt.Sprintf("%s%s", NS_PREFIX, ns)
}
