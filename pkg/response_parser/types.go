package response_parser

const (
	QUESTION_TYPE_SINGLE_CHOICE   = "single_choice"
	QUESTION_TYPE_MULTIPLE_CHOICE = "multiple_choice"
	QUESTION_TYPE_INPUT           = "input"
	QUESTION_TYPE_NUMBER_INPUT    = "numberInput"
	QUESTION_TYPE_DATE            = "dateInput"
	QUESTION_TYPE_DROPDOWN_GROUP  = "dropdown_group"
	QUESTION_TYPE_LIKERT_GROUP    = "likert_group"
	QUESTION_TYPE_EQ5D_SLIDER     = "eq5d_slider"
	QUESTION_TYPE_NUMERIC_SLIDER  = "numeric_slider"
	QUESTION_TYPE_MATRIX          = "matrix"
	QUESTION_TYPE_UNKNOWN         = "unknown"
)

const (
	OPTION_TYPE_RADIO          = "radio"
	OPTION_TYPE_CHECKBOX       = "checkbox"
	OPTION_TYPE_TEXT_INPUT     = "text_input"
	OPTION_TYPE_DATE_INPUT     = "date"
	OPTION_TYPE_NUMBER_INPUT   = "number_input"
	OPTION_TYPE_NUMERIC_SLIDER = "numeric_slider"
)

type SurveyVersionPreview struct {
	VersionID   string
	Published   int64
	Unpublished int64
	Questions   []SurveyQuestion
}

type SurveyQuestion struct {
	ID              string
	Title           string
	QuestionType    string
	ResponseOptions []ResponseOption
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
	Language      string
	Responses     map[string]string
}