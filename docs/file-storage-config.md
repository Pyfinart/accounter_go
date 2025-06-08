# 文件存储配置说明

## 配置选项

在 `config.yaml` 中的 `data.file_storage` 部分可以配置文件存储选项：

```yaml
data:
  file_storage:
    data_dir: "./data"              # 数据目录路径
    accounter_file: "accounters.json"  # 记账数据文件名
```

## 配置说明

### `data_dir`
- **类型**: string
- **默认值**: `"./data"`
- **说明**: 数据文件存储目录的路径，支持相对路径和绝对路径
- **示例**:
  - `"./data"` - 项目根目录下的data文件夹
  - `"/var/lib/accounter"` - 绝对路径
  - `"./storage/prod"` - 生产环境专用目录

### `accounter_file`
- **类型**: string
- **默认值**: `"accounters.json"`
- **说明**: 记账数据的JSON文件名
- **示例**:
  - `"accounters.json"` - 默认文件名
  - `"prod_accounters.json"` - 生产环境文件
  - `"backup_accounters.json"` - 备份文件

## 使用示例

### 开发环境配置
```yaml
data:
  file_storage:
    data_dir: "./storage/dev"
    accounter_file: "dev_accounters.json"
```

### 生产环境配置
```yaml
data:
  file_storage:
    data_dir: "/var/lib/accounter/data"
    accounter_file: "accounters.json"
```

### 测试环境配置
```yaml
data:
  file_storage:
    data_dir: "./test_data"
    accounter_file: "test_accounters.json"
```

## 默认行为

如果配置文件中没有 `file_storage` 配置，系统将使用以下默认值：
- 数据目录: `./data`
- 文件名: `accounters.json`
- 完整路径: `./data/accounters.json`

## 注意事项

1. **目录权限**: 确保应用程序对配置的数据目录有读写权限
2. **目录创建**: 如果配置的目录不存在，系统会自动创建
3. **文件路径**: 最终的文件路径为 `data_dir/accounter_file`
4. **备份建议**: 建议定期备份数据文件，特别是在生产环境中

## 切换到数据库存储

如果将来需要切换到数据库存储，只需要修改 `internal/data/data.go` 中的 `ProviderSet`：

```go
// 注释掉文件存储
// NewAccounterFileRepo,
// 启用数据库存储  
NewAccounterDbRepo,
``` 