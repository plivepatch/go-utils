package env

import (
	"io/ioutil"
	"regexp"
	"strconv"
)

// GetKernelVersion returns the major and minor versions of the Linux
func GetKernelVersion() (int, int) {
	bytes, err := ioutil.ReadFile("/proc/version")
	if err != nil {
		return -1, -1
	}
	matches := regexp.MustCompile("([0-9]+).([0-9]+)").FindSubmatch(bytes)
	if len(matches) < 3 {
		return -1, -2
	}
	major, _ := strconv.Atoi(string(matches[1]))
	minor, _ := strconv.Atoi(string(matches[2]))
	return major, minor
}
