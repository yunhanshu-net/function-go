# Number 数字输入组件文档

## 组件概述
Number 组件用于数字输入场景，支持整数和小数输入，提供丰富的配置选项控制输入范围、精度和显示格式。

## 参数说明
| 参数名 | 说明 | 是否必填 | 默认值 | 示例 |
|--------|------|----------|--------|------|
| `min` | 允许的最小值 | 否 | 0 | `min:0` |
| `max` | 允许的最大值 | 否 | 999999 | `max:100` |
| `step` | 步进值 | 否 | 1 | `step:0.5` |
| `default_value` | 默认值 | 否 | - | `default_value:10` |
| `placeholder` | 占位提示文本 | 否 | - | `placeholder:请输入数量` |
| `unit` | 单位显示 | 否 | - | `unit:个` |
| `prefix` | 前缀符号 | 否 | - | `prefix:¥` |
| `suffix` | 后缀符号 | 否 | - | `suffix:%` |
| `precision` | 小数位数精度 | 否 | 0 | `precision:2` |
| `allow_negative` | 是否允许负数 | 否 | false | `allow_negative:true` |
| `allow_decimal` | 是否允许小数 | 否 | false | `allow_decimal:true` |

## 完整示例
```go
type ProductStockReq struct {
    // 库存数量配置
    Quantity int `json:"quantity" runner:"code:quantity;name:库存数量" widget:"type:number;min:0;max:1000;step:1;unit:个;placeholder:请输入库存数量" data:"type:number;default_value:100" validate:"required,min=0"`
    
    // 折扣率配置
    Discount float64 `json:"discount" runner:"code:discount;name:折扣率" widget:"type:number;min:0;max:1;step:0.01;precision:2;suffix:%" data:"type:number;default_value:0.9" validate:"required"`
}
```

## 使用场景
1. 数值范围输入（如库存数量、年龄等）
2. 精确数值输入（如金额、百分比等）
3. 带单位的数值输入（如重量、长度等）