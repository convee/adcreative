### 通用约定

#### 1. 接口域名

- https://adcreative.convee.cn/api

#### 2. 请求方式

- GET + POST

#### 3. API 参数编码

- 所有参数采用 UTF-8 编码
- POST 请求参数为 JSON 格式，必须将 Content-type 设置为：application/json

#### 4. 公共请求参数

- Header

| 参数名称      | 参数描述  | 类型 | 备注                                                                            |
|-----------|-------| --- |-------------------------------------------------------------------------------|
| appid     | 应用ID  | string | 表示应用ID，其中与之匹配的还有appsecret，表示应用密钥，用于数据的签名加密，不同的对接项目分配不同的appid和appsecret，保证数据安全 |
| timestamp | 请求时间  | int | 表示时间戳，当请求的时间戳与服务器中的时间戳，差值在5分钟之内，属于有效请求，不在此范围内，属于无效请求                          |
| nonce     | 临时流水号 | string | 表示临时流水号，用于防止重复提交验证                                                            |
| sign      | 接口签名  | string | 加密串 32 位 md5，生成规则见下文                                                          |

- sign 生成规则  `sign = md5(appid + secret + + nonce + timestamp)` (appid,secret 由素材服务线下提供)

```shell
sign = md5(10000dfcd79dbbd08a0203dbde9a9cc591598937600)
sign = 35cdbec3e6c25814d384221e6da21de2
```

注：素材服务会对 timestamp 与 sign 参数进行校验
(1) timestamp 与 素材服务 服务器时间 gap 超过 5min，则验证失败
(2) sign 值和 素材服务 服务器按以上规则生成的值不一致，则验证失败

#### 5. 全局返回格式约定

```json
{
  "code": 200,
  "msg": "",
  "data": {
  }
}
```

### 账号接口

#### 1. 送审方新增媒体账号

送审方拿到媒体账号提交到素材服务，支持更新

- 接口地址： publisher_account/add
- 请求方式： POST
- 请求参数

| **参数名称** | **参数描述** | **必填** | **类型** | **备注** |
| --- | --- | --- | --- | --- |
| publisher | 媒体标识 | 是 | string | 通过媒体列表获取 |
| dsp_id | DSP ID | 是 | string | 由媒体提供 |
| token | Token | 是 | string | 由媒体提供 |
| callback_url | 送审素材回调 | 否 | string | 素材服务送审媒体后回调送审方接口 |

- 返回参数

```json
{
   "code": 200,
   "msg": "success",
   "data": {}
}
```

### 广告主列表

#### 1. 广告主列表

获取可用的广告主（如需新增广告主请联系素材服务进行添加）

- 接口地址： advertiser/list
- 请求方式： GET
- 请求参数：

| **参数名称** | **参数描述** | **必填** | **类型** | **备注** |
| --- | --- | --- | --- | --- |
| advertiser_id | 广告主ID | 否 | int | 通过ID搜索广告主 |

- 返回参数

```json
{
   "code": 200,
   "msg": "success",
   "data": [
      {
         "id": 1,
         "name": "广告主测试",
         "publishers": [
            {
               "id": 6,
               "info": {
                  "advertiser_name": "",
                  "qualifications": [],
                  "industrys": []
               },
               "status": 0,
               "reason": "审核失败原因"
            },
            {
               "id": 101,
               "info": {
                  "advertiser_name": "",
                  "qualifications": [
                     {
                        "file_name": "",
                        "file_url": ""
                     }
                  ],
                  "industry": []
               },
               "status": 1,
               "reason": "审核失败原因"
            }
         ]
      }
   ]
}
```

####      

#### 2. 媒体行业ID

获取可用的广告主（如需新增广告主请联系素材服务进行添加）

- 接口地址： industry/list
- 请求方式： GET
- 请求参数：

| **参数名称** | **参数描述** | **必填** | **类型** | **备注** |
| --- | --- | --- | --- | --- |
| publisher | 媒体标识 | 是 | string |

- 返回参数

