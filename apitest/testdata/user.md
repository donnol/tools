# 用户接口文档

**目录**：

* <a href="#获取用户信息"><b>获取用户信息 -- GET /api/user</b></a>

* <a href="#添加用户信息"><b>添加用户信息 -- POST /api/user</b></a>

* <a href="#导入用户信息(以csv文件格式)"><b>导入用户信息(以csv文件格式) -- POST /api/user/import</b></a>

## <a name="获取用户信息" href="#获取用户信息">获取用户信息</a>

`GET /api/user`

Request header:
- Content-Type: application/json; charset=utf-8

Response header:
- Content-Type: application/json; charset=utf-8

Param - Query

* userId (*uint*) 
* name (*string*) 

Return

* userId (*uint*) 
* name (*string*) 

Example:

<details>
<summary>Param</summary>

```json
name=jd&userId=1
```

</details>

<details>
<summary>Return</summary>

```json
{
    "userId": 1,
    "name": "jd"
}
```

</details>

## <a name="添加用户信息" href="#添加用户信息">添加用户信息</a>

`POST /api/user`

Request header:
- Content-Type: application/json; charset=utf-8

Response header:
- Content-Type: application/json; charset=utf-8

Param - Body

* userId (*uint*) 
* name (*string*) 

Return

* userId (*uint*) 
* name (*string*) 

Example:

<details>
<summary>Param</summary>

```json
{
    "userId": 1,
    "name": "jd"
}
```

</details>

<details>
<summary>Return</summary>

```json
{
    "userId": 1,
    "name": "jd"
}
```

</details>

## <a name="导入用户信息(以csv文件格式)" href="#导入用户信息(以csv文件格式)">导入用户信息(以csv文件格式)</a>

`POST /api/user/import`

Request header:
- Content-Type: text/csv; charset=utf-8

Response header:
- Content-Type: application/json; charset=utf-8

Param - Body

* userId (*uint*) 
* name (*string*) 

Return

* userId (*uint*) 
* name (*string*) 

Example:

<details>
<summary>Param</summary>

```json
{
    "userId": 1,
    "name": "jd"
}
```

</details>

<details>
<summary>Return</summary>

```json
{
    "userId": 1,
    "name": "jd"
}
```

</details>

