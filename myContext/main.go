package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// 创建一个带超时的 context，2 秒后自动取消
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // 确保函数结束时释放资源

	// 模拟一个耗时任务
	select {
	case <-time.After(3 * time.Second):
		fmt.Fprintln(w, "Finished task")
	case <-ctx.Done(): // 当 context 超时时，Done 会被触发
		http.Error(w, "Request timed out", http.StatusRequestTimeout)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
