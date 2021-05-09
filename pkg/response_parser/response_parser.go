package response_parser

import (
	"errors"

	studyAPI "github.com/influenzanet/study-service/pkg/api"
)

type responseParser struct {
	surveyVersions []SurveyVersionPreview
	responses      []ParsedResponse
}

func NewResponseParser(
	surveyDef *studyAPI.Survey,
	previewLang string,
	shortQuestionKeys bool,
	shortResponseKeys bool,
) (*responseParser, error) {
	if surveyDef == nil || surveyDef.Current == nil {
		return nil, errors.New("current survey definition not found")
	}

	rp := responseParser{
		surveyVersions: []SurveyVersionPreview{},
		responses:      []ParsedResponse{},
	}

	rp.surveyVersions = append(rp.surveyVersions, surveyDefToVersionPreview(surveyDef.Current))
	for _, v := range surveyDef.History {
		rp.surveyVersions = append(rp.surveyVersions, surveyDefToVersionPreview(v))
	}

	return &rp, errors.New("unimplemented")
}

func (rp *responseParser) Parse() error {
	return errors.New("unimplemented")
}

func (rp responseParser) GetSurveyDef() error {
	return errors.New("unimplemented")
}

func (rp responseParser) GetResponses() error {
	return errors.New("unimplemented")
}

func (rp responseParser) GetMeta() error {
	return errors.New("unimplemented")
}

/*
func getResponseTableCSV(responses studyAPI.SurveyResponses) (string, error) {

	for _, resp := range responses.Responses {
		// resp.VersionId
		for _, question := range resp.Responses {
			question.Meta.Version
		}
	}

}*/
