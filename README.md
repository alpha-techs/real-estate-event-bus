# 不动产项目助手

## 功能列表

### 已完成

- [x] 处理Github Actions发出的Webhook请求

- [x] 提供API接口用于更新后端服务

### 开发中

- [ ] 提供API接口用于更新前端服务

- [ ] 处理飞书卡片回调

- [ ] 接口认证

- [ ] 支持多环境

## 开发

### 本地运行

```bash
go run main.go
```

## 部署

### AWS EC2原生部署

```bash
nohup event-bus > /dev/null 2>&1 &
```
