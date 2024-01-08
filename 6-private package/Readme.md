### Use private package

#### Tạo 1 private package:

- `go mod init github.com/username/repo`
- `main.go`:

  ```
      package private

      import "fmt"

      func PrintInfo() {
          fmt.Println("This is a private package.")
      }
  ```

- Publish lên github ở chế độ private

#### Cấu hình sử dụng private package
- Cấu hình biến môi trường: `go env -w GOPRIVATE=go mod init github.com/username/repo/private`
- Check lại biến môi trường: `go env GOPRIVATE`
- Cấu hình SSH thay cho HTTP
    - Dùng cmd: `git config --global url."git@github.com:username/repo".insteadOf "https://github.com/username/repo"`
    - Sửa file(ubuntu): `cat ~/.gitconfig`. window: `notepad $HOME\.gitconfig`
