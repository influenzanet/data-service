package response_parser

import (
	"errors"
	"log"
	"strconv"
	"strings"

	studyAPI "github.com/influenzanet/study-service/pkg/api"
)

type responseParser struct {
	surveyKey         string
	surveyVersions    []SurveyVersionPreview
	responses         []ParsedResponse
	shortQuestionKeys bool
	shortResponseKeys bool
}

func NewResponseParser(
	surveyDef *studyAPI.Survey,
	previewLang string,
	shortQuestionKeys bool,
	shortResponseKeys bool,
) (*responseParser, error) {
	if surveyDef == nil || surveyDef.Current == nil || surveyDef.Current.SurveyDefinition == nil {
		return nil, errors.New("current survey definition not found")
	}

	rp := responseParser{
		surveyKey:         surveyDef.Current.SurveyDefinition.Key,
		surveyVersions:    []SurveyVersionPreview{},
		responses:         []ParsedResponse{},
		shortQuestionKeys: shortQuestionKeys,
		shortResponseKeys: shortResponseKeys,
	}

	rp.surveyVersions = append(rp.surveyVersions, surveyDefToVersionPreview(surveyDef.Current, previewLang))
	for _, v := range surveyDef.History {
		rp.surveyVersions = append(rp.surveyVersions, surveyDefToVersionPreview(v, previewLang))
	}

	if shortQuestionKeys && shortResponseKeys {
		for versionInd, sv := range rp.surveyVersions {
			for qInd, question := range sv.Questions {
				if shortQuestionKeys {
					rp.surveyVersions[versionInd].Questions[qInd].ID = strings.TrimPrefix(question.ID, rp.surveyKey+".")
				}

				if shortResponseKeys {
					for rInd, resp := range question.Responses {
						rIDparts := strings.Split(resp.ID, ".")
						rp.surveyVersions[versionInd].Questions[qInd].Responses[rInd].ID = rIDparts[len(rIDparts)-1]

						for oInd, option := range resp.Options {
							oIDparts := strings.Split(option.ID, ".")
							rp.surveyVersions[versionInd].Questions[qInd].Responses[rInd].Options[oInd].ID = oIDparts[len(oIDparts)-1]
						}
					}
				}
			}

		}
	}
	log.Println(rp.GetSurveyVersionDefs())

	return &rp, errors.New("test")
}

func (rp *responseParser) AddResponse(rawResp *studyAPI.SurveyResponse) error {
	parsedResponse := ParsedResponse{
		ParticipantID: rawResp.ParticipantId,
		Version:       rawResp.VersionId,
		SubmittedAt:   rawResp.SubmittedAt,
		Context:       rawResp.Context,
	}

	currentVersion, err := findSurveyVersion(rawResp.VersionId, rawResp.SubmittedAt, rp.surveyVersions)
	if err != nil {
		return err
	}

	// TODO: interpret response  from DB
	log.Println(currentVersion)

	key := "test"
	index := 0
	parsedResponse.Meta.Initialised[key] = timestampsToStr(rawResp.Responses[index].Meta.Rendered)
	parsedResponse.Meta.Displayed[key] = timestampsToStr(rawResp.Responses[index].Meta.Displayed)
	parsedResponse.Meta.Responded[key] = timestampsToStr(rawResp.Responses[index].Meta.Responded)
	parsedResponse.Meta.ItemVersion[key] = strconv.Itoa(int(rawResp.Responses[index].Meta.Version))

	rp.responses = append(rp.responses, parsedResponse)
	return errors.New("unimplemented")
}

func (rp responseParser) GetSurveyVersionDefs() []SurveyVersionPreview {
	return rp.surveyVersions
}

func (rp responseParser) GetResponses() []ParsedResponse {
	return rp.responses
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
