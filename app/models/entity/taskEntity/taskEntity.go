package taskEntity

import (
	"fmt"
	"go-auth/app/models"

	"gorm.io/gorm"
)

func GetAll(
	db *gorm.DB,
	sortBy string,
	orderBy string,
	filter map[string]interface{},
	where map[string]interface{},
	limit int,
	offset int,
) []models.Task {

	// start query
    query := "SELECT * FROM tasks WHERE deleted_at IS NULL AND user_id = ?"

	if filter["search"] != "" {
		query += " AND task LIKE ?"
	}

	if filter["startDate"] != "" && filter["endDate"] != "" {
		query += " AND DATE(task_date) BETWEEN ? AND ?"
	}

	if sortBy != "" && orderBy != "" {
        query += " ORDER BY " + sortBy + " " + orderBy
    }

	query += " LIMIT ? OFFSET ?"

	// start args
	args := []interface{}{
        where["userId"], 
    }

	if filter["search"] != "" {
        args = append(args, "%" + filter["search"].(string) + "%")
    }

	if filter["startDate"] != "" && filter["endDate"] != "" {
        args = append(args, filter["startDate"], filter["endDate"])
    }

	args = append(args, limit, offset)

    // Execute the raw query
    var taskResults []models.Task
    if err := db.Raw(query, args...).Scan(&taskResults).Error; err != nil {
        fmt.Println("Error executing query:", err)
        return nil
    }

	return taskResults
}

func CountAllList(db *gorm.DB, filter map[string]interface{},where map[string]interface{})int64{
	
	// start query
	query := "SELECT COUNT(*) FROM tasks WHERE deleted_at IS NULL AND user_id = ?"
		
	if filter["search"] != "" {
		query += " AND task LIKE ?"
	}

	if filter["startDate"] != "" && filter["endDate"] != "" {
		query += " AND DATE(task_date) BETWEEN ? AND ?"
	}

	// start args
	args := []interface{}{
        where["userId"], 
    }

	if filter["search"] != "" {
        args = append(args, "%" + filter["search"].(string) + "%")
    }

	if filter["startDate"] != "" && filter["endDate"] != "" {
        args = append(args, filter["startDate"], filter["endDate"])
    }

	var total int64
	if err := db.Raw(query, args...).Scan(&total).Error; err != nil{
		fmt.Println("Error executing query:", err)
        return 0
	}

	return total
}