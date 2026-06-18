# DNSE SDK for Go (dnse-sdk-go)

Bộ thư viện mã nguồn mở bằng ngôn ngữ Go giúp kết nối và tương tác với hệ thống OpenAPI của Công ty Cổ phần Chứng khoán DNSE, hỗ trợ đầy đủ cả REST API và WebSocket Stream Client.

## Tính năng chính

- **REST Client (`restdnse`)**: 
  - Đăng ký và gia hạn `Trading Token` bằng mã OTP/PIN.
  - Truy vấn danh sách tài khoản và số dư chi tiết.
  - Tra cứu thông số kỹ thuật mã chứng khoán (Trần/Sàn/Tham chiếu), danh sách gói vay margin.
  - Tính toán sức mua/sức bán (`PPSE`) trước khi đặt lệnh.
  - Đẩy lệnh giao dịch (`PlaceOrder`) tự động tích hợp ký chữ ký bảo mật.
- **WebSocket Client (`websocket`)**:
  - Quản lý vòng đời kết nối, tự động gửi gói tin Auth Message và cơ chế giữ kết nối (`Heartbeat`).
  - Đăng ký và lắng nghe luồng dữ liệu thị trường Public (Quote, Trade/Tick, Index).
  - Đăng ký luồng dữ liệu tài sản Private (Trạng thái lệnh, tài khoản, vị thế).
- **Bảo mật (`pkg/crypto`)**: Tự động tính toán chuỗi ký `HMAC-SHA256` chuẩn hóa theo Gateway OpenAPI v2.

## Cài đặt

Yêu cầu phiên bản Go từ `1.25.0` trở lên. 

Khởi tạo cấu hình biến môi trường bằng cách sao chép file cấu hình mẫu:
```bash
cp .env_example .env

```

Cập nhật các thông tin tài khoản của bạn vào file `.env`:

```env
DNSE_BASE_URL=[https://openapi.dnse.com.vn](https://openapi.dnse.com.vn)
DNSE_WS_URL=wss://ws-openapi.dnse.com.vn
DNSE_API_KEY=your_api_key
DNSE_API_SECRET=your_api_secret
DNSE_ACCOUNT_ID=your_account_id

```

## Hướng dẫn sử dụng nhanh

### 1. Sử dụng REST API để lấy Trading Token và đặt lệnh

```go
package main

import (
	"context"
	"fmt"
	"log"
	
	"dnse-sdk-go/restdnse"
)

func main() {
	client := restdnse.NewClient("[https://openapi.dnse.com.vn](https://openapi.dnse.com.vn)", "API_KEY", "API_SECRET")
	ctx := context.Background()

	// Đổi Smart OTP lấy Trading Token
	token, err := client.CreateTradingToken(ctx, "smart_otp", "123456")
	if err != nil {
		log.Fatalf("Lỗi lấy Token: %v", err)
	}
	fmt.Printf("Trading Token: %s\n", token)

	// Đặt lệnh mua/bán cơ sở
	orderReq := restdnse.DNSEOrderRequest{
		AccountNo:     "YOUR_ACCOUNT_ID",
		LoanPackageID: 1769,
		Symbol:        "TCB",
		Side:          "NB", // NB: Mua, NS: Bán
		OrderType:     "LO",
		Price:         33000, // Giá 33.0 VND
		Quantity:      100,
	}

	orderID, err := client.PlaceOrder(ctx, "STOCK", orderReq)
	if err != nil {
		log.Fatalf("Đặt lệnh thất bại: %v", err)
	}
	fmt.Printf("Đặt lệnh thành công, OrderID: %s\n", orderID)
}

```

### 2. Lắng nghe luồng dữ liệu Live Market qua WebSocket

```go
package main

import (
	"fmt"
	"time"
	
	"dnse-sdk-go/websocket"
)

func main() {
	streamClient := websocket.NewDNSEStreamClient("wss://ws-openapi.dnse.com.vn", "API_KEY", "API_SECRET")

	// Đăng ký callback xử lý dữ liệu bảng giá
	streamClient.OnQuote = func(symbol string, data map[string]interface{}) {
		fmt.Printf("[QUOTE] Mã: %s | Dữ liệu: %v\n", symbol, data)
	}

	// Kết nối và bắt đầu nhận luồng dữ liệu cho các mã chỉ định
	_ = streamClient.StartMarketData([]string{"HPG", "FPT"}, true, true, false)

	// Duy trì luồng chạy
	time.Sleep(10 * time.Second)
	streamClient.Close()
}

```

## Thư mục ví dụ (`examples/`)

Xem thêm các kịch bản mẫu chi tiết tại thư mục `examples/` để biết cách vận hành:

* `get_accounts.go`: Truy vấn danh sách tiểu khoản.
* `get_balances.go`: Tra cứu số dư tài sản.
* `get_loanpackage.go`: Lấy danh sách gói margin khả dụng.
* `get_ppse.go`: Tính toán sức mua bán trước lệnh.
* `get_olhcv.go`: Tải lịch sử nến đồ thị kỹ thuật.
* `websocket_stream.go`: Toàn văn kịch bản chạy luồng dữ liệu Realtime.

```

```