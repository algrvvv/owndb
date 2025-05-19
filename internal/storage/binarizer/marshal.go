package binarizer

import (
	"bytes"
	"fmt"
)

// Marshal метод для сериализации мапы.
func (b *Binarizer) Marshal(m map[string]any) ([]byte, error) {
	const op = "binarizer.Marshal"
	log := b.logger.With().Str("op", op).Logger()

	buf := new(bytes.Buffer)
	for key, value := range m {
		log.Debug().Str("key", key).Any("value", value).Msg("get new pair")
		dt, err := GetDataType(key, value)
		if err != nil {
			log.Err(err).Msg("failed to get data type. skip...")
			continue
		}

		log.Debug().Msgf("got value type: %x", dt)

		// сначала маршалим ключ
		log.Debug().Msgf("start marshal key string: %s", key)
		err = MarshalString(buf, key)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		// маршилим значение по его типу
		switch dt {
		case StrType:
			log.Debug().Msgf("start marshal string: %s", value)
			vstr, ok := value.(string)
			if !ok {
				panic("failed to convert expected string to string")
			}

			err = MarshalString(buf, vstr)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
		case IntType:
			log.Debug().Msgf("start marshal int: %d", value)
			intV, ok := value.(int)
			if !ok {
				return nil, &ErrTypeMismatch{key, value, "int"}
			}

			err = MarshalInt(buf, int64(intV))
			if err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}

		case FloatType:
			log.Debug().Msgf("start marshal float: %v", value)
			flV, ok := value.(float64)
			if !ok {
				return nil, &ErrTypeMismatch{key, value, "float64"}
			}

			err = MarshalFloat(buf, flV)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}

		case BoolType:
			log.Debug().Msgf("start marshal bool: %v", value)
			boolV, ok := value.(bool)
			if !ok {
				return nil, &ErrTypeMismatch{key, value, "bool"}
			}

			err = MarshalBool(buf, boolV)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}

		default:
			return nil, fmt.Errorf("%s: unsupported data type: %T", op, value)
		}
	}

	return buf.Bytes(), nil
}
