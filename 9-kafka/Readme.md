## Apache Kafka
**Định nghĩa:** là một nền tảng stream dữ liệu phân tán
Kafka cung cấp 3 chức năng chính cho người dùng:
- Xuất bản và đăng kí các luồng bản ghi
- Lưu trữ hiệu quả các luồng bản ghi theo thứ tự tạo bản ghi
- Xử lý các luồng bản ghi trong thời gian thực

Sơ đồ khái quát các tính năng của Kafka
<img src="https://images.viblo.asia/95772af1-1227-41ea-a13d-6ca79be78e60.png"/>

### Mô tả cấu trúc cơ bản của Kafka
<img src="https://images.viblo.asia/e688b9a2-daf0-4efd-91be-529b368a9e2c.png"/>

<img src="https://images.viblo.asia/eabf0b4b-2cf2-4398-a20e-9dd312a93fb7.png" />

**Producer:** kafka lưu, phân loại message theo topic, sử dụng producer để publish message vào các topic. Dữ liệu được gửi đến partition của topic lưu trữ trên **Broker**
**Consumer:** Kafka sử dụng `consumer` để subcribe vào topic, các consumer đợc định danh bằng các group name. Nhiều `consumer` có thể cùng đọc một topic
**Topic:** Dữ liệu truyền trong kafka theo topic, khi cần truyền dữ liệu cho các ứng dụng khác nhau thì tạo ra các topic khác nhau
**Partition:** Đây là nơi dữ liệu cho 1 `topic` được lưu trữ. Một `topic` có thể có một hoặc nhiều `partition`. Trên mỗi `partition` thì dữ liệu lưu trữ cố định và được gán cho một ID gọi là `offset`. Trong một Kafka cluster thì 1 partition có thể sao chép ra nhiều bản. Trong 1 bản leader chịu trách nhiệm đọc ghi dữ liệu và các bản còn lại gọi là `follower`. Khi bản `leader` bị lỗi thì sẽ có 1 bản `follower` lên làm `leader` thay thế. Nếu muốn dùng nhiều consumer đọc song song dữ liệu của 1 topic thì topic đó cần phải có nhiều `partition`
**Broker:** Kafka cluster là 1 set các server, mỗi một set này được gọi là 1 broker
**Zookeeper:** được dùng để quản lý và bố trí các broker

**[Một vài định nghĩa khác tham khảo tại đây](https://viblo.asia/p/kafka-nhung-khai-niem-thuat-ngu-va-giai-thich-ve-nhung-thu-ma-kafka-co-the-lam-duoc-MkNLrZ9wLgA)**


## Mô tả về ví dụ cơ bản:

Gồm 3 folders:
- `models`: gồm 2 struct:
    - `User`: thông tin user
    - `Notification`: thông tin thông báo
- `cmd/producer`: setup 1 server để gửi message 
- `cmd/consumer`: setup 1 server để nhận message

- **Các phần kiến thức kafka qua demo**:
    - Việc khởi tạo 1 topic(1 partition)
    - Producer gửi message vào topic, vì chỉ có 1 partition nên tất cả message sẽ gửi vào nó.
    - Consumer chờ và xử lý notifications nhận được
- **Các phần kiến thức sẽ làm thêm**
    - Thêm nhiều partition vào topic, việc phân phối message vào các partition theo `key` sẽ như thế nào
    - Run nhiều consumer cùng lúc(**Consumer Group**) và xem việc nó nhận message và xử lý từng partition của từng consumer như thế nào