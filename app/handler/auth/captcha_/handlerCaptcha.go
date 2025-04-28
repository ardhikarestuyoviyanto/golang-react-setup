package captcha_

import (
	"net/http"

	"github.com/dchest/captcha"
	"github.com/labstack/echo/v4"
)

func GetCaptcha(c echo.Context) error{
	captchaId := c.Param("captchaId")
	c.Response().Header().Set("Content-Type", "image/png")
	captcha.WriteImage(c.Response().Writer, captchaId, 240, 80)
	return nil
}

func GenerateCaptcha(c echo.Context)error{
	captchaId := captcha.New()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "success get captchaId",
		"data":map[string]interface{} {
			"captchaId":captchaId,
		},
	})
}
