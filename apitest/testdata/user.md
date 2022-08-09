# 用户接口文档

## 获取用户信息

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

## 添加用户信息

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

## 导入用户信息(以csv文件格式)

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

