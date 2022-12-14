package tasks

type Task struct {
	TaskID    string `dynamodbav:"task_id" json:"task_id"`
	UserID    string `dynamodbav:"user_id" json:"user_id"`
	Content   string `dynamodbav:"content" json:"content"`
	CreatedAt string `dynamodbav:"created_at" json:"created_at"`
	IsDone    bool   `dynamodbav:"is_done" json:"is_done"`
}
