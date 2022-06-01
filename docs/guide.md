# 使用指南

## 配置项说明

| 名称      | 说明                          | 类型   | 默认值 | 例子                                                |
| --------- | ----------------------------- | ------ | ------ | --------------------------------------------------- |
| project   | 项目名                        | string | -      | blog                                                |
| package   | 项目包名                      | string | -      | github.com/ychengcloud/cre/examples/blog"                |
| dsn       | 项目名                        | string | -      | mysql://root:@tcp(127.0.0.1:3306)/blog?charset=utf8 |
| root      | 模板根路径,相对于当前路径     | string | -      | path                                                |
| genRoot   | 模板生成根路径,相对于当前路径 | string | -      | path                                                |
| templates | 模板                          | array  | -      | -                                                   |
| tables    | 表                            | array  | -      | -                                                   |


模板配置: templates

| 名称    | 说明                                | 类型   | 默认值 | 例子   |
| ------- | ----------------------------------- | ------ | ------ | ------ |
| path    | 模板相对路径,相对于Root             | string | -      | path   |
| genPath | 生成路径，相对于GenRoot             | string | -      | path   |
| format  | 生成文件名格式                      | string | _      | path   |
| mode    | 生成模式, 可选值: "single", "multi" | enum   | single | single |

- 数据表配置项 : tables

| 名称       | 说明             | 类型         | 默认值                                                  | 例子               |
| ---------- | ---------------- | ------------ | ------------------------------------------------------- | ------------------ |
| name       | 数据表名         | string       | -                                                       | product            |
| skip       | 是否生成相应代码 | bool         | false                                                   | false              |
| errorCodes | 错误码           | string array | -                                                       |
| methods    | 支持的 Api 方法  | string array | ["list", "update", "delete", "batchGet", "batchDelete"] | ["list", "update"] |
| fields     | 字段             | field array  | -                                                       | -                  |

- 字段配置项 : field:

| 名称       | 说明               | 类型   | 默认值 | 例子                      |
| ---------- | ------------------ | ------ | ------ | ------------------------- |
| name       | 字段名             | string | -      | id                        |
| alias      | 别名               | string | -      | nameAlias                 |
| skip       | 是否忽略此字段     | bool   | false  | true                      |
| required   | 是否必填字段       | bool   | false  | true                      |
| sortable   | 是否可按此字段排序 | bool   | false  | true                      |
| filterable | 是否可按此字段过滤 | bool   | false  | true                      |
| operations | 排序时的可用操作   | array  | -      | true                      |
| tags       | 扩展 struct tags   | string | ""     | binding:"required,max=64" |
| relation   | 关联配置项         | object | -      | -                         |


关联配置项: relation

| 名称       | 说明                                                       | 类型   | 默认值 | 例子       |
| ---------- | ---------------------------------------------------------- | ------ | ------ | ---------- |
| name       | 关联表名                                                   | string | -      | category   |
| type       | 关联类型,取值 None, BelongsTo, HasOne, HasMany, ManyToMany | enum   | None   | ManyToMany |
| ref_table  | 指定关联表表名                                             | string | ""     | category   |
| join_table | 连接表配置项                                               | object | -      | -          |

连接表配置项: join_table

| 名称 | 说明     | 类型   | 默认值 | 例子          |
| ---- | -------- | ------ | ------ | ------------- |
| name | 连接表名 | string | -      | post_category |

