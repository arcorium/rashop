package mapper

import (
  "errors"
  "reflect"
)

const (
  skipFieldMarker = "-"
  structTag       = "map"
)

func getFieldName(field *reflect.StructField) (string, bool) {
  fieldName, ok := field.Tag.Lookup(structTag)
  if !ok {
    fieldName = field.Name
  }
  if fieldName == skipFieldMarker {
    return "", false
  }
  return fieldName, true
}

func getRealType(field reflect.Type) reflect.Type {
  for field.Kind() == reflect.Pointer {
    field = field.Elem()
  }
  return field
}

func getRealValue(value reflect.Value) reflect.Value {
  for value.Kind() == reflect.Pointer {
    value = value.Elem()
  }
  return value
}

type sourceData struct {
  index int
}

func Bind[D, S any](src S, dst D) error {
  srcType := getRealType(reflect.TypeFor[S]())
  srcVal := getRealValue(reflect.ValueOf(src))

  dstType := reflect.TypeOf(dst)

  if srcType.Kind() != reflect.Struct ||
      dstType.Kind() != reflect.Pointer ||
      dstType.Elem().Kind() != reflect.Struct {
    return errors.New("expected struct for source and pointer to struct for dest")
  }

  lookup := createLookup(srcType)

  // Set dest
  dstVal := reflect.ValueOf(dst)
  set(&srcVal, &dstVal, dstType.Elem(), lookup)
  return nil
}

func createLookup(p reflect.Type) map[string]sourceData {
  lookup := make(map[string]sourceData)
  // Scan source struct
  for i := range p.NumField() {
    currField := p.Field(i)
    currFieldName, ok := getFieldName(&currField)
    if !ok {
      continue
    }
    lookup[currFieldName] = sourceData{
      index: i,
    }
  }

  return lookup
}

func set(srcVal *reflect.Value, destVal *reflect.Value, destType reflect.Type, lookup map[string]sourceData) {
  for i := range destType.NumField() {
    currField := destType.Field(i)
    // Check if tag is exist
    currFieldName, ok := getFieldName(&currField)
    if !ok {
      continue
    }
    datas, ok := lookup[currFieldName]
    if !ok {
      continue
    }
    currVal := destVal.Elem().Field(i)
    if !currVal.CanSet() {
      continue
    }
    currVal.Set(srcVal.Field(datas.index))
  }
}

func Map[D, S any](src S) (D, error) {
  srcType := getRealType(reflect.TypeFor[S]())
  srcVal := getRealValue(reflect.ValueOf(src))

  var dst D
  dstType := reflect.TypeOf(dst)

  // Expect struct for both
  if srcType.Kind() != reflect.Struct || dstType.Kind() != reflect.Struct {
    return dst, errors.New("expected struct for both type")
  }

  lookup := createLookup(srcType)

  // Set dest
  dstVal := reflect.ValueOf(&dst)
  set(&srcVal, &dstVal, dstType, lookup)
  return dst, nil
}
