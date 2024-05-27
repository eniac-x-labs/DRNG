package common

import (
	"crypto/sha256"
	"testing"
)

func Test_Uints2Hash(t *testing.T) {
	nums := []uint64{9671254531385894685, 7757817344683588787}
	t.Log(Uint64sToStrings(sha256.New(), nums))
	// [0x2dfbca7ff946a79f2838e40d99eb2903b7047238a6adb4f4042085e1daf2fccf 0xf0156b1e1871e80aa64069aa307fbe676167991bd5df8199a772aef3fe9a0ee3]
	t.Log(Uint64sToStrings(sha256.New(), nums))
}
