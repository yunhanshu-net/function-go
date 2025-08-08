# Number 组件规则说明

## 标签参数详解（按标签类型分组）

### widget 标签参数
- **type:number** (必需)：标识为数字输入组件，示例：`widget:"type:number"`
- **min** (可选)：允许输入的最小值，示例：`min:0`
- **max** (可选)：允许输入的最大值，示例：`max:100`
- **step** (可选)：步进值，示例：`step:0.5`
- **placeholder** (可选)：提示文本，示例：`placeholder:请输入年龄`
- **unit** (可选)：单位文本，示例：`unit:岁`
- **prefix** (可选)：前缀符号，示例：`prefix:¥`
- **suffix** (可选)：后缀符号，示例：`suffix:%`
- **precision** (可选)：小数位数精度，示例：`precision:2`
- **allow_negative** (可选)：是否允许负数，示例：`allow_negative:true`
- **allow_decimal** (可选)：是否允许小数，示例：`allow_decimal:true`

### data 标签参数
- **type:number** (必需)：标识字段为数字类型，示例：`data:"type:number"`
- **default_value** (可选)：设置默认值，示例：`default_value:100`
- **example** (可选)：提供示例值，示例：`example:50`

### validate 标签参数
- **min** (可选)：最小值验证，示例：`min=0`
- **max** (可选)：最大值验证，示例：`max=100`
- **required** (可选)：必填验证，示例：`required`

## 标签规则
1. **widget 标签**：必须包含 `type:number`，可选参数如上
2. **data 标签**：必须包含 `type:number`，可包含 `default_value` 和 `example`
3. **validate 标签**：可包含数值验证规则

## 正例（完整结构体定义）
```go
// 商品价格配置
type ProductPriceConfig struct {
    Price float64 `json:"price" runner:"code:price;name:商品价格" widget:"type:number;min:0;max:9999;step:0.01;precision:2;prefix:¥" data:"type:number;default_value:99.9;example:199.99" validate:"required,min=0"`
}

// 年龄输入配置
type UserProfileConfig struct {
    Age int `json:"age" runner:"code:age;name:年龄" widget:"type:number;min:1;max:120;unit:岁" data:"type:number;example:30" validate:"required"`
}
```

## 反例（完整结构体定义）
```go
// 无效的example类型
type DiscountConfig struct {
    Discount float64 `json:"discount" runner:"code:discount;name:折扣" widget:"type:number" data:"type:number;example:abc" validate:"required"`
}

// example超出范围
type TemperatureConfig struct {
    Temp float64 `json:"temp" runner:"code:temp;name:温度" widget:"type:number;min:-10;max:50" data:"type:number;example:100" validate:"required"`
}
```

## 注意事项
1. **标签格式**：
    - 所有标签必须写在同一行
    - 多个参数用分号分隔
    - 参数值使用冒号分隔

2. **参数分类**：
    - widget：控制UI显示和行为
    - data：定义数据类型和值
    - validate：设置验证规则

3. **数值约束**：
    - min 必须小于 max
    - step 值应合理
    - precision 设置后输入会自动四舍五入

4. **默认值和示例**：
    - default_value 是实际使用的默认值
    - example 仅用于文档和提示
    - 两者都必须是有效的数值