package response_parser

import (
	"errors"
	"log"

	studyAPI "github.com/influenzanet/study-service/pkg/api"
)

func surveyDefToVersionPreview(original *studyAPI.SurveyVersion, prefLang string) SurveyVersionPreview {
	sp := SurveyVersionPreview{
		VersionID:   original.VersionId,
		Published:   original.Published,
		Unpublished: original.Unpublished,
		Questions:   []SurveyQuestion{},
	}

	sp.Questions = extractQuestions(original.SurveyDefinition, prefLang)
	return sp
}

func extractQuestions(root *studyAPI.SurveyItem, prefLang string) []SurveyQuestion {
	questions := []SurveyQuestion{}
	if root == nil {
		return questions
	}
	for _, item := range root.Items {
		if item.Type == "pageBreak" {
			continue
		}

		if isItemGroup(item) {
			questions = append(questions, extractQuestions(item, prefLang)...)
			continue
		}

		rg := getResponseGroupComponent(item)
		if rg == nil {
			continue
		}

		qType := getQuestionType(rg)

		// TODO: get response options
		responseOptions := []ResponseOption{}

		titleComp := getTitleComponent(item)
		title := ""
		if titleComp != nil {
			var err error
			title, err = getTranslation(titleComp.Content, prefLang)
			if err != nil {
				log.Printf("Question %s title error: %v", item.Key, err)
			}
		}

		question := SurveyQuestion{
			ID:              item.Key,
			Title:           title,
			QuestionType:    qType,
			ResponseOptions: responseOptions,
		}
		questions = append(questions, question)
	}
	return questions
}

func isItemGroup(item *studyAPI.SurveyItem) bool {
	return item != nil && len(item.Items) > 0
}

func getResponseGroupComponent(question *studyAPI.SurveyItem) *studyAPI.ItemComponent {
	if question.Components == nil {
		return nil
	}
	for _, c := range question.Components.Items {
		if c.Role == "responseGroup" {
			return c
		}
	}
	return nil
}

func getTitleComponent(question *studyAPI.SurveyItem) *studyAPI.ItemComponent {
	if question.Components == nil {
		return nil
	}
	for _, c := range question.Components.Items {
		if c.Role == "title" {
			return c
		}
	}
	return nil
}

func getQuestionType(rg *studyAPI.ItemComponent) string {
	if rg == nil {
		return QUESTION_TYPE_UNKNOWN
	}

	if len(rg.Items) == 1 {
		role := rg.Items[0].Role
		if role == "singleChoiceGroup" {
			return QUESTION_TYPE_SINGLE_CHOICE
		} else if role == "multipleChoiceGroup" {
			return QUESTION_TYPE_MULTIPLE_CHOICE
		}

	} else if len(rg.Items) > 1 {

	}
	return QUESTION_TYPE_UNKNOWN
}

/*
"singleChoiceGroup"
"multipleChoiceGroup"
"dropDownGroup"
"input"
"numberInput"
"dateInput"
"multilineTextInput"
"eq5d-health-indicator"
"sliderNumeric"
"matrix"
"likert"
*/

func getTranslation(content []*studyAPI.LocalisedObject, lang string) (string, error) {
	if len(content) < 1 {
		return "", errors.New("translations missing")
	}

	for _, translation := range content {
		if translation.Code == lang {
			mergedText := ""
			for _, p := range translation.Parts {
				if p.Dtype == "str" {
					mergedText += p.GetStr()
				}
			}
			return mergedText, nil
		}
	}
	return "", errors.New("translation missing")
}
