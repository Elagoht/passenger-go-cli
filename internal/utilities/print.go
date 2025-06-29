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
func PrintTable[T any](data []T, headers []string, useStderr ...bool) {
	if len(data) == 0 {
		output := os.Stdout
		if len(useStderr) > 0 && useStderr[0] {
			output = os.Stderr
		}
		fmt.Fprintln(output, "No data to display, use `passenger-go create` or `passenger-go import --file=<file>` to add data.")
		return
	}

	output := os.Stdout
	if len(useStderr) > 0 && useStderr[0] {
		output = os.Stderr
	}

	termWidth := getTerminalWidth()

	// Handle nil headers - skip header printing but still calculate widths
	if headers == nil {
		// Check if this is a key-value format (2 columns, first column is shorter)
		isKeyValue := true
		maxColumns := 0
		keyColumnMaxWidth := 0

		// Find the maximum number of columns and check if it's key-value format
		for _, row := range data {
			rowData := any(row).([]string)
			if len(rowData) > maxColumns {
				maxColumns = len(rowData)
			}
			if len(rowData) >= 2 {
				keyLength := len(rowData[0])
				if keyLength > keyColumnMaxWidth {
					keyColumnMaxWidth = keyLength
				}
			}
		}

		// If not exactly 2 columns, treat as regular table
		if maxColumns != 2 {
			isKeyValue = false
		}

		if isKeyValue {
			// Key-value format: first column fixed width, second column takes remaining space
			keyWidth := keyColumnMaxWidth
			valueWidth := termWidth - keyWidth - 3 // 3 for " | "

			// Ensure minimum widths
			if keyWidth < 8 {
				keyWidth = 8
			}
			if valueWidth < 8 {
				valueWidth = 8
			}

			// Print key-value pairs
			for _, row := range data {
				rowData := any(row).([]string)
				if len(rowData) >= 2 {
					truncatedKey := truncateString(rowData[0], keyWidth)
					truncatedValue := truncateString(rowData[1], valueWidth)
					fmt.Fprintf(output, "%-*s | %s\n", keyWidth, truncatedKey, truncatedValue)
				}
			}
		} else {
			// Regular table format - calculate max widths for all columns
			maxWidths := make([]int, maxColumns)

			// Check data rows for maximum width
			for _, row := range data {
				rowData := any(row).([]string)
				for cellIndex, cell := range rowData {
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
				availableWidth := termWidth - (len(maxWidths)-1)*3

				// Set maximum column widths based on available space
				for index := range maxWidths {
					maxColumnWidth := max(
						availableWidth/len(maxWidths),
						8,
					)

					if maxWidths[index] > maxColumnWidth {
						maxWidths[index] = maxColumnWidth
					}
				}
			}

			// Print data rows
			for _, row := range data {
				rowData := any(row).([]string)
				for index, cell := range rowData {
					if index < len(maxWidths) {
						truncatedCell := truncateString(cell, maxWidths[index])
						fmt.Fprintf(output, "%-*s", maxWidths[index], truncatedCell)
						if index < len(maxWidths)-1 {
							fmt.Fprint(output, " | ")
						}
					}
				}
				fmt.Fprintln(output)
			}
		}
		return
	}

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
		fmt.Fprintf(output, "%-*s", maxWidths[index], truncatedHeader)
		if index < len(headers)-1 {
			fmt.Fprint(output, " | ")
		}
	}
	fmt.Fprintln(output)

	// Print separator line
	totalActualWidth := 0
	for index, width := range maxWidths {
		totalActualWidth += width
		if index < len(maxWidths)-1 {
			totalActualWidth += 3 // for " | "
		}
	}
	fmt.Fprintln(output, strings.Repeat("-", totalActualWidth))

	// Print data rows
	for _, row := range data {
		rowData := any(row).([]string)
		for index, cell := range rowData {
			if index < len(maxWidths) {
				truncatedCell := truncateString(cell, maxWidths[index])
				fmt.Fprintf(output, "%-*s", maxWidths[index], truncatedCell)
				if index < len(headers)-1 {
					fmt.Fprint(output, " | ")
				}
			}
		}
		fmt.Fprintln(output)
	}
}
