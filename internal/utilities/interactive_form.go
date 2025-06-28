package utilities

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// FormField represents a single form field
type FormField struct {
	Key        string
	Label      string
	Value      string
	IsPassword bool
	IsRequired bool
}

// InteractiveForm handles the interactive form input
type InteractiveForm struct {
	fields        []*FormField
	current       int
	originalState *term.State
}

// NewInteractiveForm creates a new interactive form
func NewInteractiveForm() *InteractiveForm {
	return &InteractiveForm{
		fields:  make([]*FormField, 0),
		current: 0,
	}
}

// AddField adds a field to the form
func (f *InteractiveForm) AddField(key, label string, isPassword, isRequired bool) {
	f.fields = append(f.fields, &FormField{
		Key:        key,
		Label:      label,
		Value:      "",
		IsPassword: isPassword,
		IsRequired: isRequired,
	})
}

// GetValues returns all field values as a map
func (f *InteractiveForm) GetValues() map[string]string {
	values := make(map[string]string)
	for _, field := range f.fields {
		values[field.Key] = field.Value
	}
	return values
}

// Run starts the interactive form
func (f *InteractiveForm) Run() error {
	if len(f.fields) == 0 {
		return fmt.Errorf("no fields defined")
	}

	// Store original terminal state
	var err error
	f.originalState, err = term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to get terminal state: %w", err)
	}

	// Ensure terminal state is restored on exit
	defer func() {
		if f.originalState != nil {
			term.Restore(int(os.Stdin.Fd()), f.originalState)
		}
		fmt.Println() // Add newline at the end
	}()

	fmt.Println("Interactive Form - Use arrow keys to navigate, Enter to confirm, Ctrl+C to quit")
	fmt.Println("Guide: ↑/↓ Navigate | Enter Confirm | Ctrl+C Quit")
	fmt.Println()

	for {
		f.displayForm()
		action := f.getUserInput()

		switch action {
		case "up":
			if f.current > 0 {
				f.current--
			}
		case "down":
			if f.current < len(f.fields)-1 {
				f.current++
			}
		case "enter":
			if f.validateCurrentField() {
				if f.current == len(f.fields)-1 {
					// Last field completed, finish form
					return nil
				}
				f.current++
			}
		case "quit":
			return fmt.Errorf("form cancelled by user")
		}
	}
}

// displayForm shows the current form state
func (f *InteractiveForm) displayForm() {
	// Clear screen (simple approach)
	fmt.Print("\033[H\033[2J")

	fmt.Println("Interactive Form - Use arrow keys to navigate, Enter to confirm, Ctrl+C to quit")
	fmt.Println("Guide: ↑/↓ Navigate | Enter Confirm | Ctrl+C Quit")
	fmt.Println()

	for i, field := range f.fields {
		indicator := " "
		if i == f.current {
			indicator = ">"
		}

		displayValue := field.Value
		if field.IsPassword && field.Value != "" {
			displayValue = strings.Repeat("*", len(field.Value))
		}

		if field.Value == "" {
			fmt.Printf("%s %s: \n", indicator, field.Label)
		} else {
			fmt.Printf("%s %s: %s\n", indicator, field.Label, displayValue)
		}
	}
}

// getUserInput handles user input and returns the action
func (f *InteractiveForm) getUserInput() string {
	// Set terminal to raw mode for single character input
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "quit"
	}

	// Always restore the raw state we just set
	defer func() {
		if oldState != nil {
			term.Restore(int(os.Stdin.Fd()), oldState)
		}
	}()

	// Read single character
	var buf [1]byte
	_, err = os.Stdin.Read(buf[:])
	if err != nil {
		return "quit"
	}

	// Check for Ctrl+C
	if buf[0] == 3 {
		return "quit"
	}

	// Check for arrow keys
	if buf[0] == 27 {
		// Read next character
		_, err = os.Stdin.Read(buf[:])
		if err != nil {
			return "quit"
		}

		if buf[0] == 91 {
			// Read the arrow key identifier
			_, err = os.Stdin.Read(buf[:])
			if err != nil {
				return "quit"
			}

			switch buf[0] {
			case 65: // Up arrow
				return "up"
			case 66: // Down arrow
				return "down"
			}
		}
	}

	// Check for Enter
	if buf[0] == 13 {
		return "enter"
	}

	return "none"
}

// validateCurrentField validates and gets input for the current field
func (f *InteractiveForm) validateCurrentField() bool {
	field := f.fields[f.current]

	// Restore terminal to original state for input
	if f.originalState != nil {
		term.Restore(int(os.Stdin.Fd()), f.originalState)
	}

	fmt.Printf("\nEnter %s: ", field.Label)

	var value string
	var err error

	if field.IsPassword {
		value, err = f.readPassword()
	} else {
		value, err = f.readString()
	}

	if err != nil {
		return false
	}

	// Trim whitespace
	value = strings.TrimSpace(value)

	// Validate required fields
	if field.IsRequired && value == "" {
		fmt.Println("This field is required. Press Enter to continue...")
		bufio.NewReader(os.Stdin).ReadString('\n')
		return false
	}

	field.Value = value
	return true
}

// readPassword reads a password input with masking
func (f *InteractiveForm) readPassword() (string, error) {
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	return string(bytePassword), nil
}

// readString reads a regular string input
func (f *InteractiveForm) readString() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(input, "\n"), nil
}
