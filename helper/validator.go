package helper

import (
	"encoding/json"

	"progas-wms-be/constant"
	"progas-wms-be/dto"
	"progas-wms-be/global"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var validate = validator.New()

func ValidateStruct(request any) []dto.ValidationError {
	var validationErrors []dto.ValidationError
	err := validate.Struct(request)
	if err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, err := range errs {
				var element dto.ValidationError
				element.FailedField = err.StructNamespace()
				element.Tag = err.Tag()
				element.Value = err.Param()
				validationErrors = append(validationErrors, element)
			}
		}
	}

	if validationErrors != nil {
		return validationErrors
	}

	return nil
}

func ValidateBody(c fiber.Ctx, request any) global.ErrorResponse {
	if err := c.Bind().Body(request); err != nil {
		return global.BadRequestError("Validation Error")
	}

	if err := ValidateStruct(request); err != nil {
		return global.BadRequestErrorWithData("Validation Error", err, constant.ErrValidationError)
	}

	return nil
}

func ValidateParam(c fiber.Ctx, param any) global.ErrorResponse {
	if err := c.Bind().URI(param); err != nil {
		return global.BadRequestErrorWithData("Validation Error", err, constant.ErrValidationError)
	}

	if err := ValidateStruct(param); err != nil {
		return global.BadRequestErrorWithData("Validation Error", err, constant.ErrValidationError)
	}

	return nil
}

func ValidateHeader(c fiber.Ctx, header any) global.ErrorResponse {
	if err := c.Bind().Header(header); err != nil {
		return global.BadRequestErrorWithData("Validation Error", err, constant.ErrValidationError)
	}

	if err := ValidateStruct(header); err != nil {
		return global.BadRequestErrorWithData("Validation Error", err, constant.ErrValidationError)
	}

	return nil
}

func GetUserID(c fiber.Ctx) (string, global.ErrorResponse) {
	userId := c.Get("X-UserId")
	if userId == "" {
		return "", global.ForbiddenError()
	}
	return userId, nil
}

func ValidateQuery(c fiber.Ctx, query any) global.ErrorResponse {
	if err := c.Bind().Query(query); err != nil {
		return global.BadRequestErrorWithData("Validation Error", err, constant.ErrValidationError)
	}

	if err := ValidateStruct(query); err != nil {
		return global.BadRequestErrorWithData("Validation Error", err, constant.ErrValidationError)
	}

	return nil
}

func ValidateMessage(body []byte, payload any) global.ErrorResponse {
	if err := json.Unmarshal(body, payload); err != nil {
		return global.BadRequestErrorWithData("Validation Error", err, constant.ErrValidationError)
	}

	if err := ValidateStruct(payload); err != nil {
		return global.BadRequestErrorWithData("Validation Error", err, constant.ErrValidationError)
	}

	return nil
}
