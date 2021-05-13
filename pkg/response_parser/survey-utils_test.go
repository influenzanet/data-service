package response_parser

import (
	"testing"

	studyAPI "github.com/influenzanet/study-service/pkg/api"
)

func TestIsItemGroup(t *testing.T) {
	testLang := "en"

	testItem1 := &studyAPI.SurveyItem{Key: "weeky.G1", Items: []*studyAPI.SurveyItem{
		mockQuestion("weekly.G1.Q1", testLang, "Title of Group 1's Q1", mockLikertGroup(testLang, []MockOpionDef{
			{Key: "cat1", Label: "Category 1"},
			{Key: "cat2", Label: "Category 2"},
		}, []string{
			"o1", "o2", "o3",
		})),
	}}
	testItem2 := mockQuestion("weekly.Q2", testLang, "Title of Q2", mockMultipleChoiceGroup(testLang, []MockOpionDef{
		{Key: "1", Role: "option", Label: "Option 1"},
		{Key: "2", Role: "option", Label: "Option 2"},
		{Key: "3", Role: "input", Label: "Other"},
	}))

	t.Run("with with missing item", func(t *testing.T) {
		if isItemGroup(nil) {
			t.Error("missing item wrongly as group")
		}
	})

	t.Run("with with single item", func(t *testing.T) {
		if isItemGroup(testItem2) {
			t.Error("single item wrongly as group")
		}
	})

	t.Run("with with group item", func(t *testing.T) {
		if !isItemGroup(testItem1) {
			t.Error("group item not recognized")
		}
	})
}

func TestGetResponseGroupComponent(t *testing.T) {
	testLang := "en"

	testItem1 := mockQuestion("weekly.Q2", testLang, "Title of Q2", mockMultipleChoiceGroup(testLang, []MockOpionDef{
		{Key: "1", Role: "option", Label: "Option 1"},
		{Key: "2", Role: "option", Label: "Option 2"},
		{Key: "3", Role: "input", Label: "Other"},
	}))

	t.Run("with test items", func(t *testing.T) {
		rg := getResponseGroupComponent(testItem1)
		if rg == nil {
			t.Error("rg empty")
			return
		}
		if rg.Role != "responseGroup" {
			t.Errorf("unexpected role: %s", rg.Role)
			return
		}
	})
}

func TestGetQuestionType(t *testing.T) {

	t.Run("missing response group component", func(t *testing.T) {
		if qt := getQuestionType(nil); qt != QUESTION_TYPE_UNKNOWN {
			t.Errorf("unexpected question type: %s", qt)
		}
	})

	t.Run("missing items", func(t *testing.T) {
		rg := &studyAPI.ItemComponent{
			Key:   "rg",
			Role:  "responseGroup",
			Items: []*studyAPI.ItemComponent{},
		}
		if qt := getQuestionType(rg); qt != QUESTION_TYPE_UNKNOWN {
			t.Errorf("unexpected question type: %s", qt)
		}
	})

	t.Run("multiple items (unknown)", func(t *testing.T) {
		rg := &studyAPI.ItemComponent{
			Key:  "rg",
			Role: "responseGroup",
			Items: []*studyAPI.ItemComponent{
				{Key: "1", Role: "Text"},
				{Key: "2", Role: "Something"},
				{Key: "3", Role: "More"},
			},
		}
		if qt := getQuestionType(rg); qt != QUESTION_TYPE_UNKNOWN {
			t.Errorf("unexpected question type: %s", qt)
		}
	})

	t.Run("singleChoiceGroup", func(t *testing.T) {
		rg := &studyAPI.ItemComponent{
			Key:  "rg",
			Role: "responseGroup",
			Items: []*studyAPI.ItemComponent{
				{Key: "scg", Role: "singleChoiceGroup", Items: []*studyAPI.ItemComponent{
					{Key: "1", Role: "option"},
					{Key: "2", Role: "option"},
				}},
			},
		}
		if qt := getQuestionType(rg); qt != QUESTION_TYPE_SINGLE_CHOICE {
			t.Errorf("unexpected question type: %s", qt)
		}
	})

	t.Run("multipleChoiceGroup", func(t *testing.T) {
		rg := &studyAPI.ItemComponent{
			Key:  "rg",
			Role: "responseGroup",
			Items: []*studyAPI.ItemComponent{
				{Key: "scg", Role: "multipleChoiceGroup", Items: []*studyAPI.ItemComponent{
					{Key: "1", Role: "option"},
					{Key: "2", Role: "option"},
				}},
			},
		}
		if qt := getQuestionType(rg); qt != QUESTION_TYPE_MULTIPLE_CHOICE {
			t.Errorf("unexpected question type: %s", qt)
		}
	})

	t.Run("dropDownGroup", func(t *testing.T) {
		t.Error("test unimplemented")
	})

	t.Run("input", func(t *testing.T) {
		t.Error("test unimplemented")
	})

	t.Run("multilineTextInput", func(t *testing.T) {
		t.Error("test unimplemented")
	})

	t.Run("numberInput", func(t *testing.T) {
		t.Error("test unimplemented")
	})

	t.Run("dateInput", func(t *testing.T) {
		t.Error("test unimplemented")
	})

	t.Run("eq5d-health-indicator", func(t *testing.T) {
		t.Error("test unimplemented")
	})

	t.Run("sliderNumeric", func(t *testing.T) {
		t.Error("test unimplemented")
	})

	t.Run("matrix", func(t *testing.T) {
		t.Error("test unimplemented")
	})

	t.Run("likerts - but not likertGroup", func(t *testing.T) {
		t.Error("test unimplemented")
	})

	t.Run("likertGroup", func(t *testing.T) {
		t.Error("test unimplemented")
	})

}

