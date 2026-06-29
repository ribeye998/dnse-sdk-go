# DNSE SDK for Go (dnse-sdk-go)

Bộ thư viện mã nguồn mở bằng ngôn ngữ Go giúp kết nối và tương tác với hệ thống OpenAPI của Công ty Cổ phần Chứng khoán DNSE, hỗ trợ đầy đủ cả REST API và WebSocket Stream Client.

## Tính năng chính

- **REST Client (`dnse.Client`)**: 
  - Đăng ký và gia hạn `Trading Token` bằng mã smart_otp hoặc email_otp.
  - Truy vấn danh sách tài khoản và số dư chi tiết.
  - Tra cứu thông số kỹ thuật mã chứng khoán (Trần/Sàn/Tham chiếu), danh sách gói vay margin.
  - Tính toán sức mua/sức bán (`PPSE`) trước khi đặt lệnh.
  - Đẩy lệnh giao dịch (`PlaceOrder`), sửa lệnh (`AmendOrder`), hủy lệnh (`CancelOrder`) tự động tích hợp ký chữ ký bảo mật chuẩn hóa.
- **WebSocket Client (`dnse.StreamClient`)**:
  - Quản lý vòng đời kết nối, tự động gửi gói tin Auth Message và cơ chế giữ kết nối (`Heartbeat`).
  - Đăng ký và lắng nghe luồng dữ liệu thị trường Public (Bảng giá Quote, Lịch sử khớp lệnh Trade/Tick, Chỉ số Index, Nến đồ thị OHLC).
  - Đăng ký luồng dữ liệu tài sản Private (Trạng thái lệnh, số dư tài khoản, danh mục vị thế).
- **Bảo mật**: Tự động tính toán chuỗi ký `HMAC-SHA256` chuẩn hóa theo Gateway OpenAPI v2.

## Cài đặt

Yêu cầu phiên bản Go từ `1.21.0` trở lên. 

Thêm module vào dự án:
```bash
go get github.com/ribeye998/dnse-sdk-go
```

Khởi tạo cấu hình biến môi trường bằng cách sao chép file cấu hình mẫu:
```bash
cp .env_example .env
```

Cập nhật các thông tin tài khoản của bạn vào file `.env`:
```env
DNSE_BASE_URL=https://openapi.dnse.com.vn
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
	
	dnse "github.com/ribeye998/dnse-sdk-go"
)

func main() {
	client := dnse.NewClient("https://openapi.dnse.com.vn", "YOUR_API_KEY", "YOUR_API_SECRET")
	ctx := context.Background()

	// Đổi Smart OTP lấy Trading Token
	token, err := client.CreateTradingToken(ctx, "smart_otp", "123456")
	if err != nil {
		log.Fatalf("Lỗi lấy Token: %v", err)
	}
	fmt.Printf("Trading Token: %s\n", token)

	// Đặt lệnh mua/bán cơ sở (Stock)
	orderReq := dnse.OrderRequest{
		AccountNo:     "YOUR_ACCOUNT_ID",
		LoanPackageID: 1769,
		Symbol:        "TCB",
		Side:          dnse.SideBuy, // BUY
		OrderType:     "LO",
		Price:         33000,
		Quantity:      100,
	}

	result, err := client.PlaceOrder(ctx, dnse.MarketStock, "NORMAL", orderReq)
	if err != nil {
		log.Fatalf("Đặt lệnh thất bại: %v", err)
	}
	fmt.Printf("Đặt lệnh thành công, Kết quả: %s\n", string(result))
}
```

### 2. Lắng nghe luồng dữ liệu Live Market qua WebSocket

