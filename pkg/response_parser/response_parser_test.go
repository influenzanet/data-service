package response_parser

import (
	"strconv"
	"testing"

	studyAPI "github.com/influenzanet/study-service/pkg/api"
)

/*
"likertGroup"
"eq5d-health-indicator"
"multipleChoiceGroup"
"singleChoiceGroup"

"responseGroup"
*/

func mockQuestion(
	key string,
	lang string,
	title string,
	responseOptions *studyAPI.ItemComponent,
) *studyAPI.SurveyItem {
	q := studyAPI.SurveyItem{
		Key: key,
		Components: &studyAPI.ItemComponent{
			Role: "root",
			Items: []*studyAPI.ItemComponent{
				{Role: "title", Content: []*studyAPI.LocalisedObject{
					{Code: lang, Parts: []*studyAPI.ExpressionArg{{Data: &studyAPI.ExpressionArg_Str{Str: title}}}},
				}},
				responseOptions,
			},
		},
	}
	return &q
}

type MockOpionDef struct {
	Key   string
	Role  string
	Label string
}

func mockSingleChoiceGroup(lang string, options []MockOpionDef) *studyAPI.ItemComponent {
	rg := studyAPI.ItemComponent{
		Key:  "rg",
		Role: "responseGroup", Items: []*studyAPI.ItemComponent{
			{Key: "scg", Role: "singleChoiceGroup", Items: []*studyAPI.ItemComponent{}},
		}}
	for _, o := range options {
		rg.Items[0].Items = append(rg.Items[0].Items,
			&studyAPI.ItemComponent{Key: o.Key, Role: o.Role, Content: []*studyAPI.LocalisedObject{
				{Code: lang, Parts: []*studyAPI.ExpressionArg{{Data: &studyAPI.ExpressionArg_Str{Str: o.Label}}}},
			}},
		)
	}
	return &rg
}

func mockMultipleChoiceGroup(lang string, options []MockOpionDef) *studyAPI.ItemComponent {
	rg := studyAPI.ItemComponent{
		Key:  "rg",
		Role: "responseGroup", Items: []*studyAPI.ItemComponent{
			{Key: "mcg", Role: "multipleChoiceGroup", Items: []*studyAPI.ItemComponent{}},
		}}
	for _, o := range options {
		rg.Items[0].Items = append(rg.Items[0].Items,
			&studyAPI.ItemComponent{Key: o.Key, Role: o.Role, Content: []*studyAPI.LocalisedObject{
				{Code: lang, Parts: []*studyAPI.ExpressionArg{{Data: &studyAPI.ExpressionArg_Str{Str: o.Label}}}},
			}},
		)
	}
	return &rg
}

func mockLikertGroup(lang string, categoryLabels []MockOpionDef, optionLabels []string) *studyAPI.ItemComponent {
	rg := studyAPI.ItemComponent{
		Key:  "rg",
		Role: "responseGroup", Items: []*studyAPI.ItemComponent{
			{Key: "lg", Role: "likertGroup", Items: []*studyAPI.ItemComponent{}},
		}}

	for i, o := range categoryLabels {
		rg.Items[0].Items = append(rg.Items[0].Items,
			&studyAPI.ItemComponent{Key: strconv.Itoa(i), Role: "text", Content: []*studyAPI.LocalisedObject{
				{Code: lang, Parts: []*studyAPI.ExpressionArg{{Data: &studyAPI.ExpressionArg_Str{Str: o.Label}}}},
			}},
		)
		rg.Items[0].Items = append(rg.Items[0].Items,
			&studyAPI.ItemComponent{Key: o.Key, Role: "likert", Items: []*studyAPI.ItemComponent{}},
		)

		index := len(rg.Items[0].Items) - 1
		for j, label := range optionLabels {
			rg.Items[0].Items[index].Items = append(rg.Items[0].Items[index].Items, &studyAPI.ItemComponent{Key: strconv.Itoa(j + 1), Role: "option", Content: []*studyAPI.LocalisedObject{
				{Code: lang, Parts: []*studyAPI.ExpressionArg{{Data: &studyAPI.ExpressionArg_Str{Str: label}}}},
			}})
		}
	}
	return &rg
}

