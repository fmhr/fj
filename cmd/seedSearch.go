package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// 0-99 100-199 200-299 300-399 400-499 500-599 600-699 700-799 800-900
func seedSorting() {
	var seeds [10][]int
	for i := 0; i <= 9999; i++ {
		filename := fmt.Sprintf("tools/in/%04d.txt", i)
		file, err := os.Open(filename)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			continue
		}

		scanner := bufio.NewScanner(file)
		scanner.Scan()
		line := scanner.Text()
		file.Close()
		// スペースで区切られた３つ目の要素を取得
		str := strings.Split(line, " ")[2]
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Error converting string to int: %v\n", err)
			continue
		}
		if num < 100 {
			seeds[0] = append(seeds[0], i)
		} else if num < 200 {
			seeds[1] = append(seeds[1], i)
		} else if num < 300 {
			seeds[2] = append(seeds[2], i)
		} else if num < 400 {
			seeds[3] = append(seeds[3], i)
		} else if num < 500 {
			seeds[4] = append(seeds[4], i)
		} else if num < 600 {
			seeds[5] = append(seeds[5], i)
		} else if num < 700 {
			seeds[6] = append(seeds[6], i)
		} else if num < 800 {
			seeds[7] = append(seeds[7], i)
		} else if num <= 900 {
			seeds[8] = append(seeds[8], i)
		} else {
			log.Println("Error: seed is out of range:", num)
		}
	}
	for i := 0; i < 9; i++ {
		saveToFile(fmt.Sprintf("seeds-%d.txt", i), seeds[i])
	}
}

func saveToFile(filename string, data []int) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for _, v := range data {
		file.WriteString(strconv.Itoa(v))
		file.WriteString("\n")
	}
}
