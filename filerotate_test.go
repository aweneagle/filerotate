package filerotate

import (
	"strconv"
	"testing"
	"time"
)

func TestRotate(t *testing.T) {
	files := make(map[string]int)
	r := Rotate{
		FromFile: "./rotate_test.log",
		ToDir:    "./",
		ToFile: func() string {
			sec := time.Now().Unix() % 10
			f := "rotate_test." + strconv.FormatInt(sec, 10) + ".log"
			files[f] = 1
			return f
		},
		Permission: 0666,
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 1000; j++ {
			go func() {
				r.Write([]byte("hello world"))
			}()
		}
		time.Sleep(1 * time.Second)
	}
}
