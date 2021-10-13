# 自动部署服务

- [x] 支持多个项目部署
- [x] Slack 消息发送
- [x] 日志
- [x] 支持 Slash 命令

## 使用 Docker 安装

运行命令：

```shell
docker run -it -p 8080:8080 \
-v ~/app/app.conf:/app/app.conf \
-v ~/app/log:/app/log \
-v ~/app/shells:/app/shells \
-v /etc/timezone:/etc/timezone \
-v /etc/localtime:/etc/localtime \
--restart always \
--name godeploy \
lalolv/godeploy
```

- 映射配置文件 app.conf 和 shells 目录
- 映射本地时间和时区
- 映射本地的命令 sh

## 使用说明

- 执行路径：{域名}:{端口}/deploy/{执行脚本名称}。执行脚本名称中不带有.sh 扩展名。
- 在 shells 目录下面，保存对应的脚本文件名（扩展名为 xxx.sh）。

例如：访问 POST `http://127.0.0.1:8080/deploy/demo`，执行 `shells/demo.sh` 脚本。

## 配置

配置文件：app.conf

- [server] port 运行端口号
- [shell] path: 脚本存放的路径
- [slack] slack 接口信息，如果为空，则不会对接 API。Token 为 `OAuth & Permissions` 目录下配置项 `Bot User OAuth Token` 的值。
- [slash] Slack API 中的 Slash Commands，可以添加对应的指令，然后在 Slack 对话框中执行。

## 开发计划

- 支持插件扩展
- 添加 Token 身份验证
