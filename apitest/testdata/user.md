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

* id (*string*) id
* name (*string*) 名字
* age (*int*) 年龄
* addr (*object*) 地址
    * city (*string*) 城市
    * home (*string*) 家
* phone (*string*) 手机

Return

* id (*string*) id
* name (*string*) 名字
* age (*int*) 年龄
* addr (*object*) 地址
    * city (*string*) 城市
    * home (*string*) 家
* phone (*string*) 手机

Example:

<details>
<summary>Param</summary>

```json
addr=%7B+%7D&age=0&id=1&name=jd&phone=
```

</details>

<details>
<summary>Return</summary>

```json
{
    "id": "1", // id
    "name": "jd", // 名字
    "age": 0, // 年龄
    "addr": { // 地址
        "city": "", // 城市
        "home": "" // 家
    },
    "phone": "" // 手机
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

* id (*string*) id
* name (*string*) 名字
* age (*int*) 年龄
* addr (*object*) 地址
    * city (*string*) 城市
    * home (*string*) 家
* phone (*string*) 手机

Return

* id (*string*) id
* name (*string*) 名字
* age (*int*) 年龄
* addr (*object*) 地址
    * city (*string*) 城市
    * home (*string*) 家
* phone (*string*) 手机

Example:

<details>
<summary>Param</summary>

```json
{
    "id": "1", // id
    "name": "jd", // 名字
    "age": 0, // 年龄
    "addr": { // 地址
        "city": "", // 城市
        "home": "" // 家
    },
    "phone": "" // 手机
}
```

</details>

<details>
<summary>Return</summary>

```json
{
    "id": "1", // id
    "name": "jd", // 名字
    "age": 0, // 年龄
    "addr": { // 地址
        "city": "", // 城市
        "home": "" // 家
    },
    "phone": "" // 手机
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

* id (*string*) id
* name (*string*) 名字
* age (*int*) 年龄
* addr (*object*) 地址
    * city (*string*) 城市
    * home (*string*) 家
* phone (*string*) 手机

Return

* id (*string*) id
* name (*string*) 名字
* age (*int*) 年龄
* addr (*object*) 地址
    * city (*string*) 城市
    * home (*string*) 家
* phone (*string*) 手机

Example:

<details>
<summary>Param</summary>

```json
{
    "id": "1", // id
    "name": "jd", // 名字
    "age": 0, // 年龄
    "addr": { // 地址
        "city": "", // 城市
        "home": "" // 家
    },
    "phone": "" // 手机
}
```

</details>

<details>
<summary>Return</summary>

```json
{
    "id": "1", // id
    "name": "jd", // 名字
    "age": 0, // 年龄
    "addr": { // 地址
        "city": "", // 城市
        "home": "" // 家
    },
    "phone": "" // 手机
}
```

</details>

