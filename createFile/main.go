package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	targetSize = 1 * 1024 * 1024 * 1024 // 1GB
	letters    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func randomName(r *rand.Rand, length int) string {
	var sb strings.Builder
	sb.Grow(length)
	for i := 0; i < length; i++ {
		sb.WriteByte(letters[r.Intn(len(letters))])
	}
	return sb.String()
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	file, err := os.Create("data.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriterSize(file, 4*1024*1024) // 4MB buffer
	defer writer.Flush()

	header := "id,name,age\n"

	n, err := writer.WriteString(header)
	if err != nil {
		panic(err)
	}

	bytesWritten := int64(n)
	rows := 0
	for bytesWritten < targetSize {
		id := r.Int63()
		name := randomName(r, r.Intn(11)+5)
		age := r.Intn(99) + 1

		row := fmt.Sprintf("%d,%s,%d\n", id, name, age)
		n, err := writer.WriteString(row)
		if err != nil {
			panic(err)
		}
		bytesWritten += int64(n)
		rows++

		if rows%1_000_000 == 0 {
			fmt.Printf("Rows : %d, Size: %.2f GB\n", rows, float64(bytesWritten)/float64(1024*1024*1024))
		}
	}
	if err := writer.Flush(); err != nil {
		panic(err)
	}

	fmt.Printf("Finished writing %d rows, total size: %.2f GB\n", rows, float64(bytesWritten)/float64(1024*1024*1024))

}
