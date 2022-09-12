package getitems

type Task struct {
	TaskId    string `dynamodbav:"task_id" json:"task_id"`
	UserId    string `dynamodbav:"user_id" json:"user_id"`
	Content   string `dynamodbav:"content" json:"content"`
	CreatedAt string `dynamodbav:"created_at" json:"created_at"`
	IsDone    bool   `dynamodbav:"is_done" json:"is_done"`
}
