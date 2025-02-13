package shorter

import (
	"github.com/spaolacci/murmur3"
)

const abc = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
const abcSize = uint64(len(abc))

func UrlShorter(url string) string {
	numShort := murmur3.Sum64([]byte(url))

	var encoded string
	for numShort > 0 {
		encoded = string(abc[numShort%abcSize]) + encoded
		numShort /= abcSize
	}
	for len(encoded) <= 10 {
		encoded = string(abc[0]) + encoded
	}
	return encoded[:10]
}
