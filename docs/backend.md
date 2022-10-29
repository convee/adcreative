# 通用约定
## 1.接口域名
- 开发：https://adcreativedev.convee.cn/backend/
- 生产：https://adcreative.convee.cn/backend/

## 2.请求方式

- GET + POST
## 3.公共请求参数

- Header

| 参数 | 描述 | 类型 | 备注 |
| --- | --- | --- | --- |
| token | AM token | string | 必填 |

## 4.全局返回格式约定
```json
{
    "code": 200,
    "msg": "",
    "data": {
    }
}
```
## 5.分页约定

- 请求参数

| 参数名称 | 参数描述 | 类型 | 备注 |
| --- | --- | --- | --- |
| page | 当前页面 | int | 默认为1 |
| per_page | 分页数量 | int | 默认为20 |


- 返回参数说明

```json

{
  "code": 200,
  "msg": "ok",
  "data": {
    "lists": [
      {}
    ],
    "total": 1
  }
}
```
## 6. 系统错误码

| 错误码 | 含义 |
| --- | --- |
| 200 | 正常返回 |
| 500 | 通用错误码 |
| 400 | 参数错误 |

# 系统接口文档
### 0. 获取客户列表信息

- 接口地址: `customer/list`
- 请求方式: `GET`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| page | 当前页面 | 否 | int | 默认为1 |
| per_page | 分页数量 | 否 | int | 默认为20 |
| name | 客户名称 | 否 | string |  |
| id | 客户ID | 否 | int |  |
| is_public | 是否共有账号 | 否 | int | 默认为0 |

- 返回参数

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "data": [
          {
            "name": "客户",
            "is_public": "1",
            "token": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
          },
          {}
        ],
    "from": 1,
    "last_page": 1,
    "per_page": "20",
    "to": 4,
    "total": 4 
  }
}
```

### 1. 增加客户

- 接口地址: `customer/create`
- 请求方式: `POST`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| name | 客户名称 | 是 | string | 默认为1 |

- 返回参数

```json
{
    "code": 0,
    "msg": "",
    "data": {
      "id": 1,
      "name": "客户",
      "is_public": "1",
      "token": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }
}
```
### 
### 2. 编辑客户

- 接口地址: `customer/edit`
- 请求方式: `POST`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| id | 客户id | 是 | int |  |
| name | 客户名称 | 是 | string |

- 返回参数

```json
{

    "code": 0,
    "msg": "",
    "data": {
      "id": 1,
      "name": "客户",
      "is_public": "1",
      "token": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }
}
```

### 3. 删除客户

- 接口地址: `customer/delete`
- 请求方式: `GET`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| id | 客户id | 是 | int |

- 返回参数

```json
{

    "code": 0,
    "msg": "",
    "data": {
     
  }
}
```
### 
### 4. 媒体列表

- 接口地址: `publisher/list`
- 请求方式: `GET`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| page | 当前页面 | 否 | int | 默认为1 |
| per_page | 分页数量 | 否 | int | 默认为20 |
| name | 媒体名称 | 否 | string |  |
| id | 媒体ID | 否 | int |  |

- 返回参数

```json
{
    "code": 0,
    "msg": "",
    "data": {
      "data": [
            {
              "id": 1,
              "name": "腾讯",
              "info": "媒体方接口信息"
            },
            {}
          ],
      "from": 1,
      "last_page": 1,
      "per_page": "20",
      "to": 4,
      "total": 4 
  }
}

```
### 5. 增加媒体

- 接口地址: `publisher/create`
- 请求方式: `POST`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| name | 媒体名称 | 是 | string | 默认为1 |
| info | 送审地址相关信息 | 是 | string |  |

- 返回参数

```json
{
    "code": 0,
    "msg": "",
    "data": {
      "id": 1,
      "name": "点号通客户",
      "info": "xxxxxxxxxx"
  }
}
```
### 
### 6. 编辑媒体

- 接口地址: `publisher/edit`
- 请求方式: `POST`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| id | 媒体id | 是 | int |  |
| name | 媒体名称 | 是 | string |
| info | 送审地址相关信息 | 否 | string |  |

- 返回参数

```json
{

    "code": 0,
    "msg": "",
    "data": {
      "id": 1,
      "name": "媒体账户",
      "info": "xxxxxxx"
  }
}
```

### 7. 删除媒体

- 接口地址: `publisher/delete`
- 请求方式: `GET`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| id | 媒体id | 是 | int |

- 返回参数

```json
{

    "code": 0,
    "msg": "",
    "data": {
     
  }
}
```
### 
### 8. 广告位列表

- 接口地址: `position/list`
- 请求方式: `GET`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| page | 当前页面 | 否 | int | 默认为1 |
| per_page | 分页数量 | 否 | int | 默认为20 |
| publisher_id | 媒体id | 否 | int |  |
| name | 广告位名称 | 否 | string |  |
| position_name | 英文广告位名称 | 否 | string |  |

- 返回参数

```json
{
    "code": 0,
    "msg": "",
    "data": {
      "data": [
            {
              "id": "广告位Id",
              "publisher_id":"媒体Id",
              "name":"广告位名称",
              "type":"投放方式otv/display",
              "position":"广告位标识",
              "position_name":"position名称",
              "material_info":"素材详细信息",
              "media_type":"设备类型"
             
            },
            {}
          ],
      "from": 1,
      "last_page": 1,
      "per_page": "20",
      "to": 4,
      "total": 4 
  }
}

