package response_parser

const (
	QUESTION_TYPE_SINGLE_CHOICE       = "single_choice"
	QUESTION_TYPE_MULTIPLE_CHOICE     = "multiple_choice"
	QUESTION_TYPE_TEXT_INPUT          = "text"
	QUESTION_TYPE_NUMBER_INPUT        = "number"
	QUESTION_TYPE_DATE_INPUT          = "date"
	QUESTION_TYPE_DROPDOWN            = "dropdown"
	QUESTION_TYPE_LIKERT              = "likert"
	QUESTION_TYPE_EQ5D_SLIDER         = "eq5d_slider"
	QUESTION_TYPE_NUMERIC_SLIDER      = "slider"
	QUESTION_TYPE_MATRIX              = "matrix"
	QUESTION_TYPE_MATRIX_RADIO_ROW    = "matrix_radio_row"
	QUESTION_TYPE_MATRIX_DROPDOWN     = "matrix_dropdown"
	QUESTION_TYPE_MATRIX_INPUT        = "matrix_input"
	QUESTION_TYPE_MATRIX_NUMBER_INPUT = "matrix_number_input"
	QUESTION_TYPE_MATRIX_CHECKBOX     = "matrix_checkbox"
	QUESTION_TYPE_UNKNOWN             = "unknown"
	QUESTION_TYPE_EMPTY               = "empty"
)

const (
	OPTION_TYPE_DROPDOWN_OPTION = "option"
	OPTION_TYPE_RADIO           = "radio"
	OPTION_TYPE_CHECKBOX        = "checkbox"
	OPTION_TYPE_TEXT_INPUT      = "text"
	OPTION_TYPE_DATE_INPUT      = "date"
	OPTION_TYPE_NUMBER_INPUT    = "number"
)

const (
	RESPONSE_ROOT_KEY = "rg"
)

const (
	OPEN_FIELD_COL_SUFFIX = "open"
	TRUE_VALUE            = "TRUE"
	FALSE_VALUE           = "FALSE"
)

type SurveyVersionPreview struct {
	VersionID   string
	Published   int64
	Unpublished int64
	Questions   []SurveyQuestion
}

type SurveyQuestion struct {
	ID           string
	Title        string
	QuestionType string
	Responses    []ResponseDef
}

type ResponseDef struct {
	ID           string
	ResponseType string
	Label        string
	Options      []ResponseOption
}

type ResponseOption struct {
	ID         string
	OptionType string
	Label      string
}

type ParsedResponse struct {
	ParticipantID string
	SubmittedAt   int64
	Version       string
	Context       map[string]string // e.g. Language, or engine version
	Responses     map[string]string
	Meta          ResponseMeta
}

type ResponseMeta struct {
	Initialised map[string]string
	Displayed   map[string]string
	Responded   map[string]string
	ItemVersion map[string]string
}
