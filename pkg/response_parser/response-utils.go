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

func getResponseColumns(question SurveyQuestion, response *studyAPI.SurveyItemResponse, questionOptionSep string) map[string]string {
	switch question.QuestionType {
	case QUESTION_TYPE_SINGLE_CHOICE:
		return processResponseForSingleChoice(question, response, questionOptionSep)
	case QUESTION_TYPE_DROPDOWN:
		return processResponseForSingleChoice(question, response, questionOptionSep)
	case QUESTION_TYPE_LIKERT:
		return processResponseForSingleChoice(question, response, questionOptionSep)
	case QUESTION_TYPE_MULTIPLE_CHOICE:
		return processResponseForMultipleChoice(question, response, questionOptionSep)
	case QUESTION_TYPE_TEXT_INPUT:
		return processResponseForInputs(question, response, questionOptionSep)
		// TODO
		/*
			QUESTION_TYPE_TEXT_INPUT          = "text"
			QUESTION_TYPE_NUMBER_INPUT        = "number"
			QUESTION_TYPE_DATE_INPUT          = "date"
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
	default:
		return map[string]string{}
	}
}

func processResponseForSingleChoice(question SurveyQuestion, response *studyAPI.SurveyItemResponse, questionOptionSep string) map[string]string {
	var responseCols map[string]string

	if len(question.Responses) == 1 {
		rSlot := question.Responses[0]
		responseCols = handleSimpleSingleChoiceGroup(question.ID, rSlot, response, questionOptionSep)

	} else {
		responseCols = handleSingleChoiceGroupList(question.ID, question.Responses, response, questionOptionSep)
	}
	return responseCols
}

func handleSimpleSingleChoiceGroup(questionKey string, responseSlotDef ResponseDef, response *studyAPI.SurveyItemResponse, questionOptionSep string) map[string]string {
	responseCols := map[string]string{}

	// Prepare columns:
	responseCols[questionKey] = ""

	for _, option := range responseSlotDef.Options {
		if option.OptionType != OPTION_TYPE_RADIO &&
			option.OptionType != OPTION_TYPE_DROPDOWN_OPTION {
			responseCols[questionKey+questionOptionSep+option.ID] = ""
		}
	}

	// Find responses
	rGroup := retrieveResponseItem(response, RESPONSE_ROOT_KEY+"."+responseSlotDef.ID)
	if rGroup != nil {
		if len(rGroup.Items) != 1 {
			log.Printf("unexpected response group for question %s: %v", questionKey, rGroup)
		} else {
			selection := rGroup.Items[0]
			responseCols[questionKey] = selection.Key

			valueKey := questionKey + questionOptionSep + selection.Key
			if _, hasKey := responseCols[valueKey]; hasKey {
				responseCols[valueKey] = selection.Value
			}
		}
	}
	return responseCols
}

func handleSingleChoiceGroupList(questionKey string, responseSlotDefs []ResponseDef, response *studyAPI.SurveyItemResponse, questionOptionSep string) map[string]string {
	responseCols := map[string]string{}

	// Prepare columns:
	for _, rSlot := range responseSlotDefs {
		responseCols[questionKey+questionOptionSep+rSlot.ID] = ""
		for _, option := range rSlot.Options {
			if option.OptionType != OPTION_TYPE_RADIO &&
				option.OptionType != OPTION_TYPE_DROPDOWN_OPTION {
				responseCols[questionKey+questionOptionSep+rSlot.ID+"."+option.ID] = ""
			}
		}
	}

	// Find responses:
	for _, rSlot := range responseSlotDefs {
		rGroup := retrieveResponseItem(response, RESPONSE_ROOT_KEY+"."+rSlot.ID)
		if rGroup == nil {
			continue
		} else if len(rGroup.Items) != 1 {
			log.Printf("unexpected response group for question %s: %v", questionKey, rGroup)
			continue
		}

		selection := rGroup.Items[0]
		responseCols[questionKey+questionOptionSep+rSlot.ID] = selection.Key

		valueKey := questionKey + questionOptionSep + rSlot.ID + "." + selection.Key
		if _, hasKey := responseCols[valueKey]; hasKey {
			responseCols[valueKey] = selection.Value
		}
	}
	return responseCols
}

func processResponseForMultipleChoice(question SurveyQuestion, response *studyAPI.SurveyItemResponse, questionOptionSep string) map[string]string {
	var responseCols map[string]string

	if len(question.Responses) == 1 {
		rSlot := question.Responses[0]
		responseCols = handleSimpleMultipleChoiceGroup(question.ID, rSlot, response, questionOptionSep)

	} else {
		responseCols = handleMultipleChoiceGroupList(question.ID, question.Responses, response, questionOptionSep)
	}
	return responseCols
}

func handleSimpleMultipleChoiceGroup(questionKey string, responseSlotDef ResponseDef, response *studyAPI.SurveyItemResponse, questionOptionSep string) map[string]string {
	responseCols := map[string]string{}

	// Find responses
	rGroup := retrieveResponseItem(response, RESPONSE_ROOT_KEY+"."+responseSlotDef.ID)
	if rGroup != nil {
		if len(rGroup.Items) > 0 {
			for _, option := range responseSlotDef.Options {
				responseCols[questionKey+questionOptionSep+option.ID] = "FALSE"
			}

			for _, item := range rGroup.Items {
				value := "TRUE"
				if item.Value != "" {
					value = item.Value
				}
				responseCols[questionKey+questionOptionSep+item.Key] = value
			}
		}
	} else {
		for _, option := range responseSlotDef.Options {
			responseCols[questionKey+questionOptionSep+option.ID] = ""
		}

	}
	return responseCols
}

func handleMultipleChoiceGroupList(questionKey string, responseSlotDefs []ResponseDef, response *studyAPI.SurveyItemResponse, questionOptionSep string) map[string]string {
	responseCols := map[string]string{}

	// Prepare columns:
	for _, rSlot := range responseSlotDefs {
		// Find responses
		rGroup := retrieveResponseItem(response, RESPONSE_ROOT_KEY+"."+rSlot.ID)
		slotKeyPrefix := questionKey + questionOptionSep + rSlot.ID + "."
		if rGroup != nil {
			if len(rGroup.Items) > 0 {
				for _, option := range rSlot.Options {
					responseCols[slotKeyPrefix+option.ID] = "FALSE"
				}

				for _, item := range rGroup.Items {
					value := "TRUE"
					if item.Value != "" {
						value = item.Value
					}
					responseCols[slotKeyPrefix+item.Key] = value
				}
			}
		} else {
			for _, option := range rSlot.Options {
				responseCols[slotKeyPrefix+option.ID] = ""
			}

		}
	}

	return responseCols
}

func processResponseForInputs(question SurveyQuestion, response *studyAPI.SurveyItemResponse, questionOptionSep string) map[string]string {
	var responseCols map[string]string

	if len(question.Responses) == 1 {
		rSlot := question.Responses[0]
		responseCols = handleSimpleInput(question.ID, rSlot, response, questionOptionSep)

	} else {
		responseCols = handleInputList(question.ID, question.Responses, response, questionOptionSep)
	}
	log.Println(responseCols)
	return responseCols
}

func handleSimpleInput(questionKey string, responseSlotDef ResponseDef, response *studyAPI.SurveyItemResponse, questionOptionSep string) map[string]string {
	responseCols := map[string]string{}
	responseCols[questionKey] = ""

	// Find responses
	rValue := retrieveResponseItem(response, RESPONSE_ROOT_KEY+"."+responseSlotDef.ID)
	if rValue != nil {
		responseCols[questionKey] = rValue.Value
	}
	return responseCols
}

func handleInputList(questionKey string, responseSlotDefs []ResponseDef, response *studyAPI.SurveyItemResponse, questionOptionSep string) map[string]string {
	responseCols := map[string]string{}

	for _, rSlot := range responseSlotDefs {
		// Prepare columns:
		slotKey := questionKey + questionOptionSep + rSlot.ID
		responseCols[slotKey] = ""

		// Find responses
		rValue := retrieveResponseItem(response, RESPONSE_ROOT_KEY+"."+rSlot.ID)
		if rValue != nil {
			responseCols[slotKey] = rValue.Value
		}
	}

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