func TestResponseParser(t *testing.T) {
	testLang := "en"
	testSurveyDef := &studyAPI.SurveyItem{
		Key: "weekly",
		Items: []*studyAPI.SurveyItem{
			mockQuestion("weekly.Q1", testLang, "Title of Q1", mockSingleChoiceGroup(testLang, []MockOpionDef{
				{Key: "1", Role: "option", Label: "Yes"},
				{Key: "2", Role: "option", Label: "No"},
				{Key: "3", Role: "input", Label: "Other"},
			})),
			mockQuestion("weekly.Q2", testLang, "Title of Q2", mockMultipleChoiceGroup(testLang, []MockOpionDef{
				{Key: "1", Role: "option", Label: "Option 1"},
				{Key: "2", Role: "option", Label: "Option 2"},
				{Key: "3", Role: "input", Label: "Other"},
			})),
			{Key: "weeky.G1", Items: []*studyAPI.SurveyItem{
				mockQuestion("weekly.G1.Q1", testLang, "Title of Group 1's Q1", mockLikertGroup(testLang, []MockOpionDef{
					{Key: "cat1", Label: "Category 1"},
					{Key: "cat2", Label: "Category 2"},
				}, []string{
					"o1", "o2", "o3",
				})),
			}},
		},
	}

	t.Run("with with missing surveyDef", func(t *testing.T) {
		_, err := NewResponseParser(nil, "en", true, true)
		if err == nil {
			t.Error("error expected")
			return
		}
		if err.Error() != "current survey definition not found" {
			t.Errorf("unexpected error message: %v", err)
			return
		}
	})

	t.Run("with with missing current", func(t *testing.T) {
		testSurvey := studyAPI.Survey{
			Id:      "test-id",
			Current: nil,
			History: []*studyAPI.SurveyVersion{},
		}

		_, err := NewResponseParser(&testSurvey, "en", true, true)
		if err == nil {
			t.Error("error expected")
			return
		}
		if err.Error() != "current survey definition not found" {
			t.Errorf("unexpected error message: %v", err)
			return
		}
	})

	t.Run("with with one version", func(t *testing.T) {
		testSurvey := studyAPI.Survey{
			Id: "test-id",
			Current: &studyAPI.SurveyVersion{
				Published:        10,
				VersionId:        "1",
				SurveyDefinition: testSurveyDef,
			},
		}

		rp, err := NewResponseParser(&testSurvey, "en", true, true)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if len(rp.surveyVersions) != 1 {
			t.Errorf("unexpected number of versions: %d", len(rp.surveyVersions))
			return
		}

		if len(rp.surveyVersions[0].Questions) != 3 {
			t.Errorf("unexpected number of versions: %d", len(rp.surveyVersions[0].Questions))
			return
		}
	})

	t.Run("with with multiple versions", func(t *testing.T) {
		testSurvey := studyAPI.Survey{
			Id: "test-id",
			Current: &studyAPI.SurveyVersion{
				Published:        10,
				VersionId:        "3",
				SurveyDefinition: testSurveyDef,
			},
			History: []*studyAPI.SurveyVersion{
				{
					Published:        2,
					Unpublished:      5,
					VersionId:        "1",
					SurveyDefinition: testSurveyDef,
				},
				{
					Published:        5,
					Unpublished:      10,
					VersionId:        "2",
					SurveyDefinition: testSurveyDef,
				},
			},
		}

		rp, err := NewResponseParser(&testSurvey, "en", true, true)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if len(rp.surveyVersions) != 3 {
			t.Errorf("unexpected number of versions: %d", len(rp.surveyVersions))
			return
		}
	})
}
