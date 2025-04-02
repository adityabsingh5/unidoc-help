# PDF Form Field Handler

A Go utility to read and fill PDF form fields using the UniDoc library. This project demonstrates handling different types of PDF form fields including text, radio, checkbox, and dropdown fields.

## Project Structure
```
sample-pdf-util/
├── main.go              # Main application code
├── util/
│   └── pdf_util.go      # PDF utility functions
├── dataobject/
│   └── pdf_fields.go    # Data structures for form fields
└── constant/
    └── constant.go      # Constants for field types
```

## Prerequisites
- Go 1.16 or higher
- UniDoc license key (trial or commercial)

## Setup
1. Clone the repository:
```bash
git clone https://github.com/adityabsingh5/unidoc-help.git
cd unidoc-help
```

2. Install dependencies:
```bash
go mod tidy
```

3. Place your sample PDF file as `sample-pdf.pdf` in the root directory

## Usage
Run the application:
```bash
go run main.go
```

The program will:
1. Read the PDF form fields
2. Display all fields and their options
3. Fill the fields with test values:
   - TextInput: "Test Text Input"
   - RadioInput: "Choice1"
   - CheckboxInput: "Yes"
   - DropdownInput: "DesiredOption"
4. Save the filled PDF as `filled_sample-pdf.pdf`

## Known Issues
- Dropdown fields may not properly update from "Default" to selected value
- This is being investigated with UniDoc support

## License
This project is licensed under the MIT License - see the LICENSE file for details. 