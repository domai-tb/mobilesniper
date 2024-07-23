package enum

import (
	"fmt"
	"strings"
)

type Target struct {
	Host string
	IP   string
	Port int
}

func (t Target) String() string {
	return fmt.Sprintf("%s - %s:%d", t.Host, t.IP, t.Port)
}

func PrintTargets(t *[]Target) string {
	targets := make([]string, len(*t))
	for i, target := range *t {
		targets[i] = fmt.Sprint(target)
	}

	fmt.Printf("%s", strings.Join(targets, "\n"))
	return strings.Join(targets, "\n")
}