```
### 9. 增加广告位

- 接口地址: `position/create`
- 请求方式: `POST`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| publisher_id | 媒体id | 是 | int |
| name | 广告位名称 | 是 | string |  |
| type | 投放方式otv/display | 是 | int |  |
| position | 广告位标识 | 是 | string |  |
| position_name | position名称 | 是 | string |  |
| material_info | 素材详细信息 | 是 | string |  |
| media_type | 设备类型 | 否 | string |  |
| ad_format | 广告形式,OTV:1,OTT:2,DISPLAY:3 | 是 | int |  |
| material_is_url | 广告位对应不同活动类型素材在界面的上传格式，按位取值,1为url、0为上传素材 | 是 | int |  |
| is_rsync | 是否需要素材同步 0:否 1:是 | 是 | int |  |
| is_config_matter | 是否需要配置素材,字段含义 0:否 1:是 | 是 | int |  |
| landing_change_need_rsync | 落地页修改是否需要再次同步媒体审核 0:不需要，1:需要 | 是 | int |  |
| monitor_code_change_need_rsync | 监测代码(曝光/点击)修改是否需要再次同步媒体审核 0:不需要，1:需要 | 是 | int |  |
| monitor_position_change_need_rsync | 监测位置修改是否需要再次同步媒体审核 0:不需要，1:需要 | 是 | int |  |
| is_creative_bind | 是否创意绑定 0:否 1:是 | 是 | int |  |
| pv_limit | 曝光监测条数上限（不包括RM自己的监测），默认为0，则为不限 | 否 | string |  |
| cl_limit | 点击监测条数上限（不包括RM自己的监测），默认为0，则为不限 | 否 | string |  |

- 返回参数

```json
{
    "code": 0,
    "msg": "",
    "data": {
      
  }
}
```
### 
### 10. 编辑广告位

- 接口地址: `position/edit`
- 请求方式: `POST`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| id | 广告位id | 是 | int |  |
| publisher_id | 媒体id | 是 | int |
| name | 广告位名称 | 是 | string |  |
| type | 投放方式otv/display | 是 | int |  |
| position | 广告位标识 | 是 | string |  |
| position_name | position名称 | 是 | string |  |
| material_info | 素材详细信息 | 是 | string |  |
| media_type | 设备类型 OTT  MO  PC | 否 | string |  |
| ad_format | 广告形式,OTV:1,OTT:2,DISPLAY:3 | 是 | int |  |
| material_is_url | 广告位对应不同活动类型素材在界面的上传格式，按位取值,1为url、0为上传素材 | 是 | int |  |
| is_rsync | 是否需要素材同步 0:否 1:是 | 是 | int |  |
| is_config_matter | 是否需要配置素材,字段含义 0:否 1:是 | 是 | int |  |
| landing_change_need_rsync | 落地页修改是否需要再次同步媒体审核 0:不需要，1:需要 | 是 | int |  |
| monitor_code_change_need_rsync | 监测代码(曝光/点击)修改是否需要再次同步媒体审核 0:不需要，1:需要 | 是 | int |  |
| monitor_position_change_need_rsync | 监测位置修改是否需要再次同步媒体审核 0:不需要，1:需要 | 是 | int |  |
| is_creative_bind | 是否创意绑定 0:否 1:是 | 是 | int |  |
| pv_limit | 曝光监测条数上限（不包括RM自己的监测），默认为0，则为不限 | 否 | string |  |
| cl_limit | 点击监测条数上限（不包括RM自己的监测），默认为0，则为不限 | 否 | string |  |

- 返回参数

```json
{

    "code": 0,
    "msg": "",
    "data": {
    
  }
}
```

### 11. 删除广告位

- 接口地址: `position/delete`
- 请求方式: `GET`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| id | 广告位id | 是 | int |

- 返回参数

```json
{

    "code": 0,
    "msg": "",
    "data": {
     
  }
}
```
### 12. 送审方媒体账号列表

- 接口地址: `publisher_account/list`
- 请求方式: `GET`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| page | 当前页面 | 否 | int | 默认为1 |
| per_page | 分页数量 | 否 | int | 默认为20 |
| publisher_id | 媒体id | 否 | int |  |
| customer_id | 客户id | 否 | string |  |

- 返回参数

```json
{
    "code": 0,
    "msg": "",
    "data": {
      "data": [
            {
              "id":"account_id",
              "publisher_id":"媒体账户的ID",
              "dsp_id":"DSP的ID",
              "token":"DSP的令牌",
              "is_rsync_advertiser":"是否同步广告主 1/是,0/否",
              "is_rsync_creative":"是否同步创意 1/是,0/否",
              "customer_id":"客户ID",
              "remark":"备注",
            },
            {}
          ],
      "from": 1,
      "last_page": 1,
      "per_page": "20",
      "to": 4,
      "total": 4 
  }
}

