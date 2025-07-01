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
	Key          string
	Label        string
	Value        string
	DefaultValue string
	IsPassword   bool
	IsRequired   bool
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
func (form *InteractiveForm) AddField(key, label string, isPassword, isRequired bool) {
	form.fields = append(form.fields, &FormField{
		Key:          key,
		Label:        label,
		Value:        "",
		DefaultValue: "",
		IsPassword:   isPassword,
		IsRequired:   isRequired,
	})
}

// AddFieldWithDefault adds a field to the form with a default value
func (form *InteractiveForm) AddFieldWithDefault(key, label, defaultValue string, isPassword, isRequired bool) {
	form.fields = append(form.fields, &FormField{
		Key:          key,
		Label:        label,
		Value:        defaultValue,
		DefaultValue: defaultValue,
		IsPassword:   isPassword,
		IsRequired:   isRequired,
	})
}

// GetValues returns all field values as a map
func (form *InteractiveForm) GetValues() map[string]string {
	values := make(map[string]string)
	for _, field := range form.fields {
		values[field.Key] = field.Value
	}
	return values
}

// Run starts the interactive form
func (form *InteractiveForm) Run() error {
	if len(form.fields) == 0 {
		return fmt.Errorf("no fields defined")
	}

	// Store original terminal state
	var err error
	form.originalState, err = term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to get terminal state: %w", err)
	}

	// Ensure terminal state is restored on exit
	defer func() {
		if form.originalState != nil {
			term.Restore(int(os.Stdin.Fd()), form.originalState)
		}
		fmt.Println() // Add newline at the end
	}()

	fmt.Println("Guide: ↑/↓ Navigate | Enter Confirm | Ctrl+D Clear Field | Ctrl+S Save & Quit | Ctrl+C Quit")
	fmt.Println()

	for {
		form.displayForm()
		action := form.getUserInput()

		switch action {
		case "up":
			if form.current > 0 {
				form.current--
			}
		case "down":
			if form.current < len(form.fields)-1 {
				form.current++
			}
		case "enter":
			if form.validateCurrentField("") {
				if form.current == len(form.fields)-1 {
					// Last field completed, finish form
					return nil
				}
				form.current++
			}
		case "clear":
			if form.validateCurrentField("clear") {
				if form.current == len(form.fields)-1 {
					return nil
				}
				form.current++
			}
		case "savequit":
			return nil
		case "quit":
			return fmt.Errorf("form cancelled by user")
		}
	}
}

// displayForm shows the current form state
func (form *InteractiveForm) displayForm() {
	// Clear screen (simple approach)
	fmt.Print("\033[H\033[2J")

	fmt.Println("Guide: ↑/↓ Navigate | Enter Confirm | Ctrl+D Clear Field | Ctrl+S Save & Quit | Ctrl+C Quit")
	fmt.Println()

	for i, field := range form.fields {
		indicator := " "
		if i == form.current {
			indicator = ">"
		}

		// Use Value for display (empty string means cleared field)
		displayValue := field.Value

		if field.IsPassword && displayValue != "" {
			displayValue = strings.Repeat("*", len(displayValue))
		}

		if displayValue == "" {
			fmt.Printf("%s %s: \n", indicator, field.Label)
		} else {
			fmt.Printf("%s %s: %s\n", indicator, field.Label, displayValue)
		}
	}
}

// getUserInput handles user input and returns the action
func (form *InteractiveForm) getUserInput() string {
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

	// Check for Ctrl+C (quit)
	if buf[0] == 3 {
		return "quit"
	}
	// Check for Ctrl+D (clear)
	if buf[0] == 4 {
		return "clear"
	}
	// Check for Ctrl+S (save and quit)
	if buf[0] == 19 {
		return "savequit"
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
func (form *InteractiveForm) validateCurrentField(action string) bool {
	field := form.fields[form.current]

	// Restore terminal to original state for input
	if form.originalState != nil {
		term.Restore(int(os.Stdin.Fd()), form.originalState)
	}

	// Show current value for non-password fields
	if field.DefaultValue != "" && !field.IsPassword {
		fmt.Printf(
			"\nEnter %s (current: %s): ",
			field.Label,
			field.DefaultValue,
		)
	} else {
		fmt.Printf("\nEnter %s: ", field.Label)
	}

	if action == "clear" {
		field.Value = ""
		return true
	}

	var value string
	var err error

	if field.IsPassword {
		value, err = form.readPassword()
	} else {
		value, err = form.readString()
	}

	if err != nil {
		// If there's an error (like Ctrl+C), quit
		return false
	}

	// Trim whitespace
	value = strings.TrimSpace(value)

	// If user just pressed Enter without typing anything, keep the current value
	if value == "" {
		// Keep the current value (could be empty if field was cleared)
		return true
	}

	// Validate required fields
	if field.IsRequired && value == "" {
		fmt.Println("This field is required. Press Enter to continue...")
		bufio.NewReader(os.Stdin).ReadString('\n')
		return false
	}

	// Update the value (empty string is valid for clearing fields)
	field.Value = value
	return true
}

// readPassword reads a password input with masking
func (form *InteractiveForm) readPassword() (string, error) {
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	return string(bytePassword), nil
}

// readString reads a regular string input
func (form *InteractiveForm) readString() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(input, "\n"), nil
}
