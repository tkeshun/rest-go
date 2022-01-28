package taskstore

import (
	"testing"
	"time"
)

func TestCreateAndGet(t *testing.T) {
	ts := New()
	inputText := "テスト勉強"
	id := ts.CreateTask(inputText, nil, time.Now())
	//タスクをIDでゲットできるか試す
	task, err := ts.GetTask(id)
	if err != nil {
		t.Fatal(err)
	}
	//idが一致しているか
	if task.Id != id {
		t.Errorf("got task.Id=%d, id=%d", task.Id, id)
	}
	//打ち込んだテキストが一致しているか
	if task.Text != inputText {
		t.Errorf("got Text=%v, want %v", task.Text, "Hola")
	}
	allTasks := ts.GetAllTasks()
	//余分なタスクができてないか，idが順序と対応しているか
	if len(allTasks) != 1 || allTasks[0].Id != id {
		t.Errorf("got len(allTasks)=%d, allTasks[0].Id=%d; want 1, %d", len(allTasks), allTasks[0].Id, id)
	}
	//存在しないidを持つタスクがエラーなしで実行されないか
	_, err = ts.GetTask(id + 1)
	if err == nil {
		t.Fatal("got nil, want Error")
	}
	inputText = "テスト勉強2"
	ts.CreateTask(inputText, nil, time.Now())
	allTasks2 := ts.GetAllTasks()
	if len(allTasks2) != 2 {
		t.Errorf("got len(allTasks2)=%d; want 2", len(allTasks2))
	}
	
}
