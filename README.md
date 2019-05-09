# sendmail

`sendmail`是一个发送邮件的Golang包。没有使用外部依赖，易于安装

## 下载

```
go get github.com/flagUpDown/sendmail
```

## 使用

引入包

```go
import "github.com/flagUpDown/sendmail"
```

创建一封信邮件

```go
mail := sendmail.NewMail()
```

编写信封

```go
mail.SetFromEmail("FromEmail@email.com", "FromEmail") // 设置发件人
mail.AddRecipient("Recipient@email.com", "Recipient") // 添加发送人
mail.AddCarbonCopy("CarbonCopy@email.com", "CarbonCopy") // 添加抄送人
mail.AddBlindCarbonCopy("BlindCarbonCopy@email.com", "BlindCarbonCopy") // 添加暗抄送人
mail.SetSubject("subject") // 设置邮件主题
mail.SetContent("<h1>hello word</h1>", true) // 设置邮件内容，可选择是普通文本还是HTML格式的文本
mail.AddAttachment("/path/filename.ext", "file name") // 添加邮件附件
```

连接远程smtp服务器

```go
c, _ := sendmail.Dial("smtp.host.com", 25)
```

设置认证所需的用户名和口令

```go
c.SetAuth("example@email.com", "password")
```

发送邮件

```go
c.Send(mail)
```

关闭客户端连接

```go
c.Close()
```

