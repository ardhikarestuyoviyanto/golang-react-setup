package auth

import (
	"go-auth/app/helpers"
	"go-auth/app/middleware"
	"go-auth/app/models"
	"go-auth/app/models/entity/usersEntity"
	"net/http"
	"os"
	"time"

	"github.com/dchest/captcha"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Untuk Load View Pertama Kali
func LoginView(c echo.Context) error {
	appName := os.Getenv("APP_NAME")
	data := map[string]interface{}{
		"AppName": appName,
	}
	return c.Render(http.StatusOK, "layout/index.html", data)
}

func SignUpHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")
		password := c.FormValue("password")
		confirmPassword := c.FormValue("confirmPassword")
		captchaId := c.FormValue("captchaId")
		captchaAnswer := c.FormValue("captchaAnswer")

		// Validasi nama
		if name == "" || len(name) < 3 {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   "Nama wajib diisi & minimal 3 karakter",
				"success": false,
			})
		}

		// Validasi email
		if err := helpers.ValidateEmail(email); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   err.Error(),
				"success": false,
			})
		}

		// Validasi password dan confirm password cocok
		if password != confirmPassword {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   "Password & konfirmasi password tidak cocok",
				"success": false,
			})
		}

		// Validasi captcha
		if !captcha.VerifyString(captchaId, captchaAnswer){
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   "Captcha tidak valid",
				"success": false,
			})
		}

		// Cek duplikat email
		validateEmailDuplicate := usersEntity.GetFirstByEmail(email, db)
		if validateEmailDuplicate.ID != 0 {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   "Email sudah digunakan",
				"success": false,
			})
		}

		// Hash password
		hashedPassword, _ := helpers.HashString(password)

		// Simpan user baru
		user := models.Users{
			Name:      name,
			Email:     email,
			Password:  hashedPassword,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		db.Create(&user)

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"message": "Registrasi berhasil",
			"success": true,
		})
	}
}

func SignInHandler(db *gorm.DB) echo.HandlerFunc{
	return func(c echo.Context) error {
		email := c.FormValue("email")
		password := c.FormValue("password")
		captchaId := c.FormValue("captchaId")
		captchaAnswer := c.FormValue("captchaAnswer")

		// Validasi email
		if err := helpers.ValidateEmail(email); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   err.Error(),
				"success": false,
			})
		}

		if password == ""{
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   "Password wajib diisi",
				"success": false,
			})
		}

		// Validasi captcha
		if !captcha.VerifyString(captchaId, captchaAnswer){
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   "Captcha tidak valid",
				"success": false,
			})
		}
		
		userEmail := usersEntity.GetFirstByEmail(email, db)

		if userEmail.ID == 0 {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   "Email atau Password Salah",
				"success": false,
			})
		}

		if err := helpers.CheckHash(userEmail.Password, password); err != nil{
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   "Password Salah",
				"success": false,
			})
		}
	
		resultToken, _ := middleware.GenerateJwt(userEmail)
		db.Model(&models.Users{}).Where("id = ?", userEmail.ID).Update("token", resultToken)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Berhasil Login",
			"data": map[string]interface{}{
				"token": resultToken,
				"user": map[string]interface{}{
					"id": userEmail.ID,
					"name": userEmail.Name,
					"email": userEmail.Email,
				},
			},
		})


	}
}
