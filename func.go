package putils

import (
	"bytes"
	"strconv"

	"github.com/plivepatch/go-utils/md5"
)

func PK(endpoint, metric string, tags map[string]string) string {
	ret := bufferPool.Get().(*bytes.Buffer)
	ret.Reset()
	defer bufferPool.Put(ret)

	if tags == nil || len(tags) == 0 {
		ret.WriteString(endpoint)
		ret.WriteString("/")
		ret.WriteString(metric)

		return ret.String()
	}
	ret.WriteString(endpoint)
	ret.WriteString("/")
	ret.WriteString(metric)
	ret.WriteString("/")
	ret.WriteString(SortedTags(tags))
	return ret.String()
}

func PK2(endpoint, counter string) string {
	ret := bufferPool.Get().(*bytes.Buffer)
	ret.Reset()
	defer bufferPool.Put(ret)

	ret.WriteString(endpoint)
	ret.WriteString("/")
	ret.WriteString(counter)

	return ret.String()
}

func UUID(endpoint, metric string, tags map[string]string, dstype string, step int) string {
	ret := bufferPool.Get().(*bytes.Buffer)
	ret.Reset()
	defer bufferPool.Put(ret)

	if tags == nil || len(tags) == 0 {
		ret.WriteString(endpoint)
		ret.WriteString("/")
		ret.WriteString(metric)
		ret.WriteString("/")
		ret.WriteString(dstype)
		ret.WriteString("/")
		ret.WriteString(strconv.Itoa(step))

		return ret.String()
	}
	ret.WriteString(endpoint)
	ret.WriteString("/")
	ret.WriteString(metric)
	ret.WriteString("/")
	ret.WriteString(SortedTags(tags))
	ret.WriteString("/")
	ret.WriteString(dstype)
	ret.WriteString("/")
	ret.WriteString(strconv.Itoa(step))

	return ret.String()
}

func Checksum(endpoint string, metric string, tags map[string]string) string {
	pk := PK(endpoint, metric, tags)
	return md5.Md5(pk)
}

func ChecksumOfUUID(endpoint, metric string, tags map[string]string, dstype string, step int64) string {
	return md5.Md5(UUID(endpoint, metric, tags, dstype, int(step)))
}
