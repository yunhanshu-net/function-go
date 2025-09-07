package reflectx

import (
	"fmt"
	"reflect"
	"strings"
)

// TypeInfo 包含类型和值的详细信息
type TypeInfo struct {
	TypeName   string
	Kind       string
	Value      interface{}
	Underlying string
}

// GetTypeInfo 获取变量的详细类型信息
func GetTypeInfo(v interface{}) TypeInfo {
	t := reflect.TypeOf(v)
	return TypeInfo{
		TypeName:   t.String(),
		Kind:       t.Kind().String(),
		Value:      v,
		Underlying: getUnderlyingType(t),
	}
}

// getUnderlyingType 获取底层类型信息
func getUnderlyingType(t reflect.Type) string {
	if t == nil {
		return "nil"
	}

	switch t.Kind() {
	case reflect.Ptr:
		return "*" + getUnderlyingType(t.Elem())
	case reflect.Slice:
		return "[]" + getUnderlyingType(t.Elem())
	case reflect.Array:
		return fmt.Sprintf("[%d]", t.Len()) + getUnderlyingType(t.Elem())
	case reflect.Map:
		return "map[" + getUnderlyingType(t.Key()) + "]" + getUnderlyingType(t.Elem())
	case reflect.Chan:
		var dir string
		switch t.ChanDir() {
		case reflect.RecvDir:
			dir = "<-chan "
		case reflect.SendDir:
			dir = "chan<- "
		default:
			dir = "chan "
		}
		return dir + getUnderlyingType(t.Elem())
	default:
		return t.Kind().String()
	}
}

// PrintTypeInfo 打印变量的类型信息
func PrintTypeInfo(v interface{}) {
	info := GetTypeInfo(v)
	fmt.Printf("变量类型信息:\n")
	fmt.Printf("  类型名称: %s\n", info.TypeName)
	fmt.Printf("  基础类型: %s\n", info.Kind)
	fmt.Printf("  底层类型: %s\n", info.Underlying)
	fmt.Printf("  值: %#v\n", info.Value)
}

// AssertType 断言变量是否为特定类型
func AssertType(v interface{}, expected interface{}) bool {
	expectedType := reflect.TypeOf(expected)
	actualType := reflect.TypeOf(v)
	return actualType == expectedType
}

// SafeCast 安全类型转换
func SafeCast(v interface{}, target interface{}) (interface{}, bool) {
	targetType := reflect.TypeOf(target)
	actualType := reflect.TypeOf(v)

	if actualType == targetType {
		return v, true
	}

	// 处理指针类型
	if actualType.Kind() == reflect.Ptr && actualType.Elem() == targetType {
		return reflect.ValueOf(v).Elem().Interface(), true
	}

	return nil, false
}

// TypeTree 打印变量的类型树
func TypeTree(v interface{}, indent string) string {
	t := reflect.TypeOf(v)
	if t == nil {
		return "nil"
	}

	var builder strings.Builder
	builder.WriteString(t.Kind().String())

	switch t.Kind() {
	case reflect.Ptr:
		builder.WriteString("\n" + indent + "└── ")
		builder.WriteString(TypeTree(reflect.New(t.Elem()).Elem().Interface(), indent+"    "))
	case reflect.Slice, reflect.Array:
		builder.WriteString("\n" + indent + "└── ")
		builder.WriteString(TypeTree(reflect.New(t.Elem()).Elem().Interface(), indent+"    "))
	case reflect.Map:
		builder.WriteString("\n" + indent + "├── key: ")
		builder.WriteString(TypeTree(reflect.New(t.Key()).Elem().Interface(), indent+"│   "))
		builder.WriteString("\n" + indent + "└── value: ")
		builder.WriteString(TypeTree(reflect.New(t.Elem()).Elem().Interface(), indent+"    "))
	case reflect.Struct:
		builder.WriteString(" (")
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if i > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(field.Name + ": " + field.Type.Kind().String())
		}
		builder.WriteString(")")
	}

	return builder.String()
}
