# proxy-cos-cdn

proxy-cos-cdn 是一个用于在腾讯云服务器上设置反向代理，以提供私有读取腾讯云对象存储（COS）桶的功能。它允许服务器走内网流量，从而节省回源流量费用。

## 使用 Docker 镜像

你可以使用 Docker 镜像运行 proxy-cos-cdn，这是推荐的方式，因为它允许你轻松地设置必要的环境变量。

1. 拉取 Docker 镜像：

   ```bash
   docker pull ghcr.io/longjuan/proxy-cos-cdn:main
   ```

2. 设置必要的环境变量：

   使用 `-e` 标志来设置环境变量。以下是示例命令：

   ```bash
   docker run -e BUCKET_REGION=your_bucket_region \
              -e SECRET_ID=your_secret_id \
              -e SECRET_KEY=your_secret_key \
              -e DOMAIN_SUFFIX=your_domain_suffix \
              -e CDN_CHECK=true \  # (可选) 是否检查 CDN 域名是否正常解析，默认为 false
              ghcr.io/longjuan/proxy-cos-cdn:main
   ```
   
   请将 `your_bucket_region`、`your_secret_id`、`your_secret_key`、`your_domain_suffix` 替换为你的腾讯云 COS 凭证和域名信息。
   
3. 运行 Docker 容器，应用程序将执行检查并显示结果。

## 本地开发

要在本地进行开发，你需要具备以下先决条件：

- Go 编程环境：确保你已经安装了 Go 编程环境。
- 腾讯云 COS 凭证：需要提供腾讯云 COS 的凭证信息。

按照以下步骤进行操作：

1. 克隆存储库：

   ```bash
   git clone https://github.com/yourusername/proxy-cos-cdn.git
   cd proxy-cos-cdn
   ```

2. 构建应用程序：

   ```bash
   go build -o proxy-cos-cdn
   ```

3. 运行应用程序并使用控制台参数：

   ```bash
   ./proxy-cos-cdn -bucket-region your_bucket_region \
                   -secret-id your_secret_id \
                   -secret-key your_secret_key \
                   -domain-suffix your_domain_suffix \
                   -cdn-check true \  # (可选) 是否检查 CDN 域名是否正常解析，默认为 false
                   -port 3321 \       # (可选) 绑定端口，默认为 3321
   ```

应用程序将执行检查并显示结果。
