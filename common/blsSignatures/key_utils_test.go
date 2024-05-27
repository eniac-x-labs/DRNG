package blsSignatures

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_decode(t *testing.T) {
	pubKeyStr := "YBOQS5J342N9BwCWcvtme7v+xLE66PIBOYBsHJpom8PldLqFa7GoIWbaHBzQIXwl/Rc8hp5sz4DHx0A4tlDrTeyhP1AxC2jxLtxdcQAcJuElM096wUi3SHRlBkWHeLe0GxbClkVLAHFw3w2igcfFtjF8SVG0sAk3obur/4ogWFC/oVV0GTTeS4L/edNrBObIBgHGeE5PDOXTnQrStGV9vi3rbLEfECJk8W+A4HsWpCjwV7Xinm7VlS0yJUElA/rvkBY7KwIA2yvclH9xYcAlfGhbLwReESdaeLq58jVqQiuIQzlyZ8xbUyO1NfQtAb/sUAMgfCp/ARlqVNWsEBP/7IlPIDGJnoLhkIqryFdIJ05lDfY1J+pZB7xWjh0+/qE3FQ=="
	key, err := DecodeBase64BLSPublicKey([]byte(pubKeyStr))
	assert.NoError(t, err)
	assert.NotNil(t, key.key)
}
