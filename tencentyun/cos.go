package tencentyun

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/tencentyun/cos-go-sdk-v5"
	"proxy-cos-cdn/types"
)

func InitCosClients(bucketRegion, secretID, secretKey string) (infos []*types.BucketInfo) {
	noneBucketClient := cos.NewClient(&cos.BaseURL{
		BucketURL: &url.URL{
			Scheme: "https",
			Host:   "cos." + bucketRegion + ".myqcloud.com",
		},
	}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})
	bucketList, _, err := noneBucketClient.Service.Get(context.Background())
	if err != nil {
		panic(err)
	}
	infos = make([]*types.BucketInfo, len(bucketList.Buckets))
	for i, bucket := range bucketList.Buckets {
		displayName := strings.Split(bucket.Name, "-")[0]
		bucketClient := cos.NewClient(&cos.BaseURL{
			BucketURL: &url.URL{
				Scheme: "https",
				Host:   bucket.Name + ".cos." + bucketRegion + ".myqcloud.com",
			},
		}, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  secretID,
				SecretKey: secretKey,
			},
		})
		infos[i] = &types.BucketInfo{
			BucketName:  bucket.Name,
			DisplayName: displayName,
			COSClient:   bucketClient,
			Ak:          secretID,
			Sk:          secretKey,
		}
	}
	return
}
