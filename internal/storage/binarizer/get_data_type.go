package binarizer

func GetDataType(key string, data any) (byte, error) {
	switch data.(type) {
	case string:
		return StrType, nil
	case int, int64, int32, int8:
		return IntType, nil
	case float64:
		return FloatType, nil
	case bool:
		return BoolType, nil
	default:
		return 0, &ErrInvalidDataType{key}
	}
}
