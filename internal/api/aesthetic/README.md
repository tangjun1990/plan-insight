# 审美感知服务 API

## 简介

审美感知服务提供了一组API，用于管理用户审美数据的采集、查询和分析。主要包括两部分：小程序端API和管理后台API。

## API 列表

### 小程序端 API

#### 1. 微信小程序用户鉴权
- **URL:** `/api/wx/auth`
- **方法:** `POST`
- **描述:** 基于微信小程序授权的手机号进行鉴权，生成用户token
- **请求参数:**
  ```json
  {
    "code": "微信临时登录凭证",
    "encryptedData": "加密数据",
    "iv": "加密算法初始向量",
    "phone": "手机号码"
  }
  ```
- **响应:**
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "token": "用户token",
      "expires_in": 604800,
      "user_info": {
        "id": 1,
        "phone": "13800138000"
      }
    }
  }
  ```

#### 2. 保存审美数据
- **URL:** `/api/aesthetic/data`
- **方法:** `POST`
- **描述:** 用户提交审美数据表单
- **请求头:** `Authorization: Bearer {token}`
- **请求参数:**
  ```json
  {
    "name": "姓名",
    "gender": "性别",
    "age": 25,
    "city": "城市",
    "phone": "手机号码",
    "liked_colors": ["红色", "蓝色", "黄色"],
    "disliked_colors": ["灰色", "黑色"],
    "liked_adjectives": ["温暖", "明亮", "舒适"],
    "liked_images": ["image1.jpg", "image2.jpg"]
  }
  ```
- **响应:**
  ```json
  {
    "code": 0,
    "message": "success"
  }
  ```

#### 3. 获取用户审美数据列表
- **URL:** `/api/aesthetic/data/list`
- **方法:** `GET`
- **描述:** 小程序用户查看自己提交的审美数据列表
- **请求头:** `Authorization: Bearer {token}`
- **请求参数:**
  - `page`: 页码，默认1
  - `page_size`: 每页条数，默认10
- **响应:**
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "list": [
        {
          "id": 1,
          "user_id": 1,
          "name": "姓名",
          "gender": "性别",
          "age": 25,
          "city": "城市",
          "phone": "手机号码",
          "liked_colors": "[\"红色\",\"蓝色\",\"黄色\"]",
          "disliked_colors": "[\"灰色\",\"黑色\"]",
          "liked_adjectives": "[\"温暖\",\"明亮\",\"舒适\"]",
          "liked_images": "[\"image1.jpg\",\"image2.jpg\"]",
          "created_at": "2023-01-01T12:00:00Z",
          "updated_at": "2023-01-01T12:00:00Z"
        }
      ],
      "total": 1,
      "page": 1,
      "page_size": 10
    }
  }
  ```

#### 4. 获取审美数据详情
- **URL:** `/api/aesthetic/data/{id}`
- **方法:** `GET`
- **描述:** 小程序用户查看自己提交的审美数据详情
- **请求头:** `Authorization: Bearer {token}`
- **响应:**
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "id": 1,
      "user_id": 1,
      "name": "姓名",
      "gender": "性别",
      "age": 25,
      "city": "城市",
      "phone": "手机号码",
      "liked_colors": "[\"红色\",\"蓝色\",\"黄色\"]",
      "disliked_colors": "[\"灰色\",\"黑色\"]",
      "liked_adjectives": "[\"温暖\",\"明亮\",\"舒适\"]",
      "liked_images": "[\"image1.jpg\",\"image2.jpg\"]",
      "created_at": "2023-01-01T12:00:00Z",
      "updated_at": "2023-01-01T12:00:00Z"
    }
  }
  ```

### 管理后台 API

#### 1. 管理员登录
- **URL:** `/admin/auth/login`
- **方法:** `POST`
- **描述:** 基于固定的手机号和密码进行登录验证
- **请求参数:**
  ```json
  {
    "phone": "管理员手机号",
    "password": "密码"
  }
  ```
- **响应:**
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "token": "管理员token",
      "expires_in": 86400
    }
  }
  ```

#### 2. 获取用户列表
- **URL:** `/admin/user/list`
- **方法:** `GET`
- **描述:** 获取所有用户数据列表，分页返回
- **请求头:** `Authorization: Bearer {token}`
- **请求参数:**
  - `page`: 页码，默认1
  - `page_size`: 每页条数，默认10
  - `phone`: 手机号过滤
  - `status`: 状态过滤，1正常，0禁用
