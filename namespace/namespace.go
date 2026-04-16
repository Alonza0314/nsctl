package namespace

import (
	"fmt"

	"github.com/vishvananda/netns"
)

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

func GetNsFd(ns string) (netns.NsHandle, func(), error) {
	found, err := GetNs(ns)
	if err != nil {
		return -1, nil, err
	}
	if !found {
		return -1, nil, fmt.Errorf("namespace %s is not found", ns)
	}

	nsFd, err := netns.GetFromName(GetNsName(ns))
	if err != nil {
		return -1, nil, fmt.Errorf("get namespace %s failed: %v", ns, err)
	}

	closeFunc := func() {
		if err := nsFd.Close(); err != nil {
			fmt.Printf("failed to close %s file descriptor\n", ns)
		}
	}

	return nsFd, closeFunc, nil
}
