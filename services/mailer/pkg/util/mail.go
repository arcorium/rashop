package util

import (
  "fmt"
  vob "github.com/arcorium/rashop/services/mailer/internal/domain/valueobject"
  "time"
)

func mailBody(header, token string, expiryTime time.Time) (string, vob.BodyType) {
  return fmt.Sprintf(`
  <html>
    <body>
    <h1>%s Token:h1>
    <p>%s</p>
    <p> Token will be expired at: %v </p>
    </body>
  </html>
  `, header, token, expiryTime), vob.BodyTypeHtml
}

func CreateEmailVerificationMailBody(token string, expiryTime time.Time) (string, vob.BodyType) {
  return mailBody("Email Verification", token, expiryTime)
}

func CreateResetPasswordMailBody(token string, expiryTime time.Time) (string, vob.BodyType) {
  return mailBody("Reset Password", token, expiryTime)
}

func CreateLoginTokenMailBody(token string, expiryTime time.Time) (string, vob.BodyType) {
  return mailBody("Login", token, expiryTime)
}