```go
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	
	dnse "github.com/ribeye998/dnse-sdk-go"
)

func main() {
	stream := dnse.NewStreamClient("wss://ws-openapi.dnse.com.vn", "YOUR_API_KEY", "YOUR_API_SECRET")

	// Đăng ký callback xử lý dữ liệu
	stream.OnQuote(func(symbol string, data map[string]interface{}) {
		fmt.Printf("[QUOTE] Mã: %s | Dữ liệu: %v\n", symbol, data)
	})
	stream.OnTick(func(symbol string, data map[string]interface{}) {
		fmt.Printf("[TICK] Mã: %s | Dữ liệu: %v\n", symbol, data)
	})

	// Kết nối và bắt đầu nhận luồng dữ liệu cho HPG, FPT trên bảng G1
	err := stream.StartMarketData([]string{"HPG", "FPT"}, true, true, false)
	if err != nil {
		log.Fatalf("Lỗi kết nối stream: %v", err)
	}

	fmt.Println("Đang chạy WebSocket stream — nhấn Ctrl+C để dừng")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	stream.Close()
	fmt.Println("Đã đóng kết nối.")
}
```

---

## Danh mục mã nguồn ví dụ (`examples/`)

Xem thêm các kịch bản mẫu chi tiết tại thư mục `examples/` để biết cách vận hành:

- **Giao dịch (Trading APIs)**:
  - `get_accounts`: Truy vấn danh sách tiểu khoản.
  - `get_balances`: Tra cứu số dư tài sản chi tiết.
  - `get_loan_packages`: Lấy danh sách các gói margin khả dụng.
  - `get_ppse`: Tính toán sức mua bán trước lệnh.
  - `get_orders`: Truy vấn sổ lệnh trong ngày.
  - `get_order_detail`: Xem thông tin chi tiết một lệnh cụ thể.
  - `get_order_history`: Tra cứu lịch sử đặt lệnh.
  - `get_corporate_action_history`: Truy vấn lịch sử quyền/sự kiện doanh nghiệp.
  - `get_execution_detail`: Xem chi tiết kết quả khớp lệnh (Chỉ hỗ trợ Derivative).
  - `get_positions`: Truy vấn các vị thế danh mục đang nắm giữ.
  - `get_positions_by_id`: Xem chi tiết vị thế cụ thể theo ID.
  - `close_position`: Đóng vị thế phái sinh đang có.
  - `create_trading_token`: Xác thực bằng smart_otp hoặc email_otp lấy Trading Token.
  - `place_order`: Đặt lệnh mua/bán mới.
  - `cancel_order`: Hủy một lệnh chưa khớp.
  - `replace_order`: Sửa giá/khối lượng của lệnh chưa khớp.

- **Dữ liệu thị trường (Market Data APIs)**:
  - `get_secdef`: Truy vấn chi tiết thông số kỹ thuật mã (Giá trần/sàn/tham chiếu).
  - `get_instruments`: Tra cứu thông tin danh sách các mã giao dịch.
  - `get_trades`: Xem lịch sử giao dịch khớp lệnh của mã.
  - `get_latest_trade`: Lấy thông tin phiên khớp lệnh gần nhất của mã.
  - `get_ohlc`: Tải lịch sử nến đồ thị kỹ thuật.
  - `get_close_price`: Xem giá đóng cửa gần nhất của mã.
  - `get_working_dates`: Xem danh sách lịch ngày làm việc/giao dịch của thị trường.

- **Kênh WebSocket Stream**:
  - `websocket_sec_def`: Lắng nghe biến động thông số kỹ thuật mã chứng khoán.
  - `websocket_quote`: Lắng nghe biến động giá chào mua/chào bán tốt nhất (Top Book).
  - `websocket_trade`: Lắng nghe các giao dịch khớp lệnh thời gian thực.
  - `websocket_trade_extra`: Lắng nghe chi tiết khớp lệnh kèm khối lượng mua/bán chủ động.
  - `websocket_ohlc`: Lắng nghe cập nhật nến đồ thị thời gian thực.
  - `websocket_ohlc_closed`: Lắng nghe cập nhật nến đồ thị khi đóng nến.
  - `websocket_expected_price`: Lắng nghe giá dự kiến khớp trong phiên ATO/ATC.
  - `websocket_foreign_investor`: Lắng nghe giao dịch khối ngoại.
  - `websocket_market_index`: Lắng nghe thông số các chỉ số thị trường (e.g. `VN30`).
  - `websocket_estimated_market_index`: Lắng nghe ước lượng các chỉ số thị trường (e.g. `VN30`).
  - `websocket_trading`: Lắng nghe biến động tài sản cá nhân (Lệnh giao dịch, Danh mục vị thế, Số dư tài khoản).