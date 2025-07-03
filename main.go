package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"go-reloaded/processor"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: textprocessor <input_file> <output_file>")
		os.Exit(1)
	}

	inputFile, outputFile := os.Args[1], os.Args[2]

	if !strings.HasSuffix(inputFile, ".txt") || !strings.HasSuffix(outputFile, ".txt") {
		fmt.Println("Ошибка: оба файла должны иметь расширение .txt")
	}

	// Открытие входного файла
	input, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Failed to open input file: %v", err)
	}
	defer input.Close()

	// Создание выходного файла
	output, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer output.Close()

	// Обработка файла построчно
	scanner := bufio.NewScanner(input)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, processor.ProcessLine(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	// Запись результата
	for i, line := range lines {
		if i < len(lines)-1 {
			if _, err := output.WriteString(line + "\n"); err != nil {
				log.Fatalf("Error writing to output file: %v", err)
			}
		} else {
			if _, err := output.WriteString(line); err != nil {
				log.Fatalf("Error writing to output file: %v", err)
			}
		}
	}
}
