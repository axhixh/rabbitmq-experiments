package common

import "fmt"
import "time"

type Generator struct {
	Key   string
	Color int
}

func (g *Generator) Generate(channel chan Message, maxCount int) {
	for i := 0; i < maxCount; i++ {
		channel <- Message{Key: g.Key, Body: body(g.Key, g.Color, i)}
		time.Sleep(2 * time.Second)
	}
	channel <- Message{Key: g.Key, Body: "done"}
}

func body(key string, color int, counter int) string {
	return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m:%02d", color, key, counter)
}
