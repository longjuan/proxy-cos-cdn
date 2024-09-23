package proxy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tencentyun/cos-go-sdk-v5"
	"proxy-cos-cdn/types"
)

var clientMap = make(map[string]*types.BucketInfo)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	info := clientMap[r.Host]
	if info == nil {
		log(r, http.StatusNotFound)
		http.Error(w, "can not find bucket", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		log(r, http.StatusNotFound)
		http.Error(w, "not support http method", http.StatusNotFound)
		return
	}

	client := info.COSClient
	objectKey := r.URL.Path[1:]

	if objectKey == "" {
		log(r, http.StatusNotFound)
		http.Error(w, "not support none key", http.StatusNotFound)
		return
	}

	queryParams := r.URL.Query()

	opt := &cos.PresignedURLOptions{
		Query: &queryParams,
	}

	// 获取签名
	signedURL, err := client.Object.GetPresignedURL(context.Background(), http.MethodGet, objectKey, info.Ak, info.Sk,
		time.Hour, opt)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		log(r, http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 发起请求并转发响应
	resp, err := http.Get(signedURL.String())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		log(r, http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	// 转发响应头
	for k, v := range resp.Header {
		for _, vi := range v {
			w.Header().Set(k, vi)
		}
	}
	w.WriteHeader(resp.StatusCode)

	// 转发响应体
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return
	}

	log(r, resp.StatusCode)
}

func log(r *http.Request, statusCode int) {
	logrus.WithFields(logrus.Fields{
		"host":            r.Host,
		"url":             r.URL,
		"ip":              r.RemoteAddr,
		"X-Forwarded-For": r.Header["X-Forwarded-For"],
		"statusCode":      statusCode,
	}).Info("proxy to cos")
}

func StartProxy(infos []*types.BucketInfo, addr string) {
	for _, info := range infos {
		clientMap[info.Domain] = info
	}
	http.HandleFunc("/", proxyHandler)
	logrus.Info("Start server on " + addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		return
	}
}
