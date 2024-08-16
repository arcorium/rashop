package valueobject

type MediaType string

const (
  MediaTypeNormal  MediaType = "normal"
  MediaTypeOneTime MediaType = "one-time"
)

type Media struct {
  Name string
  Type MediaType
  Data []byte
}
