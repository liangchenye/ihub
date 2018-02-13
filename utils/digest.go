package utils

import (
	"bytes"
	"strings"

	"github.com/docker/libtrust"
	"github.com/opencontainers/go-digest"
)

// Payload gets the playload of a docker manifest file
func Payload(data []byte) ([]byte, error) {
	jsig, err := libtrust.ParsePrettySignature(data, "signatures")
	if err != nil {
		return nil, err
	}

	// Resolve the payload in the manifest.
	return jsig.Payload()
}

// DigestManifest gets the real digest content of a docker manifest
func DigestManifest(data []byte) (string, error) {
	p, err := Payload(data)
	if err != nil {
		if !strings.Contains(err.Error(), "missing signature key") {
			return "", err
		}

		p = data
	}

	d, err := digest.FromReader(bytes.NewReader(p))
	if err != nil {
		return "", err
	}

	return string(d), err
}
