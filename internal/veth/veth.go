package veth

import "fmt"

const (
	NET_PREFIX = ""
)

func GetVethName(ns1, ns2 string) string {
	return fmt.Sprintf("%s%s-%s", NET_PREFIX, ns1, ns2)
}
