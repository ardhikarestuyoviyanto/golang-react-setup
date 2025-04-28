package task

import (
	"fmt"
	"go-auth/app/helpers"
	"go-auth/app/middleware"
	"go-auth/app/models"
	"go-auth/app/models/entity/taskEntity"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func StoreHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		decodedToken, err := middleware.DecodedToken(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success":false,"error": err.Error()})
		}
		user, _ := decodedToken["user"].(map[string]interface{})

		id, ok := user["id"].(float64)

		if !ok{
			fmt.Println("Gagal Parse Ke Float")
			return nil
		} 

		userId := uint(id)


		task := c.FormValue("task")
		taskDate := c.FormValue("taskDate")
		attachmentFile, errFile := c.FormFile("attachmentFile")

		// Validate task
		if task == "" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   "Tugas Wajib Diisi",
				"success": false,
			})
		}

		// Validate taskDate
		if taskDate == "" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   "Tanggal Wajib Diisi",
				"success": false,
			})
		}

		// Validate attachmentFile
		if errFile == nil {
			// Check file extension
			ext := filepath.Ext(attachmentFile.Filename)
			if ext != ".pdf" && ext != ".PDF" { // Case-insensitive check
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"error":   "Ekstensi File Wajib Pdf",
					"success": false,
				})
			}

			// Check file size (max 5MB)
			if attachmentFile.Size > 5*1024*1024 {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"error":   "Ukuran File Maksimal 5 mb",
					"success": false,
				})
			}
		}

		// Parse taskDate
		parsedDate, err := time.Parse("2006-01-02", taskDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   "Format Tanggal Tidak Valid",
				"success": false,
			})
		}

		// Create task data
		taskData := models.Task{
			Task:          task,
			TaskDate:      parsedDate,
			UserId : uint(userId),
			AttachmentFile: nil,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		// Save task to database
		taskResult := db.Create(&taskData)
		if taskResult.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":   "Gagal Menyimpan Tugas",
				"success": false,
			})
		}

		// Handle file upload if present
		if errFile != nil {
			// No file uploaded, return success message
			return c.JSON(http.StatusCreated, map[string]interface{}{
				"message": "Tugas Berhasil Ditambah",
				"success": true,
			})
		}

		// File handling
		dstDir := "./app/views/js/public/"
		ext := filepath.Ext(attachmentFile.Filename)
		randomFileName := uuid.New().String() + ext
		dstPath := filepath.Join(dstDir, randomFileName)

		// Open the file from the request
		src, err := attachmentFile.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":   "Gagal Membuka File",
				"success": false,
			})
		}
		defer src.Close()

		// Create the destination file
		dst, err := os.Create(dstPath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":   "Gagal Menyimpan File",
				"success": false,
			})
		}
		defer dst.Close()

		// Copy content from the uploaded file to the destination
		_, err = io.Copy(dst, src)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":   "Gagal Menyalin File",
				"success": false,
			})
		}

		// Update the task with the attachment filename
		db.Model(&models.Task{}).Where("id = ?", taskData.ID).Update("attachment_file", randomFileName)

		// Return success response
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"message": "Tugas Berhasil Ditambah",
			"success": true,
		})
	}
}


func GetAllHandler(db *gorm.DB)echo.HandlerFunc{
	return func(c echo.Context) error {
		decodedToken, err := middleware.DecodedToken(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"success":false,"error": err.Error()})
		}

		user, ok := decodedToken["user"].(map[string]interface{})

		if !ok{
			fmt.Println("Invalid user data in token")
			return nil
		}

		sortBy :=  c.QueryParam("sortBy")
		orderBy:=c.QueryParam("orderBy")
		search := c.QueryParam("search")
		startDate := c.QueryParam("startDate")
		endDate := c.QueryParam("endDate")
		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		offset,_ := strconv.Atoi(c.QueryParam("offset"))

		filter := map[string]interface{}{
			"search":search,
			"startDate": startDate,
			"endDate": endDate,
		}

		where := map[string]interface{}{
			"userId": user["id"],
		}

		var taskList []map[string]interface{} = make([]map[string]interface{}, 0)
		
		tasks := taskEntity.GetAll(db, sortBy, orderBy, filter, where, limit, offset)
		taskCount := taskEntity.CountAllList(db, filter, where)

		for i := 0; i < len(tasks); i++ {
			encryptId, _ := helpers.EncryptString(fmt.Sprintf("%v", tasks[i].ID)) 
		
			parsedDate := tasks[i].TaskDate
			taskDate := parsedDate.Format("2006-01-02")
		
			taskList = append(taskList, map[string]interface{}{
				"id":            encryptId,
				"task":          tasks[i].Task,
				"taskDate":      taskDate,
				"attachmentFile": tasks[i].AttachmentFile,
			})
		}
		

		return c.JSON(http.StatusOK, map[string]interface{}{
			"success":true,
			"message": "success get tasks",
			"data":map[string]interface{}{
				"tasks": taskList,
				"total": taskCount,
			},
		})
	}
}