package conf

type OSS struct {
	Platform        string
	Host            string
	Port            string
	AccessKeyID     string // AccountName for Azure blob
	SecretAccessKey string // AccountKey for Azure blob
	ApiVersion      string // For Azure only
	Bucket          string
	Location        string
	UseSSL          bool
}
