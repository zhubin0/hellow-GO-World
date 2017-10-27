package rsa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelperEncryptByPublicKey(t *testing.T) {
	for i, p := range plainTexts {
		helper := NewHelper()
		enc, err := helper.EncryptByPublicKey(publicKey, p)
		assert.Nil(t, err)
		dec, err := helper.DecryptByPrivateKey(privateKey, enc)
		assert.Nil(t, err)
		assert.Equal(t, plainTexts[i], string(dec))
	}
}
