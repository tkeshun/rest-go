package taskstore

//taskの生成，呼び出し，消去に関わる実装
import (
	"fmt"
	"sync"
	"time"
)

//構造体：json間の受け取り用構造体
type Task struct {
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

//in-memoryのデータベースもどき
type TaskStore struct {
	sync.Mutex
	tasks  map[int]Task
	nextId int
}

func New() *TaskStore {
	ts := &TaskStore{}
	ts.tasks = make(map[int]Task)
	ts.nextId = 0
	return ts
}

//CreateTask で新しいタスクを作り，TaskStoreに格納する
func (ts *TaskStore) CreateTask(text string, tags []string, due time.Time) int {
	//更新処理，複数のアクセスが重ならないようにlockをかけて制限する
	ts.Lock()
	defer ts.Unlock()
	task := Task{
		Id:   ts.nextId,
		Text: text,
		Due:  due}
	task.Tags = make([]string, len(tags))
	copy(task.Tags, tags)
	ts.tasks[ts.nextId] = task
	ts.nextId++
	return task.Id
}

func (ts *TaskStore) GetTask(id int) (Task, error) {
	ts.Lock()
	defer ts.Unlock()

	t, ok := ts.tasks[id]
	if ok {
		return t, nil
	} else {
		return Task{}, fmt.Errorf("task with id=%d not found", id)
	}
}

func (ts *TaskStore) DeleteTask(id int) error {
	ts.Lock()
	defer ts.Unlock()

	if _, ok := ts.tasks[id]; !ok {
		return fmt.Errorf("task with id=%d not found", id)
	}

	delete(ts.tasks, id)
	return nil
}

func (ts *TaskStore) DeleteAllTasks() error {
	ts.Lock()
	defer ts.Unlock()

	ts.tasks = make(map[int]Task)
	return nil
}

//GetAllTasksは格納されたすべてのタスクを返す
func (ts *TaskStore) GetAllTasks() []Task {
	ts.Lock()
	defer ts.Unlock()
	allTasks := make([]Task, 0, len(ts.tasks))
	//taskの格納
	for _, task := range ts.tasks {
		allTasks = append(allTasks, task)
	}
	return allTasks
}

//GetTaskByTag はタグ付けされたタスクを返す
func (ts *TaskStore) GetTasksByTag(tag string) []Task {
	ts.Lock()
	defer ts.Unlock()
	var tasks []Task
taskloop:
	for _, task := range ts.tasks {
		for _, taskTag := range task.Tags {
			if taskTag == tag {
				tasks = append(tasks, task)
				continue taskloop
			}
		}

	}
	return tasks
}

// GetTasksByDueDateはTaskの作成時刻を返す
func (ts *TaskStore) GetTasksByDueDate(year int, month time.Month, day int) []Task {
	ts.Lock()
	defer ts.Unlock()
	var tasks []Task
	for _, task := range ts.tasks {
		y, m, d := task.Due.Date()
		if y == year && m == month && d == day {
			tasks = append(tasks, task)
		}
	}
	return tasks
}
