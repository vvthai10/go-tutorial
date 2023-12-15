## [GOLANG] RabbitMQ

<img src="https://images.viblo.asia/fc8f94df-220f-46be-8769-e62a5c52dd45.png"/>

- RabbitMQ là một AMQP message broker(phần mềm quản lý hàng đợi message)
- Message broker: là một chương trình trung gian được thiết kế để validating, transforming và routing messages. Phục vụ các nhu cầu giao tiếp giữa các ứng dụng với nhau
- Với các ví dụ cơ bản thì bạn có thể tìm kiếm ở các nguồn sau
  - https://viblo.asia/p/connecting-to-rabbitmq-from-go-GrLZDanBlk0
  - https://rohinivsenthil.medium.com/rabbitmq-in-golang-getting-started-34c65e6c7f92

**Ở đây mình sẽ cố gắng để cấu hình nó thành 1 module có thể tái sử dụng ở nhiều hệ thống khác nhau**

- Cấu trúc module sẽ như sau
  .
  ├── client # publisher  
  │ ├── client.go  
  │ └── options.go  
  ├── server # consumer
  | ├── server.go  
  │ └── options.go  
  ├── connection.go  
  └── errors.go

- Mình sẽ giải thích từng phần ở phía dưới đây
  - errors.go: chứa các biến để thông báo ra các lỗi cơ bản khi thực hiện.
  - connection.go: chứa các hàm để thực hiện kết nối với server. Khi kết nối với server sẽ gồm các bước: ...
  - server
    - options.go: Khởi tạo thêm các middleware khi tạo server
    - server.go: Chứa các logic cho 1 server
      - func New: Khởi tạo các phần cơ bản của server: channel, exchange, queue,
      - func consumer: Hàm đảm nhận việc đợi nhận các tin nhắn gửi vào
      - func serverCall: Thực hiện các logic cho các message gửi vào
  - client
    - options.go: Khởi tạo thêm các middleware khi tạo client
    - client.go: Chứa các logic cho 1 client
      - func New: Khởi tạo các phần cơ bản của client
      - func consumer: Thực hiện việc nhận response từ Server
      - func RemoteCall: Thực hiện việc xử lý và gửi tin nhắn tới Server
