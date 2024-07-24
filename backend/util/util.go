package util

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"reflect"
	"time"
)

// ShuffleArray shuffles an array in place.
func ShuffleArray(slice interface{}) {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		panic("Shuffle: not a slice")
	}
	n := rv.Len()
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		tmp := reflect.ValueOf(rv.Index(i).Interface())
		rv.Index(i).Set(rv.Index(j))
		rv.Index(j).Set(tmp)
	}
}

// FormatTime formats a time.Time object as a string in the format "2006-01-02 15:04:05.000Z".
func FormatTime(t time.Time) string {
	utcTime := t.UTC()
	layout := "2006-01-02 15:04:05.000Z"
	return utcTime.Format(layout)
}

// GetMD5Hash returns the MD5 hash of a string.
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