func TestGetTranslation(t *testing.T) {

	t.Run("with empty translation list", func(t *testing.T) {
		_, err := getTranslation([]*studyAPI.LocalisedObject{}, "en")
		if err == nil {
			t.Error("should return an error")
			return
		}
		if err.Error() != "translations missing" {
			t.Errorf("unexpected error: %v", err)
			return
		}
	})

	t.Run("with missing translation", func(t *testing.T) {
		_, err := getTranslation([]*studyAPI.LocalisedObject{
			{Code: "de", Parts: []*studyAPI.ExpressionArg{{Dtype: "str", Data: &studyAPI.ExpressionArg_Str{Str: "Test DE"}}}},
			{Code: "nl", Parts: []*studyAPI.ExpressionArg{{Dtype: "str", Data: &studyAPI.ExpressionArg_Str{Str: "Test NL"}}}},
		}, "en")
		if err == nil {
			t.Error("should return an error")
			return
		}
		if err.Error() != "translation missing" {
			t.Errorf("unexpected error: %v", err)
			return
		}
	})

	t.Run("with single part", func(t *testing.T) {
		tr, err := getTranslation([]*studyAPI.LocalisedObject{
			{Code: "de", Parts: []*studyAPI.ExpressionArg{{Dtype: "str", Data: &studyAPI.ExpressionArg_Str{Str: "Test DE"}}}},
			{Code: "en", Parts: []*studyAPI.ExpressionArg{{Dtype: "str", Data: &studyAPI.ExpressionArg_Str{Str: "Test EN"}}}},
			{Code: "nl", Parts: []*studyAPI.ExpressionArg{{Dtype: "str", Data: &studyAPI.ExpressionArg_Str{Str: "Test NL"}}}},
		}, "en")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if tr != "Test EN" {
			t.Errorf("unexpected value: %s", tr)
			return
		}
	})

	t.Run("with multiple parts", func(t *testing.T) {
		tr, err := getTranslation([]*studyAPI.LocalisedObject{
			{Code: "de", Parts: []*studyAPI.ExpressionArg{{Dtype: "str", Data: &studyAPI.ExpressionArg_Str{Str: "Test DE"}}}},
			{Code: "en", Parts: []*studyAPI.ExpressionArg{
				{Dtype: "str", Data: &studyAPI.ExpressionArg_Str{Str: "Test "}},
				{Dtype: "exp", Data: &studyAPI.ExpressionArg_Exp{Exp: &studyAPI.Expression{}}},
				{Dtype: "str", Data: &studyAPI.ExpressionArg_Str{Str: "EN"}},
			}},
			{Code: "nl", Parts: []*studyAPI.ExpressionArg{{Dtype: "str", Data: &studyAPI.ExpressionArg_Str{Str: "Test NL"}}}},
		}, "en")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if tr != "Test EN" {
			t.Errorf("unexpected value: %s", tr)
			return
		}
	})
}
