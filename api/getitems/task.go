package getitems

type Task struct {
	TaskId    string `dynamodbav:"task_id"`
	Content   string `dynamodbav:"content"`
	CreatedAt string `dynamodbav:"created_at"`
	IsDone    bool   `dynamodbav:"is_done"`
}
