package cron

import (
	"fmt"
	"testing"

	"github.com/zly-app/zapp"
)

func TestTask(t *testing.T) {
	task := NewTask("test", "@every 1s", true, func(ctx IContext) (err error) {
		fmt.Println("触发")
		return nil
	})
	app := zapp.NewApp("cron")
	err := task.Trigger(newContext(app, task), nil)
	if err != nil {
		t.Fatal(err)
	}
}