```
### 13. 增加送审方媒体账号

- 接口地址: `publisher_account/create`
- 请求方式: `POST`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| publisher_id | 媒体id | 是 | int |
| dsp_id | DSP的ID | 是 | string |  |
| token | DSP的令牌 | 是 | string |  |
| is_rsync_advertiser | 是否同步广告主 1/是,0/否 | 是 | int |  |
| is_rsync_creative | 是否同步创意 1/是,0/否 | 是 | int |  |
| customer_id | 客户ID | 是 | string |  |
| remark | 备注 | 否 | string |  |

- 返回参数

```json
{
    "code": 0,
    "msg": "",
    "data": {
      
  }
}
```
### 
### 14. 编辑送审方媒体账号

- 接口地址: `publisher_account/edit`
- 请求方式: `POST`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| id | 账号id | 是 | int |  |
| publisher_id | 媒体id | 是 | int |
| dsp_id | DSP的ID | 是 | string |  |
| token | DSP的令牌 | 是 | string |  |
| is_rsync_advertiser | 是否同步广告主 1/是,0/否 | 是 | int |  |
| is_rsync_creative | 是否同步创意 1/是,0/否 | 是 | int |  |
| customer_id | 客户ID | 是 | string |  |
| remark | 备注 | 否 | string |  |

- 返回参数

```json
{

    "code": 0,
    "msg": "",
    "data": {
    
  }
}
```

### 15. 删除送审方媒体账号

- 接口地址: `publisher_account/delete`
- 请求方式: `GET`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| id | 广告位id | 是 | int |

- 返回参数

```json
{

    "code": 0,
    "msg": "",
    "data": {
     
  }
}
```
### 16. 广告主送审列表

- 接口地址: `advertiser/list`
- 请求方式: `GET`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| page | 当前页面 | 否 | int | 默认为1 |
| per_page | 分页数量 | 否 | int | 默认为20 |
| name | 广告主名称 | 否 | string |  |
| info | 广告主接口信息 | 否 | string |  |
| customer_id | 客户id | 否 | int |  |

- 返回参数

```json
{
    "code": 0,
    "msg": "",
    "data": {
      "data": [
            {
              "id":"id",
              "name":"广告主名称",
              "info":"广告主接口信息",
              "customer_id":"客户ID",
            },
            {}
          ],
      "from": 1,
      "last_page": 1,
      "per_page": "20",
      "to": 4,
      "total": 4 
  }
}

