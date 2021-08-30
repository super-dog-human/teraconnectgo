package domain

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

func UUIDWithoutHypen() (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return strings.Replace(uuid.String(), "-", "", -1), nil
}

func MergeJsonToStruct(jsonDiff *map[string]interface{}, origin interface{}, allowFields *[]string) {
	allKeys := TopLevelStructKeys(origin)

	o := reflect.ValueOf(origin).Elem()
	if !o.CanSet() {
		return
	}

	for rawName, jsonValue := range *jsonDiff {
		fieldName := strings.Title(rawName) // jsonのフィールド名を比較用にアッパーキャメルにする
		if !Contains(&allKeys, fieldName) || !Contains(allowFields, fieldName) {
			continue
		}

		targetField := o.FieldByName(fieldName)

		if jsonValue == nil {
			targetField.Set(reflect.Zero(targetField.Type()))
			continue
		}

		jsonFieldType := reflect.TypeOf(jsonValue).String()
		if jsonFieldType == "map[string]interface {}" {
			childJson := jsonValue.(map[string]interface{})
			switch childTarget := targetField.Interface().(type) {
			case VoiceSynthesisConfig:
				allowChildFields := TopLevelStructKeys(&childTarget) // 親フィールドが既に許可されているので子は全て許可
				MergeJsonToStruct(&childJson, &childTarget, &allowChildFields)
				targetField.Set(reflect.ValueOf(&childTarget).Elem())
			case Position3D:
				allowChildFields := TopLevelStructKeys(&childTarget)
				MergeJsonToStruct(&childJson, &childTarget, &allowChildFields)
				targetField.Set(reflect.ValueOf(&childTarget).Elem())
			case LessonDrawingStroke:
				allowChildFields := TopLevelStructKeys(&childTarget)
				MergeJsonToStruct(&childJson, &childTarget, &allowChildFields)
				targetField.Set(reflect.ValueOf(&childTarget).Elem())
			}
		} else if jsonFieldType == "[]interface {}" {
			switch targets := targetField.Interface().(type) {
			case []LessonReference:
				targets = nil // 配列は元の値にマージせず丸ごと置き換える
				for _, v := range jsonValue.([]interface{}) {
					var targetBlankStruct LessonReference
					allowChildFields := TopLevelStructKeys(&targetBlankStruct)
					child := v.(map[string]interface{})
					MergeJsonToStruct(&child, &targetBlankStruct, &allowChildFields)
					targets = append(targets, targetBlankStruct)
				}
				targetField.Set(reflect.ValueOf(&targets).Elem())
			case []LessonAvatar:
				targets = nil
				for _, v := range jsonValue.([]interface{}) {
					var targetBlankStruct LessonAvatar
					allowChildFields := TopLevelStructKeys(&targetBlankStruct)
					child := v.(map[string]interface{})
					MergeJsonToStruct(&child, &targetBlankStruct, &allowChildFields)
					targets = append(targets, targetBlankStruct)
				}
				targetField.Set(reflect.ValueOf(&targets).Elem())
			case []LessonDrawing:
				targets = nil
				for _, v := range jsonValue.([]interface{}) {
					var targetBlankStruct LessonDrawing
					allowChildFields := TopLevelStructKeys(&targetBlankStruct)
					child := v.(map[string]interface{})
					MergeJsonToStruct(&child, &targetBlankStruct, &allowChildFields)
					targets = append(targets, targetBlankStruct)
				}
				targetField.Set(reflect.ValueOf(&targets).Elem())
			case []LessonEmbedding:
				targets = nil
				for _, v := range jsonValue.([]interface{}) {
					var targetBlankStruct LessonEmbedding
					allowChildFields := TopLevelStructKeys(&targetBlankStruct)
					child := v.(map[string]interface{})
					MergeJsonToStruct(&child, &targetBlankStruct, &allowChildFields)
					targets = append(targets, targetBlankStruct)
				}
				targetField.Set(reflect.ValueOf(&targets).Elem())
			case []LessonGraphic:
				targets = nil
				for _, v := range jsonValue.([]interface{}) {
					var targetBlankStruct LessonGraphic
					allowChildFields := TopLevelStructKeys(&targetBlankStruct)
					child := v.(map[string]interface{})
					MergeJsonToStruct(&child, &targetBlankStruct, &allowChildFields)
					targets = append(targets, targetBlankStruct)
				}
				targetField.Set(reflect.ValueOf(&targets).Elem())
			case []LessonMusic:
				targets = nil
				for _, v := range jsonValue.([]interface{}) {
					var targetBlankStruct LessonMusic
					allowChildFields := TopLevelStructKeys(&targetBlankStruct)
					child := v.(map[string]interface{})
					MergeJsonToStruct(&child, &targetBlankStruct, &allowChildFields)
					targets = append(targets, targetBlankStruct)
				}
				targetField.Set(reflect.ValueOf(&targets).Elem())
			case []LessonSpeech:
				targets = nil
				for _, v := range jsonValue.([]interface{}) {
					var targetBlankStruct LessonSpeech
					allowChildFields := TopLevelStructKeys(&targetBlankStruct)
					child := v.(map[string]interface{})
					MergeJsonToStruct(&child, &targetBlankStruct, &allowChildFields)
					targets = append(targets, targetBlankStruct)
				}
				targetField.Set(reflect.ValueOf(&targets).Elem())
			case []LessonDrawingUnit:
				targets = nil
				for _, v := range jsonValue.([]interface{}) {
					var targetBlankStruct LessonDrawingUnit
					allowChildFields := TopLevelStructKeys(&targetBlankStruct)
					child := v.(map[string]interface{})
					MergeJsonToStruct(&child, &targetBlankStruct, &allowChildFields)
					targets = append(targets, targetBlankStruct)
				}
				targetField.Set(reflect.ValueOf(&targets).Elem())
			case []Position2D:
				targets = nil
				for _, v := range jsonValue.([]interface{}) {
					var targetBlankStruct Position2D
					allowChildFields := TopLevelStructKeys(&targetBlankStruct)
					child := v.(map[string]interface{})
					MergeJsonToStruct(&child, &targetBlankStruct, &allowChildFields)
					targets = append(targets, targetBlankStruct)
				}
				targetField.Set(reflect.ValueOf(&targets).Elem())
			}
		} else {
			setValueToField(jsonValue, targetField)
		}
	}
}

