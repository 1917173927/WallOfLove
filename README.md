# WallOfLove 项目

WallOfLove 是一个基于Gin框架的社交平台后端服务，提供用户管理、帖子发布、评论互动等功能。

## 功能特性

- 用户认证(注册/登录)
- 个人资料管理
- 帖子管理(创建/更新/删除/查看)
- 评论和回复功能
- 黑名单管理
- 图片上传
- 点赞功能
- 热门帖子排行

## 技术栈

- **编程语言**: Go 1.25
- **Web框架**: Gin
- **数据库**: MySQL + GORM
- **缓存**: Redis
- **认证**: JWT
- **定时任务**: Cron
- **配置管理**: Viper
- **日志**: Zap

## 环境要求

- Go 1.25+
- MySQL 8.0+
- Redis 6.0+

## 安装指南

1. 克隆仓库:
```bash
git clone https://github.com/1917173927/WallOfLove.git
cd WallOfLove
```

2. 安装依赖:
```bash
go mod download
```

3. 配置数据库和Redis连接信息(修改conf/config/config.go)

4. 启动服务:
```bash
go run main.go
```

## 配置说明

项目使用Viper管理配置，主要配置项包括:

- 数据库连接信息
- Redis连接信息
- JWT密钥
- 文件上传限制

配置文件位于`conf/config/config.go`

## API文档

项目API遵循RESTful规范，主要端点包括:

### 用户认证
- `POST /api/register` - 用户注册
  - 参数:
    - `username`: 必填，字符串
    - `name`: 必填，字符串
    - `password`: 必填，长度8-16位
    - `gender`: 必填，整数(0:男，1:女，2:保密)
    - `avatar_path`: 可选，字符串

- `POST /api/login` - 用户登录
  - 参数:
    - `username`: 必填，字符串
    - `password`: 必填，字符串

### 帖子管理
- `POST /api/post` - 创建帖子
  - 参数:
    - `content`: 必填，字符串
    - `anonymous`: 可选，布尔值
    - `visibility`: 可选，布尔值
    - `scheduled_at`: 可选，时间字符串

- `PUT /api/post` - 更新帖子
  - 参数:
    - `id`: 必填，整数
    - `content`: 可选，字符串
    - `anonymous`: 可选，布尔值
    - `visibility`: 可选，布尔值

- `DELETE /api/post` - 删除帖子
  - 参数:
    - `post_id`: 必填，整数

- `GET /api/post/list` - 获取帖子列表
  - 参数:
    - `page_size`: 可选，整数
    - `page`: 可选，整数

- `GET /api/post/:id` - 获取单个帖子
  - URL参数:
    - `id`: 必填，整数

### 互动功能
- `POST /api/review` - 发表评论
- `POST /api/reply` - 回复评论
- `POST /api/like` - 点赞帖子

完整API文档请参考路由配置文件`conf/route/route.go`和控制器代码

## 贡献指南

欢迎提交Pull Request或Issue。提交代码前请确保:

1. 通过所有测试
2. 遵循项目代码风格
3. 更新相关文档
