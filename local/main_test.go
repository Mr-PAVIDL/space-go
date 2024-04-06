package local

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
	"testing"
	"time"
)

var masks = []string{
	`
	!!
	!!
	`,
	`
	!!
	`,
	`
	!!!!
	#!##
	#!##
	#!!#
	`,
	`
	!!!
	#!#
	#!!
	`,
	`
	#!##
	!!!!
	!###
	!###
	`,
	`
	!!!
	!##
	!!!
	#!#
	`,
	`
	##!#
	##!!
	##!#
	!!!!
	`,
	`
	!!!!
	!##!
	!##!
	!!#!
	`,
	`
	!!!#
	#!!!
	`,
	`
	!!!
	!#!
	!#!
	!#!
	`,
	`
	!!!!
	!#!#
	##!!
	####
	`,
	`
	!###
	!!!!
	`,
	`
	!#
	!!
	!!
	!#
	`,
	`
	!#!
	!!!
	!##
	!##
	`,
	`
	!#
	!!
	!!
	#!
	`,
	`
	!!!
	`,
	`
	!!!
	!##
	!##
	!!!
	`,
	`
	!!##
	#!!!
	###!
	####
	`,
	`
	#!
	#!
	!!
	!!
	`,
	`
	!!!
	`,
	`
	#!##
	!!!!
	`,
	`
	##!!
	!!!#
	!###
	!!!!
	`,
	`
	###!
	!#!!
	!#!#
	!!!#
	`,
}

func m2p(m string, id int) *packer.Polyomino {
	mat := model.Matrix{}
	for _, line := range strings.Split(m, "\n") {
		var row []int
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		for _, c := range line {
			if c == '!' {
				row = append(row, 1)
			} else {
				row = append(row, 0)
			}
		}
		mat = append(mat, row)
	}

	return &packer.Polyomino{
		Matrix:    mat,
		ID:        id,
		GarbageID: "",
	}
}

func BenchmarkName(b *testing.B) {
	file, err := os.ReadFile("C:\\Users\\ischenkx\\code\\eden\\space-go\\local\\garbage.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	var garbage map[string]model.Garbage
	if err := json.Unmarshal(file, &garbage); err != nil {
		log.Fatal(err)
	}

	polyminoes := []*packer.Polyomino{}
	for i, mask := range masks {
		polyminoes = append(polyminoes, m2p(mask, 20+i*20))
	}

	c := 15
	w, h := 8, 11

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var ps []*packer.Polyomino

		for j := 0; j < c; j++ {
			ps = append(ps, lo.Sample(polyminoes))
		}
		_ = packer.BoostedRawPack(ps, 0, 10000, w, h)
	}
}
func Bench(b testing.B) {
	file, err := os.ReadFile("C:\\Users\\ischenkx\\code\\eden\\space-go\\local\\garbage.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	var garbage map[string]model.Garbage
	if err := json.Unmarshal(file, &garbage); err != nil {
		log.Fatal(err)
	}

	polyminoes := []*packer.Polyomino{}
	for i, mask := range masks {
		polyminoes = append(polyminoes, m2p(mask, 20+i*20))
	}

	c := 15
	w, h := 8, 11
	for i := 0; i < 10; i++ {
		var ps []*packer.Polyomino

		for j := 0; j < c; j++ {
			ps = append(ps, lo.Sample(polyminoes))
		}

		fmt.Println("Starting")
		t1 := time.Now()
		data := packer.BoostedRawPack(ps, 0, 10000, w, h)
		fmt.Println("Elapsed:", time.Since(t1))

		fmt.Println("[")
		for _, row := range data {
			fmt.Printf("[%s],\n", strings.Join(lo.Map(row, func(item int, _ int) string {
				return strconv.Itoa(item)
			}), ", "))
		}
		fmt.Println("]")
		fmt.Println("----------------------")
	}

}