```json
{
   "code": 200,
   "msg": "ok",
   "data": {
      "lists": [
         {
            "id": "100000000",
            "value": "游戏"
         },
         {
            "id": "200000000",
            "value": "网电"
         },
         {
            "id": "300000000",
            "value": "其他"
         },
         {
            "id": "400000000",
            "value": "品牌"
         },
         {
            "id": "600000000",
            "value": "金融"
         },
         {
            "id": "700000000",
            "value": "教育"
         },
         {
            "id": "800000000",
            "value": "医美"
         },
         {
            "id": "900000000",
            "value": "房产"
         }
      ]
   }
}
```

####      

#### 3. 广告主送审规则列表

获取每个媒体的广告主送审规则

- 接口地址： advertiser/rules
- 请求方式： GET
- 请求参数：

| **参数名称** | **参数描述** | **必填** | **类型** | **备注** |
| --- | --- | --- | --- | --- |
| publisher | 媒体标识 | 否 | string | |

- 返回参数

```json
{
   "code": 200,
   "msg": "ok",
   "data": {
      "lists": [
         {
            "publisher_id": 2,
            "rules": [
               {
                  "desc": "广告主名称",
                  "format": [
                     "text"
                  ],
                  "key": "advertiser_name",
                  "limit": 0,
                  "measure": [
                     1,
                     100
                  ],
                  "required": 1,
                  "size": 0
               },
               {
                  "desc": "行业",
                  "format": "text",
                  "key": "industry",
                  "limit": 0,
                  "measure": [],
                  "required": 1,
                  "size": 0
               },
               {
                  "desc": "资质文件",
                  "format": [
                     "zip",
                     "rar",
                     "jpg",
                     "jpeg",
                     "png",
                     "bmp"
                  ],
                  "key": "qualifications",
                  "limit": 1,
                  "measure": [],
                  "required": 1,
                  "size": 204800
               }
            ]
         }
      ],
      "total": 9
   }
}
```

####      

#### 4. 广告主送审

按照各个媒体广告主送审规则将广告主信息送审到素材服务

- 接口地址： advertiser_audit/upload
- 请求方式： POST
- 请求参数：

| **参数名称** | **参数描述** | **必填** | **类型** | **备注** |
| --- | --- | --- | --- | --- |
| publisher | 媒体标识 | 是 | string | |
| advertiser_name | 广告主 | 是 | string | |
| advertiser_audit_info | 广告主信息 | 是 | object | |

**advertiser_info 参数说明：**

| **参数名称** | **参数描述** | **必填** | **类型** | **备注** |
| --- | --- | --- | --- | --- |
| company_name | 公司名称 | 否 | string |  |
| company_summary | 公司简介 | 否 | string |  |
| website_name | 网站名称 | 否 | string |  |
| website_address | 网站地址 | 否 | string |  |
| website_number | 网站icp备案号 | 否 | string |  |
| business_licenser | 营业执照注册号 | 否 | string |  |
| authorize_state | 代理授权书 | 否 | string |  |
| industry | 行业ID | 是 | array of string |  |
| qualifications | 广告主资质 | 是 | array of objects |  |

**请求示例：**

```json
{
   "advertiser_name": "广告主",
   "publisher": "IQIYI",
   "advertiser_info": {
      "advertiser_name": "广告主名称",
      "industry": "123",
      "qualifications": [
         {
            "file_name": "",
            "file_url": ""
         }
      ]
   }
}

```

**返回示例：**

```json
{
   "code": 200,
   "msg": "ok",
   "data": [
      {
         "publisher": "iqiyi",
         "advertiser_id": 13,
         "err_code": 0,
         "err_msg": ""
      }
   ]
}

```

#### 5. 广告主送审状态查询

广告主送审状态查询

- 接口地址： advertiser_audit/query
- 请求方式： GET
- 请求参数：

| **参数名称** | **参数描述** | **必填** | **类型** | **备注** |
| --- | --- | --- | --- | --- |
| advertiser_id | 广告主ID | 否 | int |
  |
| publisher | 媒体标识 | 否 | string | |

- 返回参数

