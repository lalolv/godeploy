# Automatic deployment service

[中文说明](./README_ZH.md)

- [x] Multiple project deployment is supported
- [x] Slack Message Sending
- [x] logs
- [x] Support for Slash commands

## Use Docker

Run in bash：

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

- Mount the configuration file app.conf and directory Shells
- Mount the local time and time zone

## How to use

- Execution URL: `{domain name}:{port}/deploy/{Execute script name}`. The script name does not contain the `.sh` extension.
- In the `shells` directory, save the script files (with the extension xxx.sh).

For example: POST `http://127.0.0.1:8080/deploy/demo`，run `shells/demo.sh` script file。

## Configuration

File name：app.conf

- [server] port: Running port number
- [shell] path: Directory where the script is saved
- [slack] slack: Slack config. If it is empty, the API is not connected. Token is `Bot User OAuth Token` value in `OAuth & Permissions`。
- [slash] Slash Commands in Slack API，You can add command and run it in the Slack input field.

## Plan

- Support for plug-in extensions
- Add Token authentication
