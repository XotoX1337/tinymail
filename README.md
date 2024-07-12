
[![Go Reference](https://pkg.go.dev/badge/github.com/XotoX1337/tinymail.svg)](https://pkg.go.dev/github.com/XotoX1337/tinymail)
[![Go Report Card](https://goreportcard.com/badge/github.com/XotoX1337/tinymail)](https://goreportcard.com/report/github.com/XotoX1337/tinymail)

# tinymail
tinymail is a small package to easily send simple emails in go.

## Download
```
go get github.com/XotoX1337/tinymail
```
## Features

* SMTP Authentification
* Email with text body
* Email from Template as String or File
* Attachments

## Examples

### Text Email
```go
import "github.com/XotoX1337/tinymail"

opts := tinymail.MailerOpts{
    User: "username",
    Password: "password",
    Host: "host",
    Port: 587
}
mailer := tinymail.New(opts)
msg := tinymail.FromText("this is a example")
msg.SetFrom("test@tinymail.test")
msg.SetTo("test.to@tinymail.test")
msg.SetSubject("TestWriteMessage")
err := mailer.SetMessage(msg).Send()
if err != nil {
    fmt.Println(err)
}
# send success
```

### Email from Template
```go
import "github.com/XotoX1337/tinymail"

opts := tinymail.MailerOpts{
    User: "username",
    Password: "password",
    Host: "host",
    Port: 587
}
mailer := tinymail.New(opts)
msg := tinymail.FromTemplateFile(path/to/template/file)
msg.SetFrom("test@tinymail.test")
msg.SetTo("test.to@tinymail.test")
msg.SetSubject("TestWriteMessage")
err := mailer.SetMessage(msg).Send()
if err != nil {
    fmt.Println(err)
}
# send success
```

### Email with Attachments
```go
import "github.com/XotoX1337/tinymail"

opts := tinymail.MailerOpts{
    User: "username",
    Password: "password",
    Host: "host",
    Port: 587
}
mailer := tinymail.New(opts)
msg := tinymail.FromText("attachment example")
msg.SetFrom("test@tinymail.test")
msg.SetTo("test.to@tinymail.test")
msg.SetSubject("TestWriteMessage")
msg.Attach(path/to/file, path/to/second/file, ...)
err := mailer.SetMessage(msg).Send()
if err != nil {
    fmt.Println(err)
}
# send success
```

