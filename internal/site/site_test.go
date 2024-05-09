package site

import "testing"

func TestRegularExpressions(t *testing.T) {
	// problem:
	// err: Get "https://en.wikipedia/w/index.php?title=HTTP_404&printable=yes": dial tcp: lookup en.wikipedia: no such host
	// err: Get "https://en.wikipedia/wiki/Wikipedia:Protection_policy#semi": dial tcp: lookup en.wikipedia: no such host
	// err: Get "https://en.wikipedia/wiki/404_Not_Found_(Mr._Robot)": dial tcp: lookup en.wikipedia: no such host
	// err: Get "https://en.wikipedia/wiki/HTTP": dial tcp: lookup en.wikipedia: no such host
}
