package middleware

import (
	"progas-wms-be/global"
	"progas-wms-be/repository"

	"github.com/gofiber/fiber/v3"
)

func Authorize(rbacRepo repository.RbacRepository, permissionKey string) fiber.Handler {
	return func(c fiber.Ctx) error {
		roleId, _ := c.Locals("role_id").(string)
		if roleId == "" {
			return global.UnauthorizedError().ToResponse(c)
		}

		isSuperAdmin, errRes := rbacRepo.IsSuperAdmin(roleId)
		if errRes != nil {
			return errRes.ToResponse(c)
		}
		if isSuperAdmin {
			return c.Next()
		}

		allowed, errRes := rbacRepo.HasPermission(roleId, permissionKey)
		if errRes != nil {
			return errRes.ToResponse(c)
		}
		if !allowed {
			return global.ForbiddenError().ToResponse(c)
		}

		return c.Next()
	}
}