```
### 17. 增加广告主送审账号

- 接口地址: `advertiser/create`
- 请求方式: `POST`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| name | 广告主名称 | 是 | string |
| info | 广告主送审信息 | 是 | string |  |
| customer_id | 客户id | 是 | int |  |

- 返回参数

```json
{
    "code": 0,
    "msg": "",
    "data": {
      "id":"id",
      "name":"广告主名称",
      "info":"广告主接口信息",
      "customer_id":"客户ID",
  }
}
```
### 
### 18. 编辑广告主送审账号

- 接口地址: `advertiser/edit`
- 请求方式: `POST`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| id | 广告主id | 是 | int |  |
| name | 广告主名称 | 是 | string |
| info | 广告主送审信息 | 是 | string |  |
| customer_id | 客户id | 是 | int |  |

- 返回参数

```json
{

    "code": 0,
    "msg": "",
    "data": {
        "id":"id",
        "name":"广告主名称",
        "info":"广告主接口信息",
        "customer_id":"客户ID",
  }
}
```

### 19. 删除广告主送审账号

- 接口地址: `advertiser/delete`
- 请求方式: `GET`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| id | 广告主id | 是 | int |
  |

- 返回参数

```json
{

    "code": 0,
    "msg": "",
    "data": {
     
  }
}
```
### 
### 20. 广告主规则列表

- 接口地址: `advertiser_rules/list`
- 请求方式: `GET`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| page | 当前页面 | 否 | int | 默认为1 |
| per_page | 分页数量 | 否 | int | 默认为20 |
| publisher_id | 媒体id | 否 | int |  |

- 返回参数

```json
{
    "code": 0,
    "msg": "",
    "data": {
      "data": [
            {
              "id":"id",
              "publisher_id":"媒体id",
              "info":"广告主规则信息"
            },
            {}
          ],
      "from": 1,
      "last_page": 1,
      "per_page": "20",
      "to": 4,
      "total": 4 
  }
}

```
### 21. 增加广告主规则

- 接口地址: `advertiser_rules/create`
- 请求方式: `POST`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| publisher_id | 媒体id | 是 | int |
| info | 广告主送审信息 | 是 | string |  |

- 返回参数

```json
{
    "code": 0,
    "msg": "",
    "data": {
      "id":"id",
      "publisher_id":"媒体id",
      "info":"广告主规则信息",
  }
}
```
### 22. 编辑广告主规则账号

- 接口地址: `advertiser_rules/edit`
- 请求方式: `POST`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| id | 广告主规则id | 是 | int |  |
| publisher_id | 媒体id | 是 | int |
| info | 广告主规则信息 | 是 | string |  |

- 返回参数

```json
{

    "code": 0,
    "msg": "",
    "data": {
        "id":"id",
        "name":"广告主名称",
        "info":"广告主接口信息",
        "customer_id":"客户ID",
  }
}
```

### 23. 删除广告主规则账号

- 接口地址: `advertiser_rules/delete`
- 请求方式: `GET`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| id | 广告主id | 是 | int |
  |

- 返回参数

```json
{

    "code": 0,
    "msg": "",
    "data": {
     
  }
}
```
### 24. 广告主送审记录

- 接口地址: `advertiser_audit/list`
- 请求方式: `GET`
- 请求参数
| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| page | 当前页面 | 否 | int | 默认为1 |
| per_page | 分页数量 | 否 | int | 默认为20 |
| advertiser_id | 广告主id | 否 | int |  |

- 返回参数

```json
{
    "code": 0,
    "msg": "",
    "data": {
      "data": [
            {
              "id":"id",
              "advertiser_id":"广告主id",
              "publisher_id":"媒体id",
              "customer_id":"客户id",
              "status":"状态",
              "info":"广告主规则信息"
            },
            {}
          ],
      "from": 1,
      "last_page": 1,
      "per_page": "20",
      "to": 4,
      "total": 4 
  }
}

```

### 25. 创意送审记录

- 接口地址: `creative/list`
- 请求方式: `GET`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| page | 当前页面 | 否 | int | 默认为1 |
| per_page | 分页数量 | 否 | int | 默认为20 |
| position_id | 广告位id | 否 | int |  |

- 返回参数

```json
{
    "code": 0,
    "msg": "",
    "data": {
      "data": [ 
            {
              "id":"id",
              "position_id":"广告位id",
              "publisher_id":"媒体id",
              "creative_id":"创意id",
              "media_cid" : "媒体方创意ID",
              "name" : "送审方创意名称",
              "template_id" : "广告位模板ID",
              "info" : "素材信息",
              "land_url" : "落地页",
              "deeplink_url"  : "app唤起地址",
              "start_date" : "开始时间",
              "end_date" : "结束时间",
              "monitor" : "曝光监测",
              "cm" : "点击监测",
              "vm" : "可见性监测",
              "action" : "广告交互类型：1-打开网页 2-下载 3-deeplink",
              "mini_program_id" : "小程序ID",
              "mini_program_path" : "小程序路径",
              "customer_id" : "客户ID",
              "advertiser_id" : "广告主ID",
              "status" : "素材状态：0代上传，1待审核，2审核中，3审核不通过",
              "reason" : "审核失败原因"
            },
            {}
          ],
      "from": 1,
      "last_page": 1,
      "per_page": "20",
      "to": 4,
      "total": 4 
  }
}

```

