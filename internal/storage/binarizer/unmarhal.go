package binarizer

import (
	"encoding/binary"
	"fmt"
)

func (b *Binarizer) Unmarshal(data []byte) (map[string]any, error) {
	const op = "binarizer.Unmarshal"
	log := b.logger.With().Str("op", op).Logger()

	// selected key for next pair
	var key string
	// result map
	m := make(map[string]any)

	// TODO: обновление
	// 1. сначала разделяем на анмаршал именно стрингов
	// 2. в стрингах уходим от большого количества буферов, а работаем чисто со слайсом
	//
	// NOTE: уходим от буфера, так как все равное нам от него нужен метод bytes.
	// а еще мы использовали только из за реализации интерфейса io.Writer
	//
	// а что будет с другими типами данных хз. посмотрим
	// TLV: type (1 byte) -> len (2 bytes)-> value (len bytes)

	for i := 0; i <= len(data); {
		if len(data) <= i {
			break
		}

		t := data[i]
		log.Debug().Msgf("got data type: %X", t)

		// NOTE: read next bytes: data len
		i++
		dataLen := binary.LittleEndian.Uint16(data[i : i+2])
		log.Debug().Msgf("got data len: %d", dataLen)

		// NOTE: read next bytes: skip 2 bytes -> data len
		i += 2

		var (
			unmarshalledData any
			err              error
		)
		switch t {
		case StrType:
			log.Debug().Msg("start unmarshal str")
			unmarshalledData = UnmarshalString(data[i : i+int(dataLen)])
		case IntType:
			log.Debug().Msg("start unmarshal int")
			unmarshalledData, err = UnmarshalInt(data[i : i+int(dataLen)])
			if err != nil {
				log.Err(err).Msg("failed to unmarshal int. skip pair")
				i += int(dataLen) // добавляем, так как перейдем к некст итерации
				key = ""          // сбрасываем ключ, если он был
				continue          // переходим к следующей паре
			}
		case FloatType:
			log.Debug().Msg("start unmarshal float")
			unmarshalledData, err = UnmarshalFloat(data[i : i+int(dataLen)])
			if err != nil {
				log.Err(err).Msg("failed to unmarshal float. skip pair")
				i += int(dataLen) // добавляем, так как перейдем к некст итерации
				key = ""          // сбрасываем ключ, если он был
				continue          // переходим к следующей паре
			}

		case BoolType:
			log.Debug().Msg("start unmarshal bool")
			unmarshalledData, err = UnmarshalBool(data[i : i+int(dataLen)])
			if err != nil {
				log.Err(err).Msg("failed to unmarshal bool. skip pair")
				i += int(dataLen) // добавляем, так как перейдем к некст итерации
				key = ""          // сбрасываем ключ, если он был
				continue          // переходим к следующей паре
			}

		default:
			log.Warn().Msgf("got unsupported type; returning nil")
			unmarshalledData = nil
		}

		// NOTE: go to next TVL
		i += int(dataLen)

		log.Debug().Msgf("unmarshalled data: %v", unmarshalledData)
		log.Debug().Msgf("idx: %d; len: %d", i, len(data))

		if key == "" {
			var ok bool
			key, ok = unmarshalledData.(string)
			if !ok {
				return nil, fmt.Errorf("%s: invalid key data type", op)
			}
		} else {
			m[key] = unmarshalledData
			key = ""
		}
	}

	log.Debug().Msgf("result map: %v", m)
	return m, nil
}
