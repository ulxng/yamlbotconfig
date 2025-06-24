package email

type SmtpConfig struct {
	Username         string `long:"username"           env:"USER"                 required:"true" description:"SMTP username"`
	Password         string `long:"password"           env:"PASS"                 required:"true" description:"SMTP password"`
	Host             string `long:"host"               env:"HOST"                 required:"true" description:"SMTP host"`
	Port             int    `long:"port"               env:"PORT"                 required:"true" description:"SMTP port"`
	DefaultEmailFrom string `long:"default-email-from" env:"DEFAULT_EMAIL_FROM"   required:"true"`
	DefaultEmailTo   string `long:"default-email-to"   env:"DEFAULT_EMAIL_TO"     required:"true"`
}
