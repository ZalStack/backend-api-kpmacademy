package handlers

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"backend-api-kpmacademy/database"
	"backend-api-kpmacademy/models"

	"github.com/gofiber/fiber/v2"
)

type MidtransNotification struct {
	OrderID           string `json:"order_id"`
	TransactionID     string `json:"transaction_id"`
	StatusCode        string `json:"status_code"`
	TransactionTime   string `json:"transaction_time"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	TransactionStatus string `json:"transaction_status"`
	Status            string `json:"status"`
	FraudStatus       string `json:"fraud_status"`
	SignatureKey      string `json:"signature_key"`
}

func verifyMidtransSignature(notification *MidtransNotification, serverKey string) bool {
	if os.Getenv("APP_ENV") != "production" {
		return true
	}
	signatureString := notification.OrderID + notification.StatusCode + notification.GrossAmount + serverKey

	h := hmac.New(sha512.New, []byte(serverKey))
	h.Write([]byte(signatureString))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(expectedSignature), []byte(notification.SignatureKey))
}

func mapMidtransStatus(status string) string {
	switch status {
	case "capture", "settlement":
		return statusPaid
	case statusPending:
		return statusPending
	case "deny", "cancel", "expire":
		return "failed"
	case "refund":
		return "refunded"
	default:
		return status
	}
}

func MidtransWebhook() fiber.Handler {
	return func(c *fiber.Ctx) error {
		serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
		if serverKey == "" && os.Getenv("APP_ENV") == "production" {
			log.Println("MIDTRANS_SERVER_KEY not configured")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Payment gateway not configured"})
		}

		var notification MidtransNotification
		if err := c.BodyParser(&notification); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		}

		if notification.OrderID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Missing required fields"})
		}

		if notification.TransactionStatus == "" && notification.Status != "" {
			notification.TransactionStatus = notification.Status
		}

		if !verifyMidtransSignature(&notification, serverKey) {
			log.Printf("Invalid Midtrans signature for order %s", notification.OrderID)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid signature"})
		}

		db := database.GetDB()
		paymentStatus := mapMidtransStatus(notification.TransactionStatus)

		var order models.Order
		if result := db.Where("order_number = ?", notification.OrderID).First(&order); result.Error != nil {
			log.Printf("Order not found: %s", notification.OrderID)
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
		}

		tx := db.Begin()
		if tx.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to start transaction"})
		}

		order.PaymentStatus = paymentStatus
		order.TransactionID = &notification.TransactionID
		order.PaymentType = notification.PaymentType
		order.PaymentReference = notification.OrderID

		if notification.TransactionTime != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", notification.TransactionTime); err == nil {
				order.PaymentTime = &t
			}
		}

		if paymentStatus == statusPaid && order.PackageID != nil {
			var pkg models.Package
			if result := tx.Where("id = ?", order.PackageID).First(&pkg); result.Error == nil {
				now := time.Now()
				dur := pkg.MembershipDurationDays
				order.MembershipDurationDays = &dur
				order.MembershipStart = &now
				end := now.AddDate(0, 0, pkg.MembershipDurationDays)
				order.MembershipEnd = &end
			}
		}

		if result := tx.Save(&order); result.Error != nil {
			tx.Rollback()
			log.Printf("Failed to update order %s: %v", notification.OrderID, result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to update order"})
		}

		if err := tx.Commit().Error; err != nil {
			log.Printf("Failed to commit order update %s: %v", notification.OrderID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to commit transaction"})
		}

		log.Printf("Payment webhook processed: order=%s status=%s", notification.OrderID, paymentStatus)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	}
}

func MidtransVideoWebhook() fiber.Handler {
	return func(c *fiber.Ctx) error {
		serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
		if serverKey == "" && os.Getenv("APP_ENV") == "production" {
			log.Println("MIDTRANS_SERVER_KEY not configured")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Payment gateway not configured"})
		}

		var notification MidtransNotification
		if err := c.BodyParser(&notification); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		}

		if notification.OrderID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Missing required fields"})
		}

		if notification.TransactionStatus == "" && notification.Status != "" {
			notification.TransactionStatus = notification.Status
		}

		if !verifyMidtransSignature(&notification, serverKey) {
			log.Printf("Invalid Midtrans signature for video order %s", notification.OrderID)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid signature"})
		}

		db := database.GetDB()
		paymentStatus := mapMidtransStatus(notification.TransactionStatus)

		var order models.VideoOrder
		if result := db.Where("order_number = ?", notification.OrderID).First(&order); result.Error != nil {
			log.Printf("Video order not found: %s", notification.OrderID)
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
		}

		tx := db.Begin()
		if tx.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to start transaction"})
		}

		order.PaymentStatus = paymentStatus
		order.TransactionID = &notification.TransactionID
		order.PaymentType = notification.PaymentType
		order.PaymentReference = notification.OrderID

		if notification.TransactionTime != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", notification.TransactionTime); err == nil {
				order.PaymentTime = &t
			}
		}

		if paymentStatus == statusPaid {
			now := time.Now()
			order.AccessGranted = true
			order.AccessStart = &now

			var video models.Video
			if result := tx.Where("id = ?", order.VideoID).First(&video); result.Error == nil {
				duration := video.AccessDurationDays
				if duration == 0 {
					duration = 30
				}
				end := now.AddDate(0, 0, duration)
				order.AccessEnd = &end
			} else {
				end := now.AddDate(0, 0, 30)
				order.AccessEnd = &end
			}
		}

		if result := tx.Save(&order); result.Error != nil {
			tx.Rollback()
			log.Printf("Failed to update video order %s: %v", notification.OrderID, result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to update order"})
		}

		if err := tx.Commit().Error; err != nil {
			log.Printf("Failed to commit video order update %s: %v", notification.OrderID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to commit transaction"})
		}

		fmt.Printf("Video payment webhook processed: order=%s status=%s\n", notification.OrderID, paymentStatus)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	}
}
