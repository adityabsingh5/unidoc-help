package dataobject

// key-value pair for filling PDF fields
type PdfFieldData struct {
	FieldName  string `json:"name"`
	FieldValue string `json:"value"`
}

// form fields metadata
type FormData struct {
	Fields []FormField `json:"formFields"`
}

type FormField struct {
	FieldName  string   `json:"name"`
	FieldType  string   `json:"type"`
	FieldLabel string   `json:"label"`
	DateFormat string   `json:"dateFormat,omitempty"`
	IsRequired bool     `json:"isRequired,omitempty"`
	Options    []string `json:"options,omitempty"`
}

// esign fields metadata
type ESignFormData struct {
	ESignFields []ESignField `json:"eSignFields"`
}

type ESignField struct {
	Identifier         string              `json:"identifier"`
	SignaturePositions []SignaturePosition `json:"signaturePositions"`
	DatePositions      []DatePosition      `json:"datePositions"`
}

type SignaturePosition struct {
	Page        int          `json:"page"`
	Coordinates []Coordinate `json:"coordinates"`
}

type DatePosition struct {
	Page        int          `json:"page"`
	Coordinates []Coordinate `json:"coordinates"`
}

type Coordinate struct {
	LLX int `json:"llx"`
	LLY int `json:"lly"`
	URX int `json:"urx"`
	URY int `json:"ury"`
}

// AddSignaturePosition adds a signature position for the given identifier
func (f *ESignFormData) AddSignaturePosition(identifier string, page int, coordinate Coordinate) {
	for i, field := range f.ESignFields {
		if field.Identifier == identifier {
			f.ESignFields[i].SignaturePositions = append(f.ESignFields[i].SignaturePositions, SignaturePosition{
				Page:        page,
				Coordinates: []Coordinate{coordinate},
			})
			return
		}
	}
	f.ESignFields = append(f.ESignFields, ESignField{
		Identifier: identifier,
		SignaturePositions: []SignaturePosition{
			{
				Page:        page,
				Coordinates: []Coordinate{coordinate},
			},
		},
	})
}

// AddDatePosition adds a date position for the given identifier
func (f *ESignFormData) AddDatePosition(identifier string, page int, coordinate Coordinate) {
	for i, field := range f.ESignFields {
		if field.Identifier == identifier {
			f.ESignFields[i].DatePositions = append(f.ESignFields[i].DatePositions, DatePosition{
				Page:        page,
				Coordinates: []Coordinate{coordinate},
			})
			return
		}
	}
	f.ESignFields = append(f.ESignFields, ESignField{
		Identifier: identifier,
		DatePositions: []DatePosition{
			{
				Page:        page,
				Coordinates: []Coordinate{coordinate},
			},
		},
	})
}

// ImageFormData holds metadata for image fields in a PDF
type ImageFormData struct {
	ImageFields []ImageField `json:"imageFields"`
}

type ImageField struct {
	Identifier string      `json:"identifier"`
	ImageData  []ImageData `json:"imageData"`
}

type ImageData struct {
	Page        int        `json:"page"`
	Coordinates Coordinate `json:"coordinates"`
}

func (f *ImageFormData) AddImageField(identifier string, page int, coordinate Coordinate) {
	for i, field := range f.ImageFields {
		if field.Identifier == identifier {
			f.ImageFields[i].ImageData = append(f.ImageFields[i].ImageData, ImageData{
				Page:        page,
				Coordinates: coordinate,
			})
			return
		}
	}
	f.ImageFields = append(f.ImageFields, ImageField{
		Identifier: identifier,
		ImageData: []ImageData{
			{
				Page:        page,
				Coordinates: coordinate,
			},
		},
	})
}

// FieldRenderingRules holds rules for rendering fields based on conditions
type FieldRenderingRules struct {
	Conditions []Condition `json:"conditions"`
}

type Condition struct {
	Field string `json:"field"`
	Type  string `json:"type"`
	Rules []Rule `json:"rules"`
}

type Rule struct {
	Value  string   `json:"value"`
	Render []string `json:"render"`
}

// CombinedFormMetadata holds all the metadata for the form
type CombinedFormMetadata struct {
	FormFields          FormData            `json:"formFieldsMetaData"`
	FieldRenderingRules FieldRenderingRules `json:"fieldRenderingRulesMetaData"`
}
