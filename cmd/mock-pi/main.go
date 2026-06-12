package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Command struct {
	Type string `json:"type"`
}

func main() {
	// Emit initial queue on startup
	emit("queue_update", map[string]interface{}{
		"steering": []string{"Audit error paths", "Add tests"},
		"follow_up": []string{"Refactor the auth module"},
	})

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		var cmd Command
		if err := json.Unmarshal([]byte(line), &cmd); err != nil {
			continue
		}
		
		if cmd.Type == "prompt" {
			// Simulate steps declaration
			emit("declare_steps", map[string]interface{}{
				"steps": []map[string]interface{}{
					{"id": 1, "description": "Read auth module"},
					{"id": 2, "description": "Propose refactor"},
					{"id": 3, "description": "Apply changes"},
				},
			})

			time.Sleep(500 * time.Millisecond)

			// Step 1 runs
			emit("tool_execution_start", map[string]string{"tool_name": "read_file"})
			emit("message_update", map[string]string{"delta": "Thinking...\n"})
			time.Sleep(500 * time.Millisecond)
			emit("tool_execution_end", map[string]string{"tool_name": "read_file", "result": "done"})

			// Step 2 runs
			emit("tool_execution_start", map[string]string{"tool_name": "propose_changes"})
			emit("message_update", map[string]string{"delta": "Here "})
			time.Sleep(200 * time.Millisecond)
			emit("message_update", map[string]string{"delta": "is "})
			time.Sleep(200 * time.Millisecond)
			emit("message_update", map[string]string{"delta": "my "})
			time.Sleep(200 * time.Millisecond)
			emit("message_update", map[string]string{"delta": "response.\n"})
			emit("tool_execution_end", map[string]string{"tool_name": "propose_changes", "result": "done"})

			// Step 3 will just stay pending in our simulation
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
	}
}

func emit(typ string, payload any) {
	b, _ := json.Marshal(payload)
	out := fmt.Sprintf(`{"type": "%s", "payload": %s}`, typ, string(b))
	fmt.Println(out)
}
