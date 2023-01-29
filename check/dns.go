package check

import (
	"context"
	"net"
	"sync"
	"time"

	"proxy-cos-cdn/types"
)

func DNSRecord(bucketInfos []*types.BucketInfo, domainSuffix string) {
	var wg sync.WaitGroup
	for _, bucketInfo := range bucketInfos {
		bucketInfo.Domain = bucketInfo.DisplayName + "." + domainSuffix
		wg.Add(1)
		go func(info *types.BucketInfo) {
			defer wg.Done()
			err := checkDNSRecord(context.Background(), info.Domain)
			if err != nil {
				info.DNSResult = false
			} else {
				info.DNSResult = true
			}
			info.CDNReady = "Unchecked"
			info.WildcardDomain = "*." + domainSuffix
		}(bucketInfo)
	}
	wg.Wait()
}

func checkDNSRecord(ctx context.Context, domain string) error {
	// 设置超时时间为 2s
	_, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := net.LookupHost(domain)
	if err != nil {
		return err
	}
	return nil
}
