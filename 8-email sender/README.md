## [GOLANG] Email Sender

- Triển khai 2 cách gửi email trong golang:

  - Dùng SMTP truyền thống
  - Dùng SendGrid

- Cấu trúc một `interface` cho cả 2 cách:

```
  type IEmailSender interface {
	Send(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
		args interface{},
	) (*Request, *Response)
}
```

- Các params trong `Send()` gổm:

  - `subject`: tên của email gửi
  - `content`: nội dung gửi đi(có thể là string, html)
  - `to`: danh sách người nhận
  - ...
  - **`args`: các cấu hình đặc biệt dùng cho các dịch vụ gửi mail khác nhau**

- Cấu trúc `Request` và `Response`:

```
  type Request struct {
    Subject     string
    Content     string
    To          []string
    Cc          []string
    Bcc         []string
    AttachFiles []string
  }

  type Response struct {
    Success bool
    Message string
  }
```

- **SMTP**

  - Sử dụng package `github.com/jordan-wright/email`

- **SendGrid**
  - Sử dụng package `github.com/sendgrid/sendgrid-go`
  - Sử dụng `args` dưới dạng 1 `struct` cung cấp khả năng lựa chọn templateID
    ```
      type SendGridOptions struct {
        TemplateID string
      }
    ```
