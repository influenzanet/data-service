package response_parser

import (
	"errors"
	"log"
	"strconv"
	"strings"

	studyAPI "github.com/influenzanet/study-service/pkg/api"
)

type responseParser struct {
	surveyKey            string
	surveyVersions       []SurveyVersionPreview
	responses            []ParsedResponse
	responseColNames     []string
	metaColNames         []string
	shortQuestionKeys    bool
	questionOptionKeySep string
}

func NewResponseParser(
	surveyDef *studyAPI.Survey,
	previewLang string,
	shortQuestionKeys bool,
	questionOptionSep string,
) (*responseParser, error) {
	if surveyDef == nil || surveyDef.Current == nil || surveyDef.Current.SurveyDefinition == nil {
		return nil, errors.New("current survey definition not found")
	}

	rp := responseParser{
		surveyKey:            surveyDef.Current.SurveyDefinition.Key,
		surveyVersions:       []SurveyVersionPreview{},
		responses:            []ParsedResponse{},
		shortQuestionKeys:    shortQuestionKeys,
		questionOptionKeySep: questionOptionSep,
	}

	rp.surveyVersions = append(rp.surveyVersions, surveyDefToVersionPreview(surveyDef.Current, previewLang))
	for _, v := range surveyDef.History {
		rp.surveyVersions = append(rp.surveyVersions, surveyDefToVersionPreview(v, previewLang))
	}

	for versionInd, sv := range rp.surveyVersions {
		for qInd, question := range sv.Questions {
			if shortQuestionKeys {
				rp.surveyVersions[versionInd].Questions[qInd].ID = strings.TrimPrefix(question.ID, rp.surveyKey+".")
			}

			//if shortResponseKeys {
			/*for rInd, resp := range question.Responses {
				rIDparts := strings.Split(resp.ID, ".")
				rp.surveyVersions[versionInd].Questions[qInd].Responses[rInd].ID = rIDparts[len(rIDparts)-1]

				for oInd, option := range resp.Options {
					oIDparts := strings.Split(option.ID, ".")
					rp.surveyVersions[versionInd].Questions[qInd].Responses[rInd].Options[oInd].ID = oIDparts[len(oIDparts)-1]
				}
			}*/

		}

	}
	log.Println(rp.surveyVersions)
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

	if rp.shortQuestionKeys {
		for i, r := range rawResp.Responses {
			rawResp.Responses[i].Key = strings.TrimPrefix(r.Key, rp.surveyKey+".")
		}
	}

	for _, question := range currentVersion.Questions {
		log.Println(question)

		resp := findResponse(rawResp.Responses, question.ID)

		// TODO: parse question
		responseColumns := getResponseColumns(question, resp, rp.questionOptionKeySep)
		for k, v := range responseColumns {
			parsedResponse.Responses[k] = v
		}

		// Set meta infos
		parsedResponse.Meta.Initialised[question.ID] = ""
		parsedResponse.Meta.Displayed[question.ID] = ""
		parsedResponse.Meta.Responded[question.ID] = ""
		parsedResponse.Meta.ItemVersion[question.ID] = ""
		arraySep := ";"
		if resp != nil && resp.Meta != nil {
			parsedResponse.Meta.Initialised[question.ID] = timestampsToStr(resp.Meta.Rendered, arraySep)
			parsedResponse.Meta.Displayed[question.ID] = timestampsToStr(resp.Meta.Displayed, arraySep)
			parsedResponse.Meta.Responded[question.ID] = timestampsToStr(resp.Meta.Responded, arraySep)
			parsedResponse.Meta.ItemVersion[question.ID] = strconv.Itoa(int(resp.Meta.Version))
		}
		rp.AddMetaColName(question.ID)
	}

	// Extend response col names:
	for k := range parsedResponse.Responses {
		rp.AddResponseColName(k)
	}

	rp.responses = append(rp.responses, parsedResponse)
	return errors.New("unimplemented")
}

func (rp *responseParser) AddResponseColName(name string) {
	for _, n := range rp.responseColNames {
		if n == name {
			return
		}
	}
	rp.responseColNames = append(rp.responseColNames, name)
}

func (rp *responseParser) AddMetaColName(name string) {
	for _, n := range rp.metaColNames {
		if n == name {
			return
		}
	}
	rp.metaColNames = append(rp.metaColNames, name)
}

func (rp responseParser) GetSurveyVersionDefs() []SurveyVersionPreview {
	return rp.surveyVersions
}

func (rp responseParser) GetResponsesCSV() []ParsedResponse {
	// TODO: merge context and responses keys
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
