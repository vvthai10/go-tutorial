## [GOLANG] Command Design Pattern

**Thuộc loại behaviroal**

- Mục đích: Dùng để gói lại tất cả các thông tin cũng như action trên các logic khác nhau
- Command design pattern có 4 điều khoản được liên kết với nhau:

  - **Command:** A command được biết về việc nhận và thực thi các method của việc nhận. Giá trị các biến của phương thức nhận được lưu trong các command
  - **Receiver:** THe receiver là object để thực thi những phương thức và cũng được store trong command thông qua command object bởi aggregation, receiver được thực hiện khi execute() method trong command được gọi
  - **Invoker:** The invoker object được biết là thực hiện execute() như thế nào, nó được xem như một người kế toán để thực thi các command và chỉ biết về các **command interface**
  - **Client:** The client quyết định sẽ thực thi tại điểm nào, command nào và pass command object đến invoker object

- Một ví dụ về một logic không dùng mẫu command design pattern
  <img src="https://images.viblo.asia/fcf991fb-559e-44d1-a80f-878b6f02dd70.png"/>
  ==> Nhìn như trên cũng biết việc mở rộng, bảo trì và xử lý logic sẽ khó khăn về sau rồi phải không

- Với ví dụ trên nếu sử dụng **command design pattern**
  <img src="https://images.viblo.asia/7f47bd93-c76d-46d6-b7c5-3de9cbec8c24.png" />
