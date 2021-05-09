package response_parser

import studyAPI "github.com/influenzanet/study-service/pkg/api"

func surveyDefToVersionPreview(original *studyAPI.SurveyVersion) SurveyVersionPreview {
	sp := SurveyVersionPreview{
		VersionID:   original.VersionId,
		Published:   original.Published,
		Unpublished: original.Unpublished,
		Questions:   []SurveyQuestion{},
	}

	sp.Questions = extractQuestions(original.SurveyDefinition)
	return sp
}

func extractQuestions(root *studyAPI.SurveyItem) []SurveyQuestion {
	questions := []SurveyQuestion{}
	if root == nil {
		return questions
	}
	for _, item := range root.Items {
		if item.Type == "pageBreak" {
			continue
		}

		if isItemGroup(item) {
			questions = append(questions, extractQuestions(item)...)
			continue
		}

		rg := getResponseGroupComponent(item)
		// TODO: find response group -> if not, continue
		// TODO: get question type (based on parsed response group)
	}
	return questions
}

func isItemGroup(item *studyAPI.SurveyItem) bool {
	return len(item.Items) > 0
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
