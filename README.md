# URL Shortener

### 接口概述
该接口用于生成一个短链接，将用户提供的原始 URL 重定向到指定的短链接路径。用户可以选择自定义路径和设置过期时间。

---

### 请求 URL
POST 8.136.122.238/api/url

---

### 请求参数
| 参数名            | 类型     | 必填 | 描述                               |
|----------------|--------|----|----------------------------------|
| `original_url` | string | 是  | 需要重定向的原始 URL。                    |
| `custom_code`  | string | 否  | 自定义的短链接路径。如果未提供，系统将随机生成一个路径。     |
| `duration`     | int    | 否  | 短链接的过期时间，单位为小时。如果未提供，将使用默认的过期时间。 |

---

### 请求示例
```json
{
  "original_url": "https://www.youtube.com",
  "custom_code": "2345",
  "duration": 1
}
```
---

### 返回参数
| 参数名            | 类型     | 描述                                |
|----------------|--------|-----------------------------------|
| `short_url`    | string | 生成的短链接 URL。                       |
| `original_url` | string | 原始 URL，即用户提供的需要重定向的 URL。          |
| `expire_time`  | string | 短链接的过期时间，格式为 YYYY-MM-DD HH:MM:SS。 |

---

### 返回示例
```json
{
  "short_url": "8.136.122.238:8080/2345",
  "original_url": "https://www.youtube.com",
  "expire_time": "2024-12-25 17:35:27"
}
```

## 第一次运行时初始化
到根目录下执行 
```makefile
make make_db
```
设置环境变量HOST_NAME为服务器的ip
编译并运行main.go
## 之后运行
```makefile
make start_db
```
直接运行可执行文件main

## 版本更新

### v1.0.0 (2024-12-25)
原始版本

