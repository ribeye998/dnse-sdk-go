package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"dnse-sdk-go/restdnse"

	"github.com/joho/godotenv"
	"golang.org/x/term"
)

func main() {
	_ = godotenv.Load()

	apiKey := os.Getenv("DNSE_API_KEY")
	apiSecret := os.Getenv("DNSE_API_SECRET")
	accountNo := os.Getenv("DNSE_ACCOUNT_ID")

	if apiKey == "" || apiSecret == "" || accountNo == "" {
		log.Fatalf("Lỗi: Vui lòng cấu hình đầy đủ DNSE_API_KEY, DNSE_API_SECRET và DNSE_ACCOUNT_ID trong file .env")
	}

	// Yêu cầu người dùng nhập mã Smart OTP từ ứng dụng DNSE
	fmt.Print("Nhập mã Smart OTP của bạn: ")
	bytePin, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		log.Fatalf("Lỗi không thể đọc mã OTP: %v", err)
	}

	myOtp := strings.TrimSpace(string(bytePin))
	if myOtp == "" {
		log.Fatalf("Lỗi: Mã OTP không được để trống")
	}

	client := restdnse.NewClient("https://openapi.dnse.com.vn", apiKey, apiSecret)
	ctx := context.Background()

	// 1. Thực hiện đổi mã Smart OTP lấy Trading Token
	log.Println("Đang yêu cầu cấp Trading Token bằng Smart OTP...")
	token, err := client.CreateTradingToken(ctx, "smart_otp", myOtp)
	if err != nil {
		log.Fatalf("Lỗi lấy Token: %v", err)
	}
	fmt.Printf("=> Lấy Trading Token thành công! Chuỗi mã: %s\n\n", token)

	// 2. Truy vấn số dư tài khoản
	var balanceRes interface{}
	log.Println("Đang truy vấn số dư tài khoản...")
	err = client.GetBalances(ctx, accountNo, &balanceRes)
	if err != nil {
		log.Fatalf("Truy vấn số dư thất bại: %v", err)
	}
	fmt.Printf("Dữ liệu tài sản trả về: %+v\n\n", balanceRes)

	// 3. Cấu hình phân loại thị trường chuẩn ĐẶC TẢ của DNSE: STOCK cho cơ sở
	marketType := "STOCK"

	// 4. Khởi tạo cấu trúc dữ liệu đặt lệnh
	orderRequest := restdnse.DNSEOrderRequest{
		AccountNo:     accountNo,
		LoanPackageID: 1769,
		Symbol:        "TCB",
		Side:          "NS", // NB Buy / NS Sell
		OrderType:     "LO",
		Price:         33000, // Đơn vị Đồng (33000 tương đương giá 33.0 trên bảng điện)
		Quantity:      100,
	}

	log.Printf("Đang thực hiện đặt lệnh kiểm thử qua REST %s: %s %d %s với giá %d...\n",
		marketType, orderRequest.Side, orderRequest.Quantity, orderRequest.Symbol, orderRequest.Price)

	// 5. Gọi hàm PlaceOrder gửi lên cổng /accounts/orders kèm query params tự động
	orderID, err := client.PlaceOrder(ctx, marketType, orderRequest)
	if err != nil {
		log.Fatalf("Lỗi trong quá trình đẩy lệnh lên sàn: %v", err)
	}

	// 6. In thông tin kết quả thành công hoàn toàn
	fmt.Println("--------------------------------------------------")
	fmt.Printf("Chúc mừng! Đẩy lệnh thành công lên hệ thống DNSE.\n")
	fmt.Printf("Mã định danh Lệnh nhận về (OrderID): %s\n", orderID)
	fmt.Println("--------------------------------------------------")
}