```json
{
   "code": 200,
   "msg": "ok",
   "data": [
      {
         "advertiser_id": 17,
         "publisher": "IQIYI",
         "reason": "",
         "status": 0
      }
   ]
}
```

###      

### 广告位接口

#### 1. 获取媒体列表

获取已经支持的媒体列表

- 接口地址： publisher/list
- 请求方式: GET
- 请求参数：
- 返回参数：

```json
{
   "code": 200,
   "msg": "success",
   "data": [
      {
         "id": 1,
         "name": "Tencent"
      },
      {
         "id": 2,
         "name": "IQIYI"
      }
   ]
}
```

#### 2. 获取广告位信息

通过该接口获得广告信息，包括广告位总量，每个广告位的 ID、名称、尺寸以及支持的物料格式。

- 接口地址： position/list
- 请求方式: GET
- 请求参数：

| **参数名称** | **参数描述** | **必填** | **类型** | **备注** |
| --- | --- | --- | --- | --- |
| publisher | 媒体标识 | 是 | string | 通过媒体查询广告位列表 |
| position | 广告位标识 | 否 | string |

- 返回示例：

```json
{
   "code": 200,
   "msg": "success",
   "data": [
      {
         "id": 7,
         "position": "Tencent_Mobile_v1_xxx",
         "templates": [
            {
               "template_id": "1",
               "template_name": "西瓜信息流大图落地页",
               "info": [
                  {
                     "attr_name": "image",
                     "attr_desc": "图片1",
                     "format": [
                        "jpg",
                        "jpeg"
                     ],
                     "measure": [
                        "690x388",
                        "1280x720"
                     ],
                     "size": 100,
                     "required": 1
                  },
                  {
                     "attr_name": "video",
                     "attr_desc": "视频1",
                     "format": [
                        "mp4"
                     ],
                     "measure": [
                        "690x388",
                        "1280x720"
                     ],
                     "size": 100,
                     "required": 1
                  }
               ]
            },
            {
               "template_id": "2",
               "template_name": "西瓜信息流落地页视频",
               "info": [
                  {
                     "attr_name": "image",
                     "attr_desc": "图片1",
                     "format": [
                        "jpg"
                     ],
                     "measure": [
                        "690x388",
                        "1280x720"
                     ],
                     "size": 100,
                     "required": 1
                  },
                  {
                     "attr_name": "title",
                     "attr_desc": "标题",
                     "format": [
                        "txt"
                     ],
                     "measure": [
                        "2",
                        "8"
                     ],
                     "size": 0,
                     "required": 1
                  }
               ]
            }
         ]
      }
   ]
}
```

### 素材接口

#### 1. 素材规格校验接口

- 接口地址： `creative/check`
- 请求方式： `POST`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| position | 广告位 | 是 | string |
  |
| template_id | 模板ID | 是 | string | |
| info | 素材信息 | 是 | array of objects | |

info 结构

| **参数名称** | **参数描述** | **必填** | **类型** | **备注** |
| --- | --- | --- | --- | --- |
| attr_name | 素材名称 | 是 | string | 广告位信息中所需要的字段名称 |
| attr_value | 素材内容 | 是 | string | 广告位信息中所需要的字段内容 |
| md5 | 素材内容md5值 | 否 | string | 图片、视频时填充文件md5值 |
| width | 素材宽 | 否 | int | 图片、视频时必填 |
| height | 素材高 | 否 | int | 图片、视频时必填 |
| ext | 素材格式 | 否 | string | 图片、视频必填 |
| duration | 素材时长 | 否 | int | 视频时必填 |

请求示例：

```json
{
   "position": "Tencent_Mobile_TXXW_Shanping_picture_logo",
   "template_id": "1",
   "info": [
      {
         "attr_name": "abstract",
         "attr_value": "摘要"
      },
      {
         "attr_name": "image",
         "attr_value": "http://v.convee.cn/image/tencent/2021/09/1632388009336762.jpg"
      }
   ]
}
```

返回示例：

