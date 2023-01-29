package tencentyun

import (
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
)

var client *cdn.Client

func InitCDNClient(secretID string, secretKey string) {
	credential := common.NewCredential(secretID, secretKey)
	client, _ = cdn.NewClient(credential, regions.Guangzhou, profile.NewClientProfile())
}

func GetDomainStatus(domain []*string) string {
	var (
		offset           = int64(0)
		limit            = int64(1)
		domainFilterName = "domain"
		baseRequest      = tchttp.NewCommonRequest("cdn", "2018-06-06", "DescribeDomains").BaseRequest
	)
	domains, err := client.DescribeDomains(&cdn.DescribeDomainsRequest{
		BaseRequest: baseRequest,
		Offset:      &offset,
		Limit:       &limit,
		Filters: []*cdn.DomainFilter{{
			Name:  &domainFilterName,
			Value: domain,
		}},
	})
	if err != nil {
		panic(err)
	}
	if len(domains.Response.Domains) == 0 {
		return "not configured"
	} else {
		return *domains.Response.Domains[0].Status
	}

}
