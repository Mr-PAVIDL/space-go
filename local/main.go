package main

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"log"
	"os"
	"space-go/internal/commander/packer"
	"space-go/internal/model"
	"strconv"
	"strings"
	"time"
)

func main() {
	file, err := os.ReadFile("C:\\Users\\ischenkx\\code\\eden\\space-go\\local\\garbage.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	var garbage map[string]model.Garbage
	if err := json.Unmarshal(file, &garbage); err != nil {
		log.Fatal(err)
	}

	t1 := time.Now()

	data, _ := packer.Pack(garbage, 8, 11)

	fmt.Println(time.Since(t1))

	fmt.Println("[")
	for _, row := range data {
		fmt.Printf("[%s],\n", strings.Join(lo.Map(row, func(item int, _ int) string {
			return strconv.Itoa(item)
		}), ", "))
	}
	fmt.Println("]")
}