func setValueToField(jsonValue interface{}, targetField reflect.Value) {
	originFieldType := targetField.Type().Kind()
	switch originFieldType {
	case reflect.Bool:
		targetField.SetBool(jsonValue.(bool))
	case reflect.Float32, reflect.Float64:
		targetField.SetFloat(jsonValue.(float64))
	case reflect.Int32, reflect.Int64:
		v := jsonValue.(float64) // jsonのintは全てfloatとして扱われるのでキャストする
		targetField.SetInt(int64(v))
	case reflect.String:
		targetField.SetString(jsonValue.(string))
	case reflect.Uint64:
		targetField.SetUint(jsonValue.(uint64))
	default:
		// 独自にUnmarshallでenumを定義した型の場合
		customField := reflect.New(targetField.Type()).Interface()
		bytes := []byte("\"" + jsonValue.(string) + "\"")
		json.Unmarshal(bytes, &customField)
		targetField.Set(reflect.ValueOf(customField).Elem())
	}
}

func TopLevelStructKeys(target interface{}) []string {
	v := reflect.ValueOf(target).Elem().Type()

	var keys []string
	for i := 0; i < v.NumField(); i++ {
		keys = append(keys, v.Field(i).Name)
	}

	return keys
}

func Contains(s *[]string, e string) bool {
	for _, a := range *s {
		if a == e {
			return true
		}
	}
	return false
}
