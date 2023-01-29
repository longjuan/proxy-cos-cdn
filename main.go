package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"proxy-cos-cdn/check"
	"proxy-cos-cdn/proxy"
	"proxy-cos-cdn/tencentyun"
)

func main() {
	var (
		bucketRegion = flag.String("bucketRegion", "", "桶的Region（必填）")
		secretID     = flag.String("secretID", "", "腾讯云SecretID（必填）")
		secretKey    = flag.String("secretKey", "", "腾讯云SecretKey（必填）")
		domainSuffix = flag.String("domainSuffix", "", "域名（必填）")
		cdnCheck     = flag.Bool("cdnCheck", false, "是否检查cdn域名是否正常解析，默认为false")
		port         = flag.Int64("port", 3321, "绑定端口，默认为3321")
	)
	flag.Parse()

	// 检查必填项，若必填项没有填就输出--help的内容
	if *bucketRegion == "" || *secretID == "" || *secretKey == "" || *domainSuffix == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	infos := tencentyun.InitCosClients(*bucketRegion, *secretID, *secretKey)

	check.DNSRecord(infos, *domainSuffix)
	if *cdnCheck {
		check.CDNReady(infos, *secretID, *secretKey)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"BucketName", "DisplayName", "Domain", "DNSResult", "CDNReady"})
	for _, info := range infos {
		table.Append([]string{info.BucketName, info.DisplayName, info.Domain, strconv.FormatBool(info.DNSResult),
			info.CDNReady})
	}
	table.Render()

	proxy.StartProxy(infos, ":"+strconv.FormatInt(*port, 10))
}
