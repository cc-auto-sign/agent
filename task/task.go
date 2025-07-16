// Package task 提供任务执行相关功能
package task

// Task 任务接口
type Task interface {
	Execute() (string, error)
}
