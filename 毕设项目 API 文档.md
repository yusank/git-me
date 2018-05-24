#毕设项目 API 文档

[TOC]

## 更新记录

2018.05.23：添加收藏模块接口

2018.05.24：更新解释 url 接口

## 用户系统

*所有的 POST 请求参数都会以 json 格式传输*

**HOST:http://45.76.169.195:17717**

###注册

| 路由 | /v1/r/register               |
| ---- | ---------------------------- |
| 参数 | {name:"", email:"", pass:""} |
| 方法 | POST                         |

成功返回：

```json
{
  "data": "success",
  "errcode": 0
}
```



### 登录

| 路由 | /v1/r/login       |
| ---- | ----------------- |
| 参数 | {name:"",pass:""} |
| 方法 | POST              |

成功返回：

```json
{
    "errcode": 0,
    "data": {
        // 用户信息
        "id": "5adc7dad65d4406e6e984585",
        "name": "yusank",
        "nickname": "iwElOCeE",
        "email": "abc@gmail.com",
        "headImg": "",
        "createdAt": 1524399750,
        "updatedAt": 0
    }
}
```



### 登出

| 路由 | /v1/user/logout |
| ---- | --------------- |
| 方法 | GET             |

成功返回：

```json
{
    "data": "success",
    "errcode": 0
}
```



### 修改用户信息

| 路由 | /v1/user/info                     |
| ---- | --------------------------------- |
| 参数 | {name:"", headImg:"",nickname:""} |
| 方法 | POST                              |

成功返回：

```json
{
    "errcode": 0,
    "data": {
        // 返回新的用户信息
        "id": "5adc7dad65d4406e6e984585",
        "name": "yusank",
        "nickname": "newnikcname",
        "email": "abc@gmail.com",
        "headImg": "",
        "createdAt": 1524399750,
        "updatedAt": 0
    }
}
```

### 更改密码

| 路由 | /v1/user/pass                     |
| ---- | --------------------------------- |
| 参数 | {name:"", oldPass:"", newPass:""} |
| 方法 | POST                              |

成功返回：

```json
{
    "data": "success",
    "errcode": 0
}
```



## 解析URL

| 路由 | /v1/download/vid                               |
| ---- | ---------------------------------------------- |
| 参数 | {"url":"http://baidu.com/a.jpg", "type":0或1}  |
| 方法 | POST                                           |
| 说明 | 参数中的 type 为0表示，检测 url；1时表示，预览 |

**注：参数为 json 格式**

成功返回：

```json
{
  "data": {
      "isUseful"：true/false	// 调接口时，type 为0时，返回该响应，表示该 url 是否可用
  },
  "errcode": 0
}
```

**如果传参的 type=1 ，则返回二进制文件流，是大小为2m 的视频文件，用来预览**

## 历史记录

| 路由 | /v1/history/list |
| ---- | ---------------- |
| 参数 | page=1&size=10   |
| 方法 | get              |

返回数据：

```json
{
    "data":{
        [
        	"id":"123",
        	"userId":"123",
        	"site":"bilibili.com",
        	"title":"视频的 title"
        	"url":"http://abc.com",
        	"size":1024,
        	"quality":"高清",
        	"type":"video",
        	"createdAt":1502020220,// 创建时间
        ],[]....
    },
    "errCode":0
}
```



## 预定下载

### 获取用户的预定列表

| 路由 | /v1/task/list/:id |
| ---- | ----------------- |
| 参数 | page=1&size=10    |
| 方法 | get               |

返回值：

```json
{
    "data":{
        [
        	"id":"123",
        	"userId":"123",
        	"status":0,// 0:未开始，1：正在下载， 2：已完成，3：下载异常
        	"url":"http://abc.com",
        	"sort":1,
        	"type":"video",
        	"createdAt":1502020220,// 创建时间
        	"updatedAt":15003494934 // 更新时间
        ],[]....
    },
    "errCode":0
}
```



请求 通用json 结构

```json
{
    "id":"sdds",
    "userId":"abc",
    "url":"http://abc.com",
    "sort":1	//用来排序
}
```



### 添加预定

| 路由 | /v1/task/add         |
| ---- | -------------------- |
| 参数 | 只传 url 和 sort即可 |
| 方法 | post                 |

返回值:

```json
{
    "data": "success",
    "errcode": 0
}
```



### 更新预定

| 路由 | /v1/task/upadte |
| ---- | --------------- |
| 参数 | 传 id           |
| 方法 | POST            |

返回值:

```json
{
    "data": "success",
    "errcode": 0
}
```



### 删除预定

| 路由 | /v1/task/del |
| ---- | ------------ |
| 参数 | 传 id        |
| 方法 | POST         |

返回值:

```json
{
    "data": "success",
    "errcode": 0
}
```



## 收藏

### 添加收藏

| 路由 | /v1/collect/add                  |
| ---- | -------------------------------- |
| 参数 | {"url":"http://baidu.com/a.jpg"} |
| 方法 | POST                             |

**注：参数为 json 格式**

返回值:

```json
{
    "data": "success",
    "errcode": 0
}
```



### 获取收藏列表

| 路由 | /v1/collect/list |
| ---- | ---------------- |
| 参数 | page=1&size=10   |
| 方法 | GET              |

返回参数

```json
{
    "data":{
        [
        	"id":"123",
        	"userId":"123",
        	"url":"http://abc.com",
        	"site":"baidu.com",
        	"size":1024,
        	"title":"资源的 title"
        	"quality":"高清",
        	"createdAt":1502020220,// 创建时间
        ],[]....
    },
    "errCode":0
}
```

