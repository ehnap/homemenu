# CLAUDE.md - HomeMenu（家味）家庭菜谱与菜单规划应用

## 项目概述

家庭菜谱管理与一周菜单智能规划工具。支持菜谱录入、智能/规则菜单生成、拖拽编辑、购物清单生成与分享。

**目标用户**：家庭成员（2-5人），自托管部署。未来可能开放给更多用户使用。

## 技术栈

- **前端**：Vue 3 + TypeScript + Vite + Tailwind CSS + Pinia
- **后端**：Go (Gin 或 Echo) + SQLite（默认）
- **部署**：Docker Compose + 裸机单二进制 双支持
- **拖拽库**：vuedraggable / @vueuse/integrations 的 useSortable

## 架构原则

### 前后端分离
- Go 提供 REST API
- 生产模式 Go 内嵌前端静态资源（go:embed）
- 开发模式前后端独立运行，Vite proxy 转发 API

### Repository 接口抽象（重要）
所有数据访问通过 Repository 接口，service 层不直接写 SQL：

```go
type RecipeRepo interface {
    Create(ctx context.Context, recipe *Recipe) error
    GetByID(ctx context.Context, id int64) (*Recipe, error)
    List(ctx context.Context, filters RecipeFilters) ([]Recipe, error)
    Update(ctx context.Context, recipe *Recipe) error
    Delete(ctx context.Context, id int64) error
}
```

MVP 实现 SQLite 版本，后续可无缝切换 PostgreSQL/MySQL，只需新增实现并在配置中切换。所有 Repo 接口统一遵循此模式。

## 数据模型

### users 用户表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PK | 自增ID |
| username | TEXT UNIQUE | 用户名 |
| password_hash | TEXT | 密码哈希 |
| nickname | TEXT | 昵称 |
| created_at | DATETIME | 创建时间 |

### recipes 菜谱表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PK | 自增ID |
| user_id | INTEGER FK | 创建者 |
| name | TEXT NOT NULL | 菜名 |
| steps | JSON | 步骤数组 [{order, description, image_url?}] |
| cook_time | INTEGER | 烹饪时长（分钟），可选 |
| difficulty | TEXT | 简单/中等/复杂，可选 |
| tags | JSON | 标签数组（口味：咸/辣/清淡，菜系，荤/素等），可选 |
| cover_image | TEXT | 封面图URL，可选 |
| calories | INTEGER | 卡路里，可选 |
| notes | TEXT | 备注，可选 |
| created_at | DATETIME | |
| updated_at | DATETIME | |

### recipe_ingredients 食材表（独立表，支持按食材检索）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PK | 自增ID |
| recipe_id | INTEGER FK | 关联菜谱 |
| name | TEXT NOT NULL | 食材名（如"菠菜"） |
| amount | TEXT | 用量（如"200"） |
| unit | TEXT | 单位（如"克"） |

**对 name 建索引**，支持多食材组合查询：
```sql
SELECT recipe_id FROM recipe_ingredients WHERE name LIKE '%菠菜%'
INTERSECT
SELECT recipe_id FROM recipe_ingredients WHERE name LIKE '%猪肉%'
```

### meal_plans 周菜单表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PK | 自增ID |
| user_id | INTEGER FK | 创建者 |
| name | TEXT | 菜单名（如"第7周菜单"） |
| start_date | DATE | 开始日期 |
| end_date | DATE | 结束日期 |
| config | JSON | 生成条件（口味偏好、食材限制、荤素比例等） |
| share_token | TEXT UNIQUE | 分享令牌，用于无登录访问 |
| created_at | DATETIME | |
| updated_at | DATETIME | |

### meal_plan_items 菜单明细表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PK | 自增ID |
| meal_plan_id | INTEGER FK | 关联菜单 |
| recipe_id | INTEGER FK | 关联菜谱 |
| date | DATE | 哪天 |
| meal_type | TEXT | breakfast/lunch/dinner |
| sort_order | INTEGER | 同一餐内排序（拖拽用） |

## 核心功能

### 1. 菜谱管理（CRUD）
- **列表页**：卡片式展示，支持按标签、口味、荤素筛选，支持按食材名搜索（单个或多个食材组合）
- **详情页**：食材清单 + 分步骤做法 + 封面图
- **编辑页**：
  - 必填：菜名、至少一个食材
  - 可选：步骤、烹饪时间、难度、标签、封面图、卡路里、备注
  - 食材行动态增删
  - 步骤支持拖拽排序
  - 图片上传（存本地 /uploads 目录）

### 2. 智能菜单生成
**配置面板**（生成前设置）：
- 日期范围（默认当周一到周日）
- 每餐菜品数量（如午餐：2荤1素1汤）
- 餐次选择（勾选需要规划的：早/午/晚）
- 口味偏好（如本周偏清淡）
- 食材偏好（优先使用某些食材）
- 食材排除（不要某些食材）

