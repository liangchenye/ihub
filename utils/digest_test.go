// digest.go is from docker/docker project
package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlgAvailable(t *testing.T) {
	assert.Equal(t, true, SHA256.Available())
	assert.Equal(t, false, TarsumV1SHA256.Available())
}

func TestDigestManifest(t *testing.T) {
	data := `
{
   "schemaVersion": 2,
   "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
   "config": {
      "mediaType": "application/vnd.docker.container.image.v1+json",
      "size": 4,
      "digest": "sha256:ccc7a11d65b1b5874b65adb4b2387034582d08d65ac1817ebc5fb9be1baa5f88"
   },
   "layers": [
      {
         "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
         "size": 46057336,
         "digest": "sha256:95f4beaf474661dfe65828d8903f1a907dc7bcfdcf8f502425422f582a2c4135"
      }
	  ]
}
`
	p, err := DigestManifest([]byte(data))
	fmt.Println(p, err)

	_, err = DigestManifest([]byte("invalid"))
	assert.NotNil(t, err)
}
