package base

import (
	"strings"
	"time"
)

func TimeNow() string {
	str := strings.Split(time.Now().String(), " ")[0] + " " + strings.Split(time.Now().String(), " ")[1]
	for len(str) < 30 {
		str += " "
	}
	return str
}
