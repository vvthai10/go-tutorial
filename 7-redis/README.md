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
  ├── publisher # publisher  
  │ ├── publisher.go  
  │ └── options.go  
  ├── consumer # consumer
  | ├── consumer.go  
  │ └── options.go  
  ├── connection.go  
  └── errors.go

- Mình sẽ giải thích từng phần ở phía dưới đây

  - errors.go: chứa các biến để thông báo ra các lỗi cơ bản khi thực hiện.
  - connection.go: chứa các hàm để thực hiện kết nối với consumer. Khi kết nối với consumer sẽ gồm các bước: ...
  - consumer
    - options.go: Khởi tạo thêm các middleware khi tạo consumer
    - consumer.go: Chứa các logic cho 1 consumer
      - func New: Khởi tạo các phần cơ bản của consumer: channel, exchange, queue,
      - func consumer: Hàm đảm nhận việc đợi nhận các tin nhắn gửi vào
      - func consumerCall: Thực hiện các logic cho các message gửi vào
  - publisher
    - options.go: Khởi tạo thêm các middleware khi tạo publisher
    - publisher.go: Chứa các logic cho 1 publisher
      - func New: Khởi tạo các phần cơ bản của publisher
      - func consumer: Thực hiện việc nhận response từ consumer
      - func RemoteCall: Thực hiện việc xử lý và gửi tin nhắn tới consumer

- **connection.go:**
  - Gồm struct `Connection`, dùng để khởi tạo 1 kết nối đến server rabbitmq.
  - Trong hàm khởi tạo, thì sẽ gồm các bước:
    - Tạo 1 exchange
    - Tạo 1 queue
    - Thêm queue và exchange
    - Cho queue bắt đầu nhận tin

### Version 1

- **consumer:**

  - `type CallHandler func(*amqp.Delivery) (interface{}, error)`: format hàm để xử lý các nhiệm vụ.
  - Class `Consumer`: chứa các phần khởi tạo và phân luồng cho các xử lý nhiệm vụ khi nhận được message
    - Hàm `New`: gọi để khởi tạo 1 server
    - Hàm `consumer`: gọi để bắt đầu chạy server, lắng nghe khi nhận được tin hoặc reconnect
    - Hàm `consumerCall`: Thực thi hàm xử lý logic cho nội dung message
    - Hàm `publish`: Gửi lại phản hồi cho `publisher`

- **publisher:**
  - Struct `Message`: dùng để format thông tin gửi đi
  - Struct `pendingCall`: dùng để chứa thông tin về tin nhắn đang ở trạng thái nào
  - Struct `Publisher`: dùng để khởi tạo 1 publisher
    - Hàm `New`: Dùng để khởi tạo vả run
    - Hàm `consumer`: dùng để nhận các luồng xử lý: reconnect, nhận message phản hồi
    - Hàm `getCall`: dùng để xem xét cập nhật tình hình xử lý tin nhắn
    - Hàm `publish`: gửi tin nhắn đến exchange chỉ định
    - Hàm `RemoteCall`: Dùng để gửi tin nhắn - dùng ở main.go
      - Khởi tạo 1 `pendingCall`
      - Gửi tin nhắn bằng `publish`
      - Hàm `addCall` và `deleteCall` dùng để đẩy và xóa tin nhắn trước và sau khi xử lý xong

### Version 2

**Cập nhật:**

- Cho phép tạo các exchange với các type khác nhau(fanout, direct, topic)
- Thêm nhiều queue vào exchange

- `connection.go`

  - Gồm các structs:
    - `Queue`: gồm các params: `name`, `key`, `amqp.Queue`
    - `Exchange`: gồm các params: `name`, `type`
    - `Connection`:
  - Gồm các functions:
    - `New`, `Attempts`, `connect`: đảm nhận việc tạo kết nối đến server, khởi tạo exchange và các queue. Kết nối các queue vào exchange

- `consumer`: tương tự như version 1, update việc **lắng nghe message từ nhiều channel**
- `publisher`: đơn giản hơn version 1, do chỉ có mục đích là gửi tin nhắn.
