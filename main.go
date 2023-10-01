package main

import (
	"flag"
	"fmt"
	"os"
	"proxy-cos-cdn/check"
	"proxy-cos-cdn/proxy"
	"proxy-cos-cdn/tencentyun"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

func main() {
	var (
		bucketRegion = flag.String("bucket-region", "", "Bucket的Region（必填）")
		secretID     = flag.String("secret-id", "", "腾讯云Secret ID（必填）")
		secretKey    = flag.String("secret-key", "", "腾讯云Secret Key（必填）")
		domainSuffix = flag.String("domain-suffix", "", "域名后缀（必填）")
		cdnCheck     = flag.Bool("cdn-check", false, "是否检查CDN域名是否正常解析，默认为false")
		port         = flag.Int64("port", 3321, "绑定端口，默认为3321")
	)

	*bucketRegion = getEnvVar("BUCKET_REGION", *bucketRegion)
	*secretID = getEnvVar("SECRET_ID", *secretID)
	*secretKey = getEnvVar("SECRET_KEY", *secretKey)
	*domainSuffix = getEnvVar("DOMAIN_SUFFIX", *domainSuffix)
	*cdnCheck = getEnvVarBool("CDN_CHECK", *cdnCheck)
	*port = getEnvVarInt("PORT", *port)

	flag.Parse()

	if *bucketRegion == "" || *secretID == "" || *secretKey == "" || *domainSuffix == "" {
		fmt.Println("缺少必填参数，请提供必填参数:")
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

func getEnvVar(envName string, defaultValue string) string {
	envValue := os.Getenv(envName)
	if envValue != "" && defaultValue == "" {
		return envValue
	}
	return defaultValue
}

func getEnvVarBool(envName string, defaultValue bool) bool {
	envValue := os.Getenv(envName)
	if envValue != "" {
		value, err := strconv.ParseBool(envValue)
		if err == nil {
			return value
		}
	}
	return defaultValue
}

func getEnvVarInt(envName string, defaultValue int64) int64 {
	envValue := os.Getenv(envName)
	if envValue != "" {
		value, err := strconv.ParseInt(envValue, 10, 64)
		if err == nil {
			return value
		}
	}
	return defaultValue
}
