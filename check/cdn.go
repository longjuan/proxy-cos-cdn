package check

import (
	"sync"

	"proxy-cos-cdn/tencentyun"
	"proxy-cos-cdn/types"
)

func CDNReady(infos []*types.BucketInfo, secretID string, secretKey string) {
	tencentyun.InitCDNClient(secretID, secretKey)
	var wg sync.WaitGroup
	for _, info := range infos {
		wg.Add(1)
		go func(info *types.BucketInfo) {
			status := tencentyun.GetDomainStatus([]*string{&info.Domain, &info.WildcardDomain})
			info.CDNReady = status
			wg.Done()
		}(info)
	}
	wg.Wait()
}
