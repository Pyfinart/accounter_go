# 💰 元元

一个基于Go-Kratos框架的现代化个人记账系统，提供完整的Web界面和RESTful API。

## ✨ 功能特性

### 📝 记账功能
- ✅ 添加收入/支出记录
- ✅ 支持多种分类（餐饮、交通、购物等）
- ✅ 自定义交易描述和日期
- ✅ 删除交易记录

### 📊 数据统计
- ✅ 总收入、总支出、余额统计
- ✅ 分类统计图表（饼图）
- ✅ 按时间范围筛选
- ✅ 按类型和分类筛选

### 🎨 用户界面
- ✅ 现代化响应式Web界面
- ✅ 移动端友好设计
- ✅ 实时数据更新
- ✅ 直观的图表展示

### 💾 数据存储
- ✅ 文件存储（JSON格式）
- ✅ 配置化存储路径
- ✅ 数据库支持（预留）
- ✅ 自动数据备份

## 🚀 快速开始

### 方法一：一键启动（推荐）
```bash
./start.sh
```

### 方法二：手动启动
1. **启动API服务器**
```bash
go build -o build/accounter ./cmd/accounter
./build/accounter -conf ./configs/config.yaml
```

2. **启动Web界面**
```bash
cd web
go run server.go
```

### 访问系统
- 🌐 **Web界面**: http://localhost:3000
- 🔌 **API服务**: http://localhost:8000

## 📋 API接口

### 添加交易记录
```bash
curl -X POST http://localhost:8000/api/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "type": 2,
    "category": 2,
    "desc": "午餐",
    "amount": 25.50,
    "date": "2024-01-15"
  }'
```

### 查询交易记录
```bash
curl http://localhost:8000/api/transactions
```

### 获取统计数据
```bash
curl http://localhost:8000/api/stats
```

### 删除交易记录
```bash
curl -X DELETE http://localhost:8000/api/transactions/1
```

## 🔧 配置说明

### 文件存储配置
在 `configs/config.yaml` 中配置：
```yaml
data:
  file_storage:
    data_dir: "./data"              # 数据目录
    accounter_file: "accounters.json"  # 数据文件名
```

### 不同环境配置
- `configs/config.yaml` - 生产环境
- `configs/config-dev.yaml` - 开发环境

## 📊 数据分类

### 交易类型
- `1` - 收入
- `2` - 支出

### 分类列表
| 编号 | 分类 | 编号 | 分类 |
|------|------|------|------|
| 1 | 游戏 | 10 | 投资 |
| 2 | 餐饮 | 11 | 借款 |
| 3 | 旅行 | 12 | 工资 |
| 4 | 教育 | 13 | 其他收入 |
| 5 | 健康 | 14 | 应用 |
| 6 | 购物 | 15 | 住房 |
| 7 | 其他 | 16 | 水电费 |
| 8 | 交通 | 17 | 礼物 |
| 9 | 娱乐 | 18 | 零食 |

## 🏗️ 项目结构

```
accounter_go/
├── api/                    # API定义
│   └── accounter/v1/      # Proto文件和生成代码
├── cmd/accounter/         # 主程序入口
├── configs/               # 配置文件
├── internal/              # 内部代码
│   ├── biz/              # 业务逻辑层
│   ├── data/             # 数据访问层
│   ├── service/          # 服务层
│   └── server/           # 服务器配置
├── web/                   # Web界面
│   ├── static/           # 静态文件
│   └── server.go         # Web服务器
├── data/                  # 数据文件目录
└── docs/                  # 文档
```

## 🔄 切换存储方式

### 当前：文件存储
数据保存在JSON文件中，便于查看和备份。

### 切换到数据库存储
1. 修改 `internal/data/data.go` 中的 `ProviderSet`：
```go
// 注释掉文件存储
// NewAccounterFileRepo,
// 启用数据库存储
NewAccounterDbRepo,
```

2. 配置数据库连接信息
3. 重新编译运行

## 🛠️ 开发说明

### 重新生成Proto代码
```bash
cd api/accounter/v1
protoc --proto_path=. \
  --proto_path=../../../third_party \
  --go_out=paths=source_relative:. \
  --go-http_out=paths=source_relative:. \
  --go-grpc_out=paths=source_relative:. \
  accounter.proto
```

### 重新生成依赖注入代码
```bash
cd cmd/accounter
go generate
```

## 📱 移动端支持

Web界面采用响应式设计，完美支持：
- 📱 手机浏览器
- 📱 平板设备
- 💻 桌面浏览器

## 🔒 安全说明

- 当前版本为单用户系统
- 数据存储在本地文件
- 建议定期备份数据文件
- 生产环境请配置HTTPS


