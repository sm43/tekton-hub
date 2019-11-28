package models

import (
	"log"
	"strconv"
)

// Task is a database model representing task data
type Task struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Downloads   int      `json:"downloads"`
	Rating      float64  `json:"rating"`
	Github      string   `json:"github"`
	Tags        []string `json:"tags"`
}

func addTask(task *Task) {
	sqlStatement := `
	INSERT INTO TASK (NAME,DESCRIPTION,DOWNLOADS,RATING,GITHUB)
	VALUES ($1, $2, $3, $4, $5)`
	_, err := DB.Exec(sqlStatement, task.Name, task.Description, task.Downloads, task.Rating, task.Github)
	if err != nil {
		log.Println(err)
	}
}

// GetAllTasks will return all the tasks
func GetAllTasks() []Task {
	tasks := []Task{}
	sqlStatement := `
	SELECT * FROM TASK`
	rows, err := DB.Query(sqlStatement)
	defer rows.Close()
	for rows.Next() {
		task := Task{}
		err = rows.Scan(&task.ID, &task.Name, &task.Description, &task.Downloads, &task.Rating, &task.Github)
		if err != nil {
			log.Println(err)
		}
		tasks = append(tasks, task)
	}
	sqlStatement = `SELECT T.ID,TG.NAME FROM TAG TG JOIN TASK_TAG TT ON TT.TAG_ID=TG.ID JOIN TASK T ON T.ID=TT.TASK_ID`
	rows, err = DB.Query(sqlStatement)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var tag string
		var taskID int
		err := rows.Scan(&taskID, &tag)
		if err != nil {
			log.Println(err)
		}
		tasks[taskID-1].Tags = append(tasks[taskID-1].Tags, tag)
	}
	return tasks
}

// GetTaskWithName returns a task with requested ID
func GetTaskWithName(name string) Task {
	task := Task{}
	id, err := strconv.Atoi(name)
	if err != nil {
		log.Fatalln(err)
	}
	var taskTagMap map[int][]string
	taskTagMap = make(map[int][]string)
	taskTagMap = getTaskTagMap()
	sqlStatement := `
	SELECT * FROM TASK WHERE ID=$1;`
	err = DB.QueryRow(sqlStatement, id).Scan(&task.ID, &task.Name, &task.Description, &task.Downloads, &task.Rating, &task.Github)
	if err != nil {
		return Task{}
	}
	task.Tags = taskTagMap[task.ID]
	return task
}

// GetTaskNameFromID returns name from given ID
func GetTaskNameFromID(taskID string) string {
	id, err := strconv.Atoi(taskID)
	if err != nil {
		log.Println(err)
	}
	sqlStatement := `SELECT NAME FROM TASK WHERE ID=$1`
	var taskName string
	err = DB.QueryRow(sqlStatement, id).Scan(&taskName)
	if err != nil {
		return ""
	}
	log.Println(taskName)
	return taskName
}

// IncrementDownloads will increment the number of downloads
func IncrementDownloads(taskID string) {
	id, err := strconv.Atoi(taskID)
	if err != nil {
		log.Println(err)
	}
	sqlStatement := `UPDATE TASK SET DOWNLOADS = DOWNLOADS + 1 WHERE ID=$1`
	_, err = DB.Exec(sqlStatement, id)
	if err != nil {
		log.Println(err)
	}
}

func updateAverageRating(taskID int, rating float64) {
	sqlStatement := `UPDATE TASK SET RATING=$2 WHERE ID=$1`
	_, err := DB.Exec(sqlStatement, taskID, rating)
	if err != nil {
		log.Println(err)
	}
}
