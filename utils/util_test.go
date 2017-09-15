package utils

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnap(t *testing.T) {
	cases := []struct {
		input string
		o1    string
		o2    string
	}{
		{"sha256:1234", "12", "1234"},
		{"1234", "12", "1234"},
		{"sha256:1", "", ""},
	}

	for _, c := range cases {
		o1, o2 := Snap(c.input)
		assert.Equal(t, c.o1, o1)
		assert.Equal(t, c.o2, o2)
	}
}

func TestGetDigest(t *testing.T) {
	expected := "e637dbddc5be31a623306e9e21e4d7c77878ba425445649d2a6804755976f29d"
	data, _ := ioutil.ReadFile("testdata/" + expected)
	assert.Equal(t, expected, GetDigest("sha256", data))

	assert.Equal(t, "", GetDigest("sha512", data))
}