**双引擎**：
- **规则引擎**（默认，无需配置）：
  - 从菜谱库随机抽取，满足荤素数量约束
  - 口味标签权重过滤
  - 同一周内菜品不重复
  - 食材偏好/排除过滤
- **AI 引擎**（可选，需配置 LLM）：
  - 将菜谱库摘要 + 约束条件发给 LLM
  - 兼容 OpenAI API 协议（base_url + api_key + model）
  - LLM 返回推荐的 recipe_id 列表
  - 若 LLM 不可用，自动降级到规则引擎

**配置判断逻辑**：config.yaml 中有 llm 配置且 api_key 非空 → 使用 AI 引擎，否则使用规则引擎。也可在前端配置面板手动切换。

### 3. 菜单编辑与拖拽
- **周视图**：7列（周一到周日）× N行（已选餐次）的网格布局
- 每个格子显示该餐的菜品卡片
- **拖拽**：菜品可在任意格子间拖拽移动
- **操作**：
  - 单道菜"换一道"（重新随机/AI推荐一道符合约束的）
  - 手动添加菜品（从菜谱库选择）
  - 删除某道菜
- 拖拽排序通过更新 sort_order 和 date/meal_type 实现

### 4. 购物清单
- 从菜单自动汇总所有食材
- 同名食材合并数量（单位相同时数值相加）
- **两种视图**：
  - 按天拆分：显示每天需要买的食材（适合每天采购）
  - 一周汇总：所有食材合并（适合一次性采购）
- 每项可勾选"已买"（状态存前端 localStorage 即可）

### 5. 分享功能
- 生成菜单分享链接：`{base_url}/share/{share_token}`
- **无需登录**即可查看
- 分享页展示：
  - 一周菜单总览
  - 每日购物清单
  - 点击菜名可查看详细做法和食材
- share_token 在创建菜单时自动生成（UUID）

## 用户系统（框架级，不深入实现）
- JWT 认证：登录返回 access_token + refresh_token
- 中间件校验 token
- 所有数据表预留 user_id 外键
- MVP 阶段：注册、登录、token 刷新三个接口即可
- 后续可扩展：角色权限、家庭组、邀请码等

## API 设计（RESTful）

```
# 用户
POST   /api/auth/register
POST   /api/auth/login
POST   /api/auth/refresh

# 菜谱
GET    /api/recipes              # 列表（支持 ?tag=&ingredient=&q= 筛选）
POST   /api/recipes              # 创建
GET    /api/recipes/:id          # 详情
PUT    /api/recipes/:id          # 更新
DELETE /api/recipes/:id          # 删除

# 菜单
GET    /api/meal-plans           # 列表
POST   /api/meal-plans           # 创建
GET    /api/meal-plans/:id       # 详情（含 items）
PUT    /api/meal-plans/:id       # 更新
DELETE /api/meal-plans/:id       # 删除
POST   /api/meal-plans/generate  # 智能生成

# 菜单明细
PUT    /api/meal-plans/:id/items # 批量更新（拖拽后保存）
POST   /api/meal-plans/:id/items/:itemId/reroll  # 换一道

# 购物清单
GET    /api/meal-plans/:id/shopping-list?mode=daily|weekly

# 分享（无需认证）
GET    /api/share/:token         # 获取分享的菜单数据

# 图片上传
POST   /api/upload               # 上传图片，返回 URL
```

## 项目结构

