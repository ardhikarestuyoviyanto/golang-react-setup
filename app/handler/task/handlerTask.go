package task

import (
	"fmt"
	"go-auth/app/helpers"
	"go-auth/app/models"
	"go-auth/app/models/entity/taskEntity"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func DestroyHandler(db *gorm.DB)echo.HandlerFunc{
	return func(c echo.Context) error {
		taskIdHash := c.Param("taskId")
		taskIdStr, err := helpers.DecryptString(taskIdHash)
		if err != nil{
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   "taskId tidak valid",
				"success": false,
			})
		}
		
		taskId, err := strconv.Atoi(taskIdStr)

		if err != nil{
			return fmt.Errorf("Gagal Konversi taskId ke int")
		}

		var taskModel models.Task
		db.First(&taskModel, taskId)
		if taskModel.ID == 0{
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   "task tidak ada",
				"success": false,
			})
		}

		db.Delete(&taskModel)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message":   "Berhasil Hapus",
			"success": true,
		})
		
	}
}

func UpdateHandler(db *gorm.DB)echo.HandlerFunc{
	return func(c echo.Context) error {
		taskIdHash := c.Param("taskId")
		taskIdStr, err := helpers.DecryptString(taskIdHash)
		if err != nil{
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   "taskId tidak valid",
				"success": false,
			})
		}
		
		taskId, err := strconv.Atoi(taskIdStr)

		if err != nil{
			return fmt.Errorf("Gagal Konversi taskId ke int")
		}

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
			// Validasi file
			err := helpers.ValidateExtFile([]string{".pdf", ".PDF"}, attachmentFile, 5)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"error":   err.Error(),
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

		var taskModel models.Task
		db.First(&taskModel, taskId)
		if taskModel.ID == 0{
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   "task tidak ada",
				"success": false,
			})
		}

		taskModel.Task = task
		taskModel.TaskDate = parsedDate
		if errFile == nil{
			// Update File
			dstDir := "./app/views/storage/file/"
			fileName, err := helpers.UploadFile(dstDir, attachmentFile)
			if err != nil{
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"error":  err.Error(),
					"success": false,
				})
			}
			// Update
			taskModel.AttachmentFile = &fileName
		}

		db.Save(&taskModel)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":   "task berhasil update",
			"success": true,
		})
	}
}

func GetHandler(db *gorm.DB) echo.HandlerFunc{
	return func(c echo.Context) error {
		taskIdHash := c.Param("taskId")
		taskIdStr, err := helpers.DecryptString(taskIdHash)
		if err != nil{
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   "taskId tidak valid",
				"success": false,
			})
		}
		
		taskId, err := strconv.Atoi(taskIdStr)

		if err != nil{
			return fmt.Errorf("Gagal Konversi taskId ke int")
		}

		var task models.Task
		db.First(&task, taskId)

		if task.ID == 0{
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   "Task tidak ada",
				"success": false,
			})
		}

		parsedDate := task.TaskDate
		taskDate := parsedDate.Format("2006-01-02")

		result := map[string]interface{}{
			"data":map[string]interface{}{
				"id":taskIdHash,
				"task": task.Task,
				"taskDate": taskDate,
				"attachmentFile": task.AttachmentFile,
			},
			"success": true,
		}

		return c.JSON(http.StatusOK, result)
	}
}

func StoreHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(map[string]interface{})

		if !ok{
			return fmt.Errorf("user ga ada")
		}

		id, ok := user["id"].(float64)

		if !ok{
			return fmt.Errorf("Gagal Parse Ke Float")
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
			// Validasi file
			err := helpers.ValidateExtFile([]string{".pdf", ".PDF"}, attachmentFile, 5)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"error":   err.Error(),
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

		dstDir := "./app/views/storage/file/"
		fileName, err := helpers.UploadFile(dstDir, attachmentFile)
		if err != nil{
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":  err.Error(),
				"success": false,
			})
		}

		// Update the task with the attachment filename
		db.Model(&models.Task{}).Where("id = ?", taskData.ID).Update("attachment_file", fileName)

		// Return success response
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"message": "Tugas Berhasil Ditambah",
			"success": true,
		})
	}
}


func GetAllHandler(db *gorm.DB)echo.HandlerFunc{
	return func(c echo.Context) error {
		user, ok := c.Get("user").(map[string]interface{})
		if !ok{
			return fmt.Errorf("user ga ada")
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
			encryptId, err := helpers.EncryptString(fmt.Sprintf("%v", tasks[i].ID)) 
			
			if err != nil{
				log.Print(err.Error())
			}

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