```json
{
   "code": 200,
   "msg": "ok",
   "data": {
      "err_code": 1,
      "err_msg": [
         "标题必填"
      ],
      "position": "Tencent_Mobile_TXXW_Shanping_picture_logo"
   }
}
```

#### 2. 素材上传接口

- 接口地址： `creative/upload`
- 请求方式： `POST`
- 请求参数

| 参数名称 | 参数描述 | 必填 | 类型 | 备注 |
| --- | --- | --- | --- | --- |
| material | 物料信息 | 是 | array of objects | 物料信息 |

**material 参数说明**

| **参数名称** | **参数描述** | **必填** | **类型** | **备注** |
| --- | --- | --- | --- | --- |
| position | 广告位标识 | 是 | string | 从广告位列表接口获取 |
| advertiser_id | 广告主ID | 否 | int | 广告主送审返回 |
| publisher | 媒体标识 | 是 | string | 从媒体列表接口获取 |
| creative_id | 送审方创意ID | 是 | string | 送审方的创意 ID |
| media_info | 媒体方所需信息 | 否 | string | 腾讯广告主名称、微博行业ID、b612媒体dealID |
| name | 送审方创意名称 | 是 | string | 送审方的创意名称 |
| template_id | 模板ID | 是 | string | 广告位信息中的展现形式ID |
| info | 素材信息 | 是 | array of objects | 素材信息，根据广告位列表中的素材要求传递 |
| land_url | 落地页 | 是 | string | 广告落地页(广告点击后跳转的地址) |
| deeplink_url | 应用直达URL | 否 | string | 应用直达URL，当返回了deeplinkurl，优先唤醒本地app，如果无法唤醒，则调用land_url(打开或者下载) |
| start_date | 生效日期 | 是 | string | 素材的生效时间，格式为Y-m-d，例如2020-09-01 |
| end_date | 失效日期 | 是 | string | 素材的失效时间，格式为Y-m-d，例如2020-09-30 |
| monitor | 曝光监测 | 是 | array of objects | 用于填写第三方曝光监测地址，由url和t组成，url是监测地址，t是监测点位 t=0时，为首帧监测，t=1时，为第2秒监测，以此类推。注意送审方自身曝光监测地址必须放该数组第一位。 |
| cm | 点击监测 | 否 | array of string | 用于填写第三方点击监测地址。 注意， 送审方自身点击监测地址必须放该数组第一位 |
| vm | 可见性监测 | 否 | array of string | 用于填写第三方可见性监测地址。 注意， 送审方自身可见性监测地址必须放该数组第一位 |
| action | 广告交互类型 | 是 | int | 1-打开网页 2-下载 3-deeplink |
| mini_program_id | 小程序id | 否 | string | 小程序的ID |
| mini_program_path | 小程序路径 | 否 | string | 小程序的路径 |

**备注** 如果有多条监测的时候，请将自身监测放在最后一条。

**info 参数说明**

| **参数名称** | **参数描述** | **必填** | **类型** | **备注** |
| --- | --- | --- | --- | --- |
| attr_name | 素材名称 | 是 | string | 广告位信息中所需要的字段名称 |
| attr_value | 素材内容 | 是 | string | 广告位信息中所需要的字段内容 |
| md5 | 素材内容md5值 | 是 | string | 图片、视频时填充文件md5值 |
| width | 素材宽 | 否 | int | 图片、视频时必填 |
| height | 素材高 | 否 | int | 图片、视频时必填 |
| ext | 素材格式 | 否 | string | 图片、视频必填 |
| duration | 素材时长 | 否 | int | 视频时必填 |

- 请求示例

