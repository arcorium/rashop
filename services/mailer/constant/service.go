package constant

const (
  SERVICE_NAME         = "mailer"
  SERVICE_VERSION      = "v1"
  SERVICE_COMMAND_NAME = SERVICE_NAME + "::command"
  SERVICE_QUERY_NAME   = SERVICE_NAME + "::query"
)

const (
  DEFAULT_EMAIL_SENDER     = "rashop@no-reply.com"
  EmailVerificationSubject = "Rashop - Email Verification"
  ResetPasswordSubject     = "Rashop - Reset Password"
  LoginSubject             = "Rashop - Login Token"
)
