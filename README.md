# urlShortener
接口概述
该接口用于生成一个短链接，将用户提供的原始URL重定向到指定的短链接路径。用户可以选择自定义路径和设置过期时间。

请求URL
复制
POST 8.136.122.238:8080/api/url
请求参数
参数名	类型	必填	描述
original_url	string	是	需要重定向的原始URL。
custom_code	string	否	自定义的短链接路径。如果未提供，系统将随机生成一个路径。
duration	int	否	短链接的过期时间，单位为小时。如果未提供，将使用默认的过期时间。

请求示例
json
复制
{
"original_url": "https://www.youtube.com",
"custom_code": "2345",
"duration": 1
}
返回参数
参数名	类型	描述
short_url	string	生成的短链接URL。
original_url	string	原始URL，即用户提供的需要重定向的URL。
expire_time	string	短链接的过期时间，格式为YYYY-MM-DD HH:MM:SS。
返回示例
json
复制
{
"short_url": "8.136.122.238:8080/2345",
"original_url": "https://www.youtube.com",
"expire_time": "2024-12-25 17:35:27"
}
错误响应
HTTP状态码	错误码	描述
400	1001	original_url 参数缺失或格式不正确。
400	1002	custom_code 参数已存在，请选择其他路径。
400	1003	duration 参数无效，必须为正整数。
500	2001	服务器内部错误，无法生成短链接。
错误响应示例
json
复制
{
"error_code": 1001,
"message": "original_url parameter is missing or invalid"
}
注意事项
如果未提供 custom_code，系统将自动生成一个唯一的路径。

如果未提供 duration，短链接将使用系统默认的过期时间（例如24小时）。

custom_code 必须是唯一的，如果已存在，系统将返回错误。

示例调用
bash
复制
curl -X POST "http://8.136.122.238:8080/generate_short_url" \
-H "Content-Type: application/json" \
-d '{
"original_url": "https://www.youtube.com",
"custom_code": "2345",
"duration": 1
}'
返回结果
json
复制
{
"short_url": "8.136.122.238:8080/2345",
"original_url": "https://www.youtube.com",
"expire_time": "2024-12-25 17:35:27"
}
