package response_parser

import (
	"errors"
	"log"
	"strconv"
	"strings"

	studyAPI "github.com/influenzanet/study-service/pkg/api"
)

func findSurveyVersion(versionID string, submittedAt int64, versions []SurveyVersionPreview) (sv SurveyVersionPreview, err error) {
	if versionID == "" {
		return findVersionBasedOnTimestamp(submittedAt, versions)
	} else {
		sv, err = findVersionBasedOnVersionID(versionID, versions)
		if err != nil {
			return findVersionBasedOnTimestamp(submittedAt, versions)
		}
	}
	return sv, nil
}

func findVersionBasedOnTimestamp(submittedAt int64, versions []SurveyVersionPreview) (sv SurveyVersionPreview, err error) {
	for _, v := range versions {
		if v.Unpublished == 0 {
			if v.Published <= submittedAt {
				return v, nil
			}
		} else {
			if v.Published <= submittedAt && v.Unpublished > submittedAt {
				return v, nil
			}
		}
	}
	return sv, errors.New("no survey version found")
}

func findVersionBasedOnVersionID(versionID string, versions []SurveyVersionPreview) (sv SurveyVersionPreview, err error) {
	for _, v := range versions {
		if v.VersionID == versionID {
			return v, nil
		}
	}
	return sv, errors.New("no survey version found")
}

func timestampsToStr(ts []int64, sep string) string {
	if len(ts) == 0 {
		return ""
	}

	b := make([]string, len(ts))
	for i, v := range ts {
		b[i] = strconv.Itoa(int(v))
	}
	return strings.Join(b, sep)
}

func findResponse(responses []*studyAPI.SurveyItemResponse, key string) *studyAPI.SurveyItemResponse {
	for _, r := range responses {
		if r.Key == key {
			return r
		}
	}
	return nil
}

func getResponseColumns(question SurveyQuestion, response *studyAPI.SurveyItemResponse) map[string]string {
	switch question.QuestionType {
	case QUESTION_TYPE_SINGLE_CHOICE:
		return generateResponseForSingleChoice(question, response)
	default:
		return map[string]string{}
	}
}

func generateResponseForSingleChoice(question SurveyQuestion, response *studyAPI.SurveyItemResponse) map[string]string {
	responseCols := map[string]string{}

	// Prepare response columns
	if len(question.Responses) == 1 {
		responseCols[question.ID] = ""
		rSlot := question.Responses[0]

		for _, option := range rSlot.Options {
			if option.OptionType != OPTION_TYPE_RADIO {
				responseCols[question.ID+"-"+option.ID] = ""
			}
		}
	} else {
		for _, rSlot := range question.Responses {
			responseCols[question.ID+"-"+rSlot.ID] = ""
			for _, option := range rSlot.Options {
				if option.OptionType != OPTION_TYPE_RADIO {
					responseCols[question.ID+"-"+rSlot.ID+"."+option.ID] = ""
				}
			}
		}
	}

	// Find responses
	if len(question.Responses) == 1 {
		responseCols[question.ID] = ""
		rSlot := question.Responses[0]
		rGroup := retrieveResponseItem(response, "rg."+rSlot.ID)
		if rGroup != nil {
			if len(rGroup.Items) != 1 {
				log.Printf("unexpected response group for question %s: %v", question.ID, rGroup)
			} else {
				selection := rGroup.Items[0]
				responseCols[question.ID] = selection.Key
				responseCols[question.ID+"-"+selection.Key] = selection.Value
			}
		}
	} else {
		for _, rSlot := range question.Responses {
			rGroup := retrieveResponseItem(response, "rg."+rSlot.ID)
			if rGroup == nil {
				continue
			} else if len(rGroup.Items) != 1 {
				log.Printf("unexpected response group for question %s: %v", question.ID, rGroup)
				continue
			}

			selection := rGroup.Items[0]
			responseCols[question.ID+"-"+rSlot.ID] = selection.Key
			responseCols[question.ID+"-"+rSlot.ID+"."+selection.Key] = selection.Value
		}
	}

	log.Println(responseCols)
	return responseCols
}

func retrieveResponseItem(response *studyAPI.SurveyItemResponse, fullKey string) *studyAPI.ResponseItem {
	if response == nil {
		return nil
	}
	keyParts := strings.Split(fullKey, ".")

	var result *studyAPI.ResponseItem
	for _, key := range keyParts {
		if result == nil {
			if key != response.Response.Key {
				return nil
			}
			result = response.Response
			continue
		}
		found := false
		for _, item := range result.Items {
			if item.Key == key {
				result = item
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}
	return result
}

/*
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
*/
