# Automatic deployment service

[中文说明](./README_ZH.md)

- [x] Multiple project deployment is supported
- [x] Slack Message Sending
- [x] logs

## How to use

- Execution URL: `{domain name}:{port}/deploy/{Execute script name}`. The script name does not contain the `.sh` extension.
- In the `shells` directory, save the script files (with the extension xxx.sh).

For example: POST `http://127.0.0.1:8080/deploy/demo`，run `shells/demo.sh` script file。

## Plan

- Support for Slack commands
- Support for plug-in extensions
- Add Token authentication
