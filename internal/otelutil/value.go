package otelutil

import (
	"fmt"
	"reflect"
	"time"

	otellog "go.opentelemetry.io/otel/log"
)

// AnyToValue converts common Go values into an OpenTelemetry log value.
func AnyToValue(value any) otellog.Value {
	switch v := value.(type) {
	case nil:
		return otellog.Value{}
	case otellog.Value:
		return v
	case string:
		return otellog.StringValue(v)
	case bool:
		return otellog.BoolValue(v)
	case int:
		return otellog.IntValue(v)
	case int8:
		return otellog.Int64Value(int64(v))
	case int16:
		return otellog.Int64Value(int64(v))
	case int32:
		return otellog.Int64Value(int64(v))
	case int64:
		return otellog.Int64Value(v)
	case uint:
		return otellog.Int64Value(int64(v))
	case uint8:
		return otellog.Int64Value(int64(v))
	case uint16:
		return otellog.Int64Value(int64(v))
	case uint32:
		return otellog.Int64Value(int64(v))
	case uint64:
		return otellog.StringValue(fmt.Sprintf("%d", v))
	case uintptr:
		return otellog.StringValue(fmt.Sprintf("%d", v))
	case float32:
		return otellog.Float64Value(float64(v))
	case float64:
		return otellog.Float64Value(v)
	case []byte:
		return otellog.BytesValue(v)
	case time.Time:
		return otellog.StringValue(v.Format(time.RFC3339Nano))
	case time.Duration:
		return otellog.StringValue(v.String())
	case error:
		return otellog.StringValue(v.Error())
	case fmt.Stringer:
		return otellog.StringValue(v.String())
	case []string:
		values := make([]otellog.Value, 0, len(v))
		for _, item := range v {
			values = append(values, otellog.StringValue(item))
		}
		return otellog.SliceValue(values...)
	case []int:
		values := make([]otellog.Value, 0, len(v))
		for _, item := range v {
			values = append(values, otellog.IntValue(item))
		}
		return otellog.SliceValue(values...)
	case []int64:
		values := make([]otellog.Value, 0, len(v))
		for _, item := range v {
			values = append(values, otellog.Int64Value(item))
		}
		return otellog.SliceValue(values...)
	case []float64:
		values := make([]otellog.Value, 0, len(v))
		for _, item := range v {
			values = append(values, otellog.Float64Value(item))
		}
		return otellog.SliceValue(values...)
	case []any:
		values := make([]otellog.Value, 0, len(v))
		for _, item := range v {
			values = append(values, AnyToValue(item))
		}
		return otellog.SliceValue(values...)
	case map[string]any:
		return MapToValue(v)
	default:
		return reflectedValue(reflect.ValueOf(value))
	}
}

// MapToAttributes converts a flat map into OTel attributes.
func MapToAttributes(fields map[string]any) []otellog.KeyValue {
	attrs := make([]otellog.KeyValue, 0, len(fields))
	for key, value := range fields {
		attrs = append(attrs, otellog.KeyValue{
			Key:   key,
			Value: AnyToValue(value),
		})
	}
	return attrs
}

// MapToValue converts a map to an OTel map value.
func MapToValue(fields map[string]any) otellog.Value {
	return otellog.MapValue(MapToAttributes(fields)...)
}

func reflectedValue(value reflect.Value) otellog.Value {
	if !value.IsValid() {
		return otellog.Value{}
	}

	switch value.Kind() {
	case reflect.Pointer:
		if value.IsNil() {
			return otellog.Value{}
		}
		return reflectedValue(value.Elem())
	case reflect.Slice, reflect.Array:
		if value.Type().Elem().Kind() == reflect.Uint8 {
			return otellog.BytesValue(value.Bytes())
		}
		values := make([]otellog.Value, 0, value.Len())
		for i := 0; i < value.Len(); i++ {
			values = append(values, AnyToValue(value.Index(i).Interface()))
		}
		return otellog.SliceValue(values...)
	case reflect.Map:
		if value.Type().Key().Kind() != reflect.String {
			return otellog.StringValue(fmt.Sprint(value.Interface()))
		}
		fields := make(map[string]any, value.Len())
		iter := value.MapRange()
		for iter.Next() {
			fields[iter.Key().String()] = iter.Value().Interface()
		}
		return MapToValue(fields)
	case reflect.Struct:
		return otellog.StringValue(fmt.Sprint(value.Interface()))
	default:
		return otellog.StringValue(fmt.Sprint(value.Interface()))
	}
}
