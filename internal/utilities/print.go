package utilities

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	if maxLength <= 3 {
		return s[:maxLength]
	}
	return s[:maxLength-3] + "..."
}

func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80 // Most terminals default to 80x24
	}
	return width
}

// calculate the max width of each column and print the table
func PrintTable[T any](data []T, headers []string) {
	if len(data) == 0 {
		fmt.Println("No data to display")
		return
	}

	termWidth := getTerminalWidth()

	// Calculate initial max widths
	maxWidths := make([]int, len(headers))
	for index, header := range headers {
		maxWidths[index] = len(header)
	}

	// Check data rows for maximum width
	for index := range data {
		row := any(data[index]).([]string)
		for cellIndex, cell := range row {
			if cellIndex < len(maxWidths) && len(cell) > maxWidths[cellIndex] {
				maxWidths[cellIndex] = len(cell)
			}
		}
	}

	// Calculate total width needed including separators
	totalWidth := 0
	for index, width := range maxWidths {
		totalWidth += width
		if index < len(maxWidths)-1 {
			totalWidth += 3 // for " | "
		}
	}

	// If table is too wide, proportionally reduce column widths
	if totalWidth > termWidth {
		// Reserve space for separators
		availableWidth := termWidth - (len(headers)-1)*3

		// Set maximum column widths based on available space
		for index := range maxWidths {
			maxColumnWidth := max(
				availableWidth/len(headers),
				8,
			)

			if maxWidths[index] > maxColumnWidth {
				maxWidths[index] = maxColumnWidth
			}
		}
	}

	// Print headers with proper padding
	for index, header := range headers {
		truncatedHeader := truncateString(header, maxWidths[index])
		fmt.Printf("%-*s", maxWidths[index], truncatedHeader)
		if index < len(headers)-1 {
			fmt.Print(" | ")
		}
	}
	fmt.Println()

	// Print separator line
	totalActualWidth := 0
	for index, width := range maxWidths {
		totalActualWidth += width
		if index < len(maxWidths)-1 {
			totalActualWidth += 3 // for " | "
		}
	}
	fmt.Println(strings.Repeat("-", totalActualWidth))

	// Print data rows
	for _, row := range data {
		rowData := any(row).([]string)
		for index, cell := range rowData {
			if index < len(maxWidths) {
				truncatedCell := truncateString(cell, maxWidths[index])
				fmt.Printf("%-*s", maxWidths[index], truncatedCell)
				if index < len(headers)-1 {
					fmt.Print(" | ")
				}
			}
		}
		fmt.Println()
	}
}