- **响应:**
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "list": [
        {
          "id": 1,
          "name": "姓名",
          "phone": "手机号码",
          "gender": "性别",
          "age": 25,
          "city": "城市",
          "status": 1,
          "first_login_time": "2023-01-01T12:00:00Z",
          "last_login_time": "2023-01-01T12:00:00Z",
          "created_at": "2023-01-01T12:00:00Z",
          "updated_at": "2023-01-01T12:00:00Z"
        }
      ],
      "total": 1,
      "page": 1,
      "page_size": 10
    }
  }
  ```

#### 3. 禁用用户
- **URL:** `/admin/user/{id}/disable`
- **方法:** `PUT`
- **描述:** 将单个用户改为禁用状态
- **请求头:** `Authorization: Bearer {token}`
- **响应:**
  ```json
  {
    "code": 0,
    "message": "success"
  }
  ```

#### 4. 启用用户
- **URL:** `/admin/user/{id}/enable`
- **方法:** `PUT`
- **描述:** 将单个用户改为启用状态
- **请求头:** `Authorization: Bearer {token}`
- **响应:**
  ```json
  {
    "code": 0,
    "message": "success"
  }
  ```

#### 5. 获取审美数据列表
- **URL:** `/admin/aesthetic/data/list`
- **方法:** `GET`
- **描述:** 分页获取审美数据表中的数据
- **请求头:** `Authorization: Bearer {token}`
- **请求参数:**
  - `page`: 页码，默认1
  - `page_size`: 每页条数，默认10
  - `name`: 姓名过滤
  - `gender`: 性别过滤
  - `age_min`: 最小年龄过滤
  - `age_max`: 最大年龄过滤
  - `city`: 所在城市过滤
  - `phone`: 手机号过滤
- **响应:**
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "list": [
        {
          "id": 1,
          "user_id": 1,
          "name": "姓名",
          "gender": "性别",
          "age": 25,
          "city": "城市",
          "phone": "手机号码",
          "liked_colors": "[\"红色\",\"蓝色\",\"黄色\"]",
          "disliked_colors": "[\"灰色\",\"黑色\"]",
          "liked_adjectives": "[\"温暖\",\"明亮\",\"舒适\"]",
          "liked_images": "[\"image1.jpg\",\"image2.jpg\"]",
          "created_at": "2023-01-01T12:00:00Z",
          "updated_at": "2023-01-01T12:00:00Z"
        }
      ],
      "total": 1,
      "page": 1,
      "page_size": 10
    }
  }
  ```

#### 6. 获取审美数据统计分析
- **URL:** `/admin/aesthetic/data/analysis`
- **方法:** `GET`
- **描述:** 基于审美数据表中的数据，进行数据统计和分析
- **请求头:** `Authorization: Bearer {token}`
- **请求参数:**
  - `analysis_type`: 分析类型，可选值:
    - `color`: 喜欢的颜色分析
    - `disliked_color`: 讨厌的颜色分析
    - `adjective`: 喜欢的形容词分析
    - `image`: 喜欢的图片分析
    - `region`: 用户地域分布分析
  - `dimension`: 分析维度，可选值:
    - `count`: 按数量统计
    - `top`: 取前N项
    - `percent`: 按百分比统计
    - `map`: 地图数据格式（仅用于region类型）
  - `top`: 取前N条数据，默认10
  - `gender`: 按性别过滤
  - `age_min`: 最小年龄过滤
  - `age_max`: 最大年龄过滤
  - `city`: 按城市过滤
- **响应:**
  ```json
  {
    "code": 0,
    "message": "success",
    "data": [
      {
        "name": "红色",
        "count": 10,
        "percent": 25.0
      },
      {
        "name": "蓝色",
        "count": 8,
        "percent": 20.0
      }
    ]
  }
  ```

## 数据库设计

服务使用了以下三个主要的数据表：

### 1. 用户表 (users)
存储用户基本信息、登录状态和鉴权信息。

### 2. 审美数据表 (aesthetic_data)
存储用户提交的审美数据，包括喜欢的颜色、讨厌的颜色、喜欢的形容词和喜欢的图片。

### 3. 管理员表 (admins)
存储管理后台的管理员账号信息和鉴权信息。

## 部署说明

1. 确保MySQL数据库已正确配置
2. 启动服务，系统会自动创建表结构并初始化一个管理员账号
3. 默认管理员账号：
   - 手机号：13800138000
   - 密码：admin123