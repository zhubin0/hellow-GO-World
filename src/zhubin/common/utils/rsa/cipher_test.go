package rsa

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

var privateKey = `MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBAISX9Vk8y8Wu+CVBdA2YYWNi8cJhhpFT0QTFXSrgBN5YW3xpoRKW/1z5bIEHyfg0bt9oN/5wps1i4oxu48hsWeLWj+sCv4kQv2CknQHZiFLvtzrw5WjlWfKegKsY8YkbtZ6uMz4tcMdWoeAPvFhrMRUmvCo05XVP4aILxzUCd1G1AgMBAAECgYB14H5PWjwyP4392QWqfHjAGZuiWn9+vYwJ+MgOMOBDJzwWC/YVh8X4SwoKX/lPPpX+6TE2c8Hmv+12ObMpYCI7vwlo5eVKZANE91KW9NiFwDuN3w0z3NFVTO86+/aRre9k+Gba51oMap5rR9UpswDBtzoYV8YnWKjrAb+h0LC/lQJBAL7KzagzqeGzbiXnZLQ8lR7sTOFLylWoEtS+1MSHCxymSa6mObFq0KvshT/cZGdIh/4OIW/x1EOpe9u2fC2zu5cCQQCx6R5Jyq5QVyud9whu6tzhnrAGj8KJyMHFvSUaKN3uCsfHuj4Fcjdgdidc84MpkzXkcIEgdMpuFTOpK7Tn7XaTAkBEJgp5hyKqFL5GWbWVz4HwTrVTUBAQsn0vco5rOFVWwWrWMFexMJciodQisGVIoxa4P3HgG4AXPwWXwEHwzR83AkBmNPiTh/7QZOPH4i1UG1U9wL57ZodqRI0dnmX8O1IT+NmA4nvTASTTI83FVpgZgFrLm95y2OWajE+bdmJ9gyxFAkAwkBdmAJMSPMb4CoefFv4XMI1aZIuxgUMN39G0/fvq8b4ZEl8z2rxhMBv6X9WlWyMoHdVIDAmgJuIoDd4yNNCb`

var publicKey = `MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCEl/VZPMvFrvglQXQNmGFjYvHCYYaRU9EExV0q4ATeWFt8aaESlv9c+WyBB8n4NG7faDf+cKbNYuKMbuPIbFni1o/rAr+JEL9gpJ0B2YhS77c68OVo5VnynoCrGPGJG7WerjM+LXDHVqHgD7xYazEVJrwqNOV1T+GiC8c1AndRtQIDAQAB`

var plainTexts = []string{`test text`, "我们是DataMesh", "2834834*(!@*(#*(#$@&@$(#*)!@#(!@#*(!#&(　｀２３９３２"}

func TestEncryptAndDecrypt(t *testing.T) {
	pubKey, err := base64.StdEncoding.DecodeString(publicKey)
	assert.Nil(t, err)
	priKey, err := base64.StdEncoding.DecodeString(privateKey)
	assert.Nil(t, err)
	key, err := ParsePKCS8Key(pubKey, priKey)
	assert.Nil(t, err)
	for i, p := range plainTexts {
		cipher, err := NewRSA(key)
		assert.Nil(t, err)
		enT, err := cipher.Encrypt([]byte(p))
		assert.Nil(t, err)
		deT, err := cipher.Decrypt(enT) // by private key
		assert.Nil(t, err)
		assert.Equal(t, plainTexts[i], string(deT))
	}
}
