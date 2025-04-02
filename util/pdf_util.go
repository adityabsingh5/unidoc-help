package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/constant"
	"main/dataobject"
	"strings"

	"github.com/unidoc/unipdf/v3/annotator"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/fjson"
	"github.com/unidoc/unipdf/v3/model"
)

func ReadPdfFields(pdfReader *model.PdfReader) (*dataobject.FormData, error) {
	formData := &dataobject.FormData{
		Fields: []dataobject.FormField{},
	}

	// read acroForm
	acroForm := pdfReader.AcroForm
	if acroForm == nil {
		fmt.Errorf("Inside func FetchPdfFields. No fields in PDF form")
		return formData, nil
	}

	// track unique field names
	uniqueFields := make(map[string]bool)

	// fetch all fields
	fields := acroForm.AllFields()
	for _, field := range fields {
		fieldName, _ := field.FullName()
		lowerFieldName := strings.ToLower(fieldName)

		// Skip processing if the field has already been processed, or
		// if it matches specific prefixes or contains certain substrings
		if uniqueFields[fieldName] ||
			strings.HasPrefix(lowerFieldName, "esign_") ||
			strings.HasPrefix(lowerFieldName, "image_") ||
			strings.Contains(lowerFieldName, "field_rendering_rules") { // TODO: add strings into constants
			continue
		}

		ctx := field.GetContext()

		// Extract the tooltip for FE field label rendering
		tooltip := ""
		if dict, ok := core.GetDict(field.GetContainingPdfObject()); ok {
			if tuObj := dict.Get("TU"); tuObj != nil {
				if tuStr, ok := core.GetString(tuObj); ok {
					tooltip = tuStr.Decoded()
				}
			}
		}

		// check if the field is required
		isRequired := false
		fflags := field.Flags()
		if fflags&(1<<1) != 0 {
			isRequired = true
		}

		switch t := ctx.(type) {
		case *model.PdfFieldText:
			fieldType := constant.TextField
			dateFormat := ""
			// check if field is a date picker
			if strings.Contains(strings.ToLower(fieldName), "_af_date") {
				fieldType = constant.DateField
				// check for json date format
				if aaDict, ok := core.GetDict(field.AA); ok {
					if kObj, ok := aaDict.Get("K").(*core.PdfIndirectObject); ok {
						if dict, ok := core.GetDict(kObj.PdfObject); ok {
							if jsStr, ok := core.GetString(dict.Get("JS")); ok {
								dateFormat = extractDateFormatFromJS(jsStr.Decoded()) // Store the extracted format
							}
						}
					}
				}
			}
			formData.Fields = append(formData.Fields, dataobject.FormField{
				FieldName:  fieldName,
				FieldType:  fieldType,
				FieldLabel: tooltip,
				DateFormat: dateFormat,
				IsRequired: isRequired,
				Options:    nil,
			})
		case *model.PdfFieldButton:
			// add only radio or checkbox field to formData
			var fieldType string
			switch {
			case t.IsCheckbox():
				fieldType = constant.CheckboxField
			case t.IsRadio():
				fieldType = constant.RadioField
			default:
				continue
			}
			// process checkboxes and radio buttons
			var possibleStates []string
			for _, wa := range field.Annotations {
				if apDict, has := core.GetDict(wa.AP); has {
					if dDict, has := core.GetDict(apDict.Get("D")); has {
						for _, key := range dDict.Keys() {
							if key != "Off" { // exclude default "Off" state
								possibleStates = append(possibleStates, key.String())
							}
						}
					}
				}
			}

			formData.Fields = append(formData.Fields, dataobject.FormField{
				FieldName:  fieldName,
				FieldType:  fieldType,
				FieldLabel: tooltip,
				IsRequired: isRequired,
				Options:    possibleStates,
			})

		case *model.PdfFieldChoice:
			// process dropdown fields
			var optionsList []string
			if optArray, ok := core.GetArray(t.Opt); ok {
				for _, item := range optArray.Elements() {
					if optionPair, ok := core.GetArray(item); ok {
						exportValue := ""
						if len(optionPair.Elements()) > 0 {
							if exportObj, ok := core.GetString(optionPair.Get(0)); ok {
								exportValue = exportObj.Decoded()
							}
						}

						optionsList = append(optionsList, exportValue)
					}
				}
			}

			formData.Fields = append(formData.Fields, dataobject.FormField{
				FieldName:  fieldName,
				FieldType:  constant.DropdownField,
				FieldLabel: tooltip,
				IsRequired: isRequired,
				Options:    optionsList,
			})

		default:
			// unsupported field types
			continue
		}
		uniqueFields[fieldName] = true
	}

	return formData, nil
}

func FillPdfFields(pdfReader *model.PdfReader, fieldData map[string]string) error {

	flattenField := make([]string, 0)
	// Convert to slice of KeyValue structs
	var result []dataobject.PdfFieldData
	for k, v := range fieldData {
		result = append(result, dataobject.PdfFieldData{
			FieldName:  k,
			FieldValue: v,
		})
		// Check if key doesn't start with "Esign_" and value is not empty
		if v != "" && !strings.HasPrefix(k, "Esign_") {
			flattenField = append(flattenField, k)
		}
	}

	fmt.Println("Flattened Fields:")
	for _, field := range flattenField {
		fmt.Println(field)
	}

	// Convert to JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		fmt.Errorf("Error marshalling JSON:", "err", err)
		return err
	}

	// Create a reader from the encoded JSON
	jsonReader := bytes.NewReader(jsonData)

	// Load field data from JSON reader
	fdata, err := fjson.LoadFromJSON(jsonReader)
	if err != nil {
		fmt.Errorf("Error loading field data from JSON:", "err", err)
		return err
	}

	// Populate the form data.
	err = pdfReader.AcroForm.Fill(fdata)
	if err != nil {
		fmt.Errorf("Error loading filling AcroForm from fdata:", "err", err)
		return err
	}

	// // Flatten form partially.
	// err = partialFlattenPdf(pdfReader, flattenField)
	// if err != nil {
	// 	logger.AsyncLog(slog.LevelError, "Error in flattening the pdf form:", "err", err)
	// 	return err
	// }

	// Flatten form.
	fieldAppearance := annotator.FieldAppearance{OnlyIfMissing: true, RegenerateTextFields: true}
	// NOTE: To customize certain styles try:
	// style := fieldAppearance.Style()
	// style.CheckmarkGlyph = "a22"
	// style.AutoFontSizeFraction = 0.70
	// fieldAppearance.SetStyle(style)
	//
	// or for specifying a full set of appearance styles:
	// fieldAppearance.SetStyle(annotator.AppearanceStyle{
	//     CheckmarkGlyph:       "a22",
	//     AutoFontSizeFraction: 0.70,
	//     FillColor:            model.NewPdfColorDeviceGray(0.8),
	//     BorderColor:          model.NewPdfColorDeviceRGB(1, 0, 0),
	//     BorderSize:           2.0,
	//     AllowMK:              false,
	// })
	err = pdfReader.FlattenFields(true, fieldAppearance)
	if err != nil {
		fmt.Errorf("Error in flattening the pdf form:", "err", err)
		return err
	}

	return nil
}

// utility function to check if a slice contains a value
func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// Helper function to extract date format
func extractDateFormatFromJS(jsCode string) string {
	start := `AFDate_KeystrokeEx("`
	end := `"`
	if startIdx := strings.Index(jsCode, start); startIdx != -1 {
		startIdx += len(start)
		if endIdx := strings.Index(jsCode[startIdx:], end); endIdx != -1 {
			return jsCode[startIdx : startIdx+endIdx]
		}
	}
	return ""
}
