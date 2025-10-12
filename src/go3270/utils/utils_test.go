package utils_test

import (
	"go3270/utils"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var m1 = map[byte]string{
	0x01: "1",
	0x02: "2",
}

var m2 = map[string]byte{
	"1": 0x01,
	"2": 0x02,
}

var m3 = map[byte]string{
	0x01: "1",
	0x02: "2",
	0x03: "2",
}

var m4 = map[string][]byte{
	"1": {0x01},
	"2": {0x02, 0x03},
}

func Test_Invert(t *testing.T) {
	assert.True(t, reflect.DeepEqual(m2, utils.Invert(m1)))
	assert.True(t, reflect.DeepEqual(m4, utils.InvertMulti(m3)))
}

func Test_Ternary(t *testing.T) {
	assert.True(t, utils.Ternary(true, "a", "b") == "a")
	assert.True(t, utils.Ternary(false, "a", "b") == "b")
}
