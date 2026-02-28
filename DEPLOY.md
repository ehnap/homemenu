# HomeMenu 部署指南

## Docker 部署（推荐）

### 前置条件

- Docker >= 20.10
- Docker Compose >= 2.0

### 快速启动

```bash
# 1. 构建镜像
make docker-build

# 2. 启动服务
make docker-up

# 3. 验证运行状态
curl http://localhost:8080/api/health
# 返回 {"status":"ok"} 即为正常

# 查看日志
docker compose logs -f

# 停止服务
make docker-down
```

### 配置方式

支持两种配置方式，可混合使用（环境变量优先级高于 config.yaml）。

#### 方式一：环境变量（推荐用于 Docker）

编辑 `docker-compose.yml` 中的 `environment` 部分：

```yaml
services:
  homemenu:
    environment:
      - HOMEMENU_JWT_SECRET=your-secret-key    # 必须修改！
      - HOMEMENU_PORT=8080
      - HOMEMENU_DB_PATH=./data/homemenu.db
      - HOMEMENU_UPLOAD_DIR=./data/uploads
      - HOMEMENU_SHARE_BASE_URL=https://your-domain.com
      # LLM 配置（可选，不配则使用规则引擎生成菜单）
      - HOMEMENU_LLM_BASE_URL=https://api.openai.com/v1
      - HOMEMENU_LLM_API_KEY=sk-xxx
      - HOMEMENU_LLM_MODEL=gpt-4o
```

#### 方式二：配置文件

创建 `config.yaml` 并挂载到容器中：

```yaml
server:
  port: 8080
  jwt_secret: "your-secret-key"

db:
  driver: sqlite
  path: ./data/homemenu.db

llm:
  base_url: ""
  api_key: ""
  model: ""

upload:
  dir: ./data/uploads
  max_size_mb: 10

share:
  base_url: "https://your-domain.com"
```

在 `docker-compose.yml` 中取消注释配置文件挂载：

```yaml
volumes:
  - ./data:/app/data
  - ./config.yaml:/app/config.yaml  # 取消注释此行
```

### 环境变量一览

| 环境变量 | 默认值 | 说明 |
|---------|--------|------|
| `HOMEMENU_PORT` | `8080` | 服务端口 |
| `HOMEMENU_JWT_SECRET` | `change-me-in-production` | JWT 密钥，**生产环境必须修改** |
| `HOMEMENU_DB_PATH` | `./data/homemenu.db` | SQLite 数据库路径 |
| `HOMEMENU_UPLOAD_DIR` | `./data/uploads` | 图片上传目录 |
| `HOMEMENU_SHARE_BASE_URL` | `http://localhost:8080` | 分享链接前缀 |
| `HOMEMENU_LLM_BASE_URL` | 空 | LLM API 地址（兼容 OpenAI 协议） |
| `HOMEMENU_LLM_API_KEY` | 空 | LLM API Key |
| `HOMEMENU_LLM_MODEL` | 空 | LLM 模型名 |

### 数据持久化

容器内 `/app/data` 目录包含所有持久数据：

- `homemenu.db` — SQLite 数据库
- `uploads/` — 用户上传的图片

通过 `volumes: ./data:/app/data` 映射到宿主机，确保容器重建后数据不丢失。

### 反向代理（Nginx 示例）

```nginx
server {
    listen 80;
    server_name your-domain.com;

    client_max_body_size 10m;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

使用反向代理时，记得将 `HOMEMENU_SHARE_BASE_URL` 设为实际域名。

---

## 裸机部署

### 前置条件

- Go >= 1.24
- Node.js >= 20

### 构建与运行

```bash
# 构建（前端 + 后端，输出单二进制 homemenu）
make build

# 创建数据目录
mkdir -p data/uploads

# 运行（可选指定配置文件路径）
./homemenu
# 或
./homemenu /path/to/config.yaml
```

不提供 config.yaml 时使用默认值，也可通过环境变量配置：

```bash
HOMEMENU_JWT_SECRET=your-secret HOMEMENU_SHARE_BASE_URL=https://your-domain.com ./homemenu
```

### systemd 服务（可选）

```ini
# /etc/systemd/system/homemenu.service
[Unit]
Description=HomeMenu
After=network.target

[Service]
Type=simple
User=homemenu
WorkingDirectory=/opt/homemenu
ExecStart=/opt/homemenu/homemenu
Environment=HOMEMENU_JWT_SECRET=your-secret-key
Environment=HOMEMENU_SHARE_BASE_URL=https://your-domain.com
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable --now homemenu
```
