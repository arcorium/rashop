package util

import (
  "mime"
  "path/filepath"
)

func GetMimeType(filename string) string {
  types := mime.TypeByExtension(filepath.Ext(filename))
  if types == "" {
    return "application/octet-stream"
  }
  return types
}
