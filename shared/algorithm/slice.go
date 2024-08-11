package algo

func IndexOf[T ~[]E, E comparable](s T, e E) (int, bool) {
  for i, v := range s {
    if v != e {
      continue
    }
    return i, true
  }
  return -1, false
}

func IndexOfFunc[S ~[]E, E, T any](s S, b T, pred func(*E, T) bool) (int, bool) {
  for i, v := range s {
    if !pred(&v, b) {
      continue
    }
    return i, true
  }
  return -1, false
}

// FilterByIndices will filter out data from slices that is not in indices parameter
func FilterByIndices[S ~[]E, E any](s S, indices ...int) S {
  result := make(S, len(indices))
  for _, idx := range indices {
    result = append(result, s[idx])
  }
  return result
}

// FilterByIndicesPointing works like FilterByIndices, but instead of copying each element it will use pointer
// to point the real object. This function better used for temporary data as view
func FilterByIndicesPointing[S ~[]E, E any](s S, indices ...int) []*E {
  result := make([]*E, len(indices))
  for _, idx := range indices {
    result = append(result, &s[idx])
  }
  return result
}