```
homemenu/
├── backend/
│   ├── main.go                  # 入口，初始化配置/DB/路由
│   ├── config/
│   │   └── config.go            # 配置加载
│   ├── db/
│   │   ├── sqlite.go            # SQLite 连接与迁移
│   │   └── migrations/          # SQL 迁移文件
│   ├── model/
│   │   ├── user.go
│   │   ├── recipe.go
│   │   ├── ingredient.go
│   │   ├── meal_plan.go
│   │   └── meal_plan_item.go
│   ├── repository/
│   │   ├── interfaces.go        # 所有 Repo 接口定义
│   │   └── sqlite/              # SQLite 实现
│   │       ├── recipe_repo.go
│   │       ├── ingredient_repo.go
│   │       ├── meal_plan_repo.go
│   │       └── user_repo.go
│   ├── service/
│   │   ├── recipe_service.go
│   │   ├── meal_plan_service.go
│   │   ├── generator/
│   │   │   ├── interface.go     # MenuGenerator 接口
│   │   │   ├── rule_engine.go   # 规则引擎实现
│   │   │   └── ai_engine.go     # AI 引擎实现
│   │   └── shopping_service.go
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── recipe_handler.go
│   │   ├── meal_plan_handler.go
│   │   ├── share_handler.go
│   │   └── upload_handler.go
│   ├── middleware/
│   │   └── auth.go              # JWT 中间件
│   └── static/                  # 构建后的前端文件（go:embed）
├── frontend/
│   ├── src/
│   │   ├── App.vue
│   │   ├── main.ts
│   │   ├── router/
│   │   │   └── index.ts
│   │   ├── stores/
│   │   │   ├── auth.ts
│   │   │   ├── recipe.ts
│   │   │   └── mealPlan.ts
│   │   ├── api/
│   │   │   ├── client.ts        # axios 实例 + 拦截器
│   │   │   ├── recipe.ts
│   │   │   ├── mealPlan.ts
│   │   │   └── auth.ts
│   │   ├── views/
│   │   │   ├── RecipeList.vue
│   │   │   ├── RecipeDetail.vue
│   │   │   ├── RecipeEdit.vue
│   │   │   ├── MealPlanList.vue
│   │   │   ├── MealPlanEditor.vue    # 周视图 + 拖拽
│   │   │   ├── MealPlanGenerate.vue  # 生成配置面板
│   │   │   ├── ShoppingList.vue
│   │   │   ├── ShareView.vue         # 分享页（无需登录）
│   │   │   └── Login.vue
│   │   ├── components/
│   │   │   ├── RecipeCard.vue
│   │   │   ├── MealSlot.vue          # 单个餐位格子
│   │   │   ├── DraggableMealCard.vue # 可拖拽菜品卡片
│   │   │   ├── IngredientInput.vue   # 食材动态输入组件
│   │   │   ├── TagFilter.vue
│   │   │   └── ImageUpload.vue
│   │   └── types/
│   │       └── index.ts              # TypeScript 类型定义
│   ├── index.html
│   ├── vite.config.ts
│   ├── tailwind.config.js
│   └── tsconfig.json
├── config.yaml                  # 运行配置
├── docker-compose.yml
├── Dockerfile                   # 多阶段构建（前端build + Go build）
├── Makefile                     # dev/build/run 快捷命令
└── README.md
```

## 配置文件 config.yaml

```yaml
server:
  port: 8080
  jwt_secret: "change-me-in-production"

db:
  driver: sqlite           # 预留切换: sqlite | postgres | mysql
  path: ./data/homemenu.db # SQLite 专用

llm:                       # 可选，不配或 api_key 为空则用规则引擎
  base_url: ""
  api_key: ""
  model: ""

upload:
  dir: ./data/uploads
  max_size_mb: 10

share:
  base_url: "http://localhost:8080"
```

## 部署

### Docker Compose
```yaml
version: "3"
services:
  homemenu:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data        # SQLite + 上传图片持久化
      - ./config.yaml:/app/config.yaml
```

### 裸机
```bash
make build    # 构建前端 + Go 二进制
./homemenu    # 直接运行，读取 config.yaml
```

### Dockerfile（多阶段构建）
```dockerfile
# Stage 1: 前端构建
FROM node:20-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/ .
RUN npm ci && npm run build

# Stage 2: Go 构建
FROM golang:1.22-alpine AS backend
WORKDIR /app
COPY backend/ ./backend/
COPY --from=frontend /app/frontend/dist ./backend/static/
RUN cd backend && go build -o /homemenu .

# Stage 3: 运行
FROM alpine:3.19
COPY --from=backend /homemenu /app/homemenu
WORKDIR /app
EXPOSE 8080
CMD ["./homemenu"]
```

## 开发顺序建议

1. **后端基础**：项目初始化、配置加载、SQLite 连接、用户认证框架
2. **菜谱 CRUD**：Model + Repo + Service + Handler + 前端页面
3. **食材搜索**：recipe_ingredients 表 + 按食材筛选
4. **菜单生成 - 规则引擎**：生成逻辑 + 配置面板
5. **菜单编辑 + 拖拽**：周视图网格 + vuedraggable
6. **购物清单**：食材汇总 + 按天/周视图
7. **分享功能**：share_token + 无登录查看页
8. **AI 引擎**：LLM 集成 + 降级逻辑
9. **图片上传**：上传接口 + 前端组件
10. **Docker 部署**：Dockerfile + docker-compose

## 编码规范

- Go 代码遵循标准 Go 项目布局
- 前端组件使用 Vue 3 `<script setup>` + Composition API
- 所有 API 返回统一 JSON 格式：`{ code: 0, data: {}, message: "" }`
- 错误处理：Go 侧统一错误码，前端 axios 拦截器统一处理
- 数据库迁移：使用 SQL 文件，按版本号顺序执行
