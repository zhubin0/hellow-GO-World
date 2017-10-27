package mail

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSend(t *testing.T) {
	scm := sendCloudMail{
		From:     "xueguo@datamesh.com",
		Fromname: "xueguo",
		To:       "1203897532@qq.com",
		Subject:  "hahah-0",
		Html:     "<html><body>this is a html page</body></html>",
	}

	err := scm.send("http://192.168.2.36:8094/sendmail")
	assert.Nil(t, err)

}
