package utils

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetPagination(c *fiber.Ctx) (int, int, int) {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	return page, perPage, offset
}

func CalcLastPage(total int64, perPage int) int {
	lp := int(math.Ceil(float64(total) / float64(perPage)))
	if lp < 1 {
		lp = 1
	}
	return lp
}

func GetUserID(c *fiber.Ctx) (uuid.UUID, bool) {
	uid, ok := c.Locals("user_id").(uuid.UUID)
	return uid, ok
}

func GetUserRole(c *fiber.Ctx) string {
	role, _ := c.Locals("user_role").(string)
	return role
}

func IsAdmin(c *fiber.Ctx) bool {
	return GetUserRole(c) == "admin"
}

func ValidateStatus(status string, allowed []string) bool {
	for _, a := range allowed {
		if status == a {
			return true
		}
	}
	return false
}
