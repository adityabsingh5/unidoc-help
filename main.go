package main

import (
	"fmt"
	"log"
	"main/util"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/model"
)

func init() {
	err := license.SetMeteredKey("add-unidoc-api-key-here")
	if err != nil {
		panic(err)
	}
}

func main() {
	// Step 1: Read PDF and extract fields
	file, err := os.Open("sample-pdf.pdf")
	if err != nil {
		log.Fatal("Error opening PDF:", err)
	}
	defer file.Close()

	pdfReader, err := model.NewPdfReader(file)
	if err != nil {
		log.Fatal("Error creating PDF reader:", err)
	}

	// Step 2: Extract form fields
	formData, err := util.ReadPdfFields(pdfReader)
	if err != nil {
		log.Fatal("Error reading PDF fields:", err)
	}

	// Step 3: Print form fields for debugging
	fmt.Println("PDF Form Fields:")
	for _, field := range formData.Fields {
		fmt.Printf("\nField: %s\nType: %s\nOptions: %v\n", 
			field.FieldName, field.FieldType, field.Options)
	}

	// Step 4: Fill fields with test values
	fieldValues := map[string]string{
		"TextInput":     "Test Text Input",
		"RadioInput":    "Choice1",
		"CheckboxInput": "Yes",
		"DropdownInput": "DesiredOption", // This should select "DesiredOption" but shows "Default"
	}

	// Step 5: Fill the PDF
	err = util.FillPdfFields(pdfReader, fieldValues)
	if err != nil {
		log.Fatal("Error filling PDF fields:", err)
	}

	// Step 6: Save the filled PDF
	pdfWriter := model.NewPdfWriter()
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		log.Fatal("Error getting number of pages:", err)
	}

	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			log.Fatal("Error getting page:", err)
		}
		pdfWriter.AddPage(page)
	}

	f, err := os.Create("filled_sample-pdf.pdf")
	if err != nil {
		log.Fatal("Error creating output file:", err)
	}
	defer f.Close()

	err = pdfWriter.Write(f)
	if err != nil {
		log.Fatal("Error writing PDF:", err)
	}

	fmt.Println("\nPDF has been filled and saved as filled_sample-pdf.pdf")
}