```json
{
   "is_only_valid": 1,
   "material": [
      {
         "position": "Tencent_mobile_v1_xxx",
         "template_id": "模板id",
         "creative_id": "10001",
         "name": "测试物料",
         "info": [
            {
               "attr_name": "image",
               "attr_value": "http://mat.convee.cn/image/sina/2020/09/851599640038003647.jpg",
               "md5": "d41d8cd98f00b204e9800998ecf8427e",
               "width": 1080,
               "height": 720,
               "ext": "jpg",
               "duration": 0
            },
            {
               "attr_name": "video",
               "attr_value": "http://mat.convee.cn/image/sina/2020/09/851599640038003647.mp4",
               "md5": "d41d8cd98f00b204e9800998ecf8427e",
               "width": 1080,
               "height": 720,
               "ext": "mp4",
               "duration": 15
            },
            {
               "attr_name": "title",
               "attr_value": "标题",
               "md5": "d41d8cd98f00b204e9800998ecf8427e",
               "width": 0,
               "height": 0,
               "ext": "txt",
               "duration": 0
            }
         ],
         "start_date": "2015-05-07",
         "end_date": "2018-08-30",
         "land_url": "http://tv.sohu.com/20150505/n412440832.shtml",
         "action": 1,
         "monitor": [
            {
               "t": 0,
               "url": "http://g.cn.miaozhen.com/x/k=2006958&p=6wxzh&dx=0&rt=2&ns=__IP__&ni=__IESID__&v=__LOC__&nd=__DRA__&np=__POS__&nn=__APP__&o="
            }
         ],
         "cm": [
            "http://e.cn.miaozhen.com/r/k=2012716&p=6wn6u&dx=0&rt=2&ns=__IP__&ni=__IESID__&v=__LOC__&nd=__DRA__&np=__POS__&nn=__APP__&o="
         ],
         "vm": [
            "http://e.cn.miaozhen.com/r/k=2012716&p=6wn6u&dx=0&rt=2&ns=__IP__&ni=__IESID__&v=__LOC__&nd=__DRA__&np=__POS__&nn=__APP__&o="
         ],
         "deeplink_url": "",
         "mini_program_id": "小程序id",
         "mini_program_path": "小程序路径"
      }
   ]
}
```

- 返回成功示例

```json
{
   "code": 200,
   "msg": "success",
   "data": [
      {
         "position": "Tencent_mobile_v1",
         "creative_id": "88n0f2202c3ddc54",
         "err_code": 1,
         "media_cid": "xxxx",
         "err_msg": "image1 素材尺寸为 【1242x2208】, 正确尺寸应为【720*1280】误差最大正负 5"
      }
   ]
}
```

- 返回失败示例

```json
{
   "code": 400,
   "msg": "请求参数错误",
   "data": [
      "Material[0].PublisherId:PublisherId不存在",
      "Material[1].PublisherId:PublisherId不存在",
      "Material[1].PositionId:PositionId为必填字段"
   ]
}
```

#### 3. 获取指定物料的审核结果

通过该 API 获取 请求中指定的物料审核的结果，最多查询20个创意

- 接口地址： creative/query
- 请求方式: POST
- 请求参数

| **参数名称** | **参数描述** | **必填** | **类型** | **备注** |
| --- | --- | --- | --- | --- |
| creative_id | 送审方创意ID | 是 | string | 送审方素材ID |

- 请求示例

```json
{
   "creative_id": [
      "a",
      "b"
   ]
}
```

- 返回示例

```json
{
   "code": 200,
   "msg": "success",
   "data": [
      {
         "creative_id": "10002",
         "media_cid": "888888",
         "status": 1,
         "reason": ""
      }
   ]
}
```

#### 4. 审核回调接口

通过该 API 将物料审核的结果回调给送审方，只调一次，不论成功失败

- 接口地址：由送审方提供
- 请求方式：POST
- 请求参数：

| **参数名称** | **参数描述** | **必填** | **类型** | **备注** |
| --- | --- | --- | --- | --- |
| material | 物料信息 | 是 | array of objects | 物料信息，支持创意批量回调 |

**material 参数说明：**

| **参数名称** | **参数描述** | **必填** | **类型** | **备注** |
| --- | --- | --- | --- | --- |
| creative_id | 送审方创意ID | 是 | string |
| status | 状态码 | 是 | int | 状态码 1初始化 2待送审 3审核中 4审核通过 5审核不通过 6发送媒体失败 |
| reason | 审核拒绝原因 | 否 | string |  |

- 返回示例

```json
{
   "code": 200,
   "msg": "success",
   "data": {}
}
```
