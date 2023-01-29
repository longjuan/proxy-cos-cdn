package types

import (
	"github.com/tencentyun/cos-go-sdk-v5"
)

type BucketInfo struct {
	BucketName     string
	DisplayName    string
	Domain         string
	WildcardDomain string
	DNSResult      bool
	COSClient      *cos.Client
	Ak             string
	Sk             string
	CDNReady       string
}
