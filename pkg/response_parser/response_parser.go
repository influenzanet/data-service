package response_parser

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	studyAPI "github.com/influenzanet/study-service/pkg/api"
)

type ResponseParser struct {
	surveyKey            string
	surveyVersions       []SurveyVersionPreview
	responses            []ParsedResponse
	contextColNames      []string
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
) (*ResponseParser, error) {
	if surveyDef == nil || surveyDef.Current == nil || surveyDef.Current.SurveyDefinition == nil {
		return nil, errors.New("current survey definition not found")
	}

	rp := ResponseParser{
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
		}

	}
	return &rp, nil
}

func (rp *ResponseParser) AddResponse(rawResp *studyAPI.SurveyResponse) error {
	parsedResponse := ParsedResponse{
		ParticipantID: rawResp.ParticipantId,
		Version:       rawResp.VersionId,
		SubmittedAt:   rawResp.SubmittedAt,
		Context:       rawResp.Context,
		Responses:     map[string]string{},
		Meta: ResponseMeta{
			Initialised: map[string]string{},
			Displayed:   map[string]string{},
			Responded:   map[string]string{},
			ItemVersion: map[string]string{},
		},
	}

	currentVersion, err := findSurveyVersion(rawResp.VersionId, rawResp.SubmittedAt, rp.surveyVersions)
	if err != nil {
		return err
	}

	if rp.shortQuestionKeys {
		for i, r := range rawResp.Responses {
			rawResp.Responses[i].Key = strings.TrimPrefix(r.Key, rp.surveyKey+".")
		}
	}

	for _, question := range currentVersion.Questions {
		resp := findResponse(rawResp.Responses, question.ID)

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
	for k := range parsedResponse.Context {
		rp.AddContextColName(k)
	}

	rp.responses = append(rp.responses, parsedResponse)
	return nil
}

func (rp *ResponseParser) AddResponseColName(name string) {
	for _, n := range rp.responseColNames {
		if n == name {
			return
		}
	}
	rp.responseColNames = append(rp.responseColNames, name)
}

func (rp *ResponseParser) AddContextColName(name string) {
	for _, n := range rp.contextColNames {
		if n == name {
			return
		}
	}
	rp.contextColNames = append(rp.contextColNames, name)
}

func (rp *ResponseParser) AddMetaColName(name string) {
	for _, n := range rp.metaColNames {
		if n == name {
			return
		}
	}
	rp.metaColNames = append(rp.metaColNames, name)
}

func (rp ResponseParser) GetSurveyVersionDefs() []SurveyVersionPreview {
	return rp.surveyVersions
}

func (rp ResponseParser) GetResponses() []ParsedResponse {
	// TODO: merge context and responses keys
	return rp.responses
}

func (rp ResponseParser) GetResponsesCSV(writer io.Writer) error {
	if len(rp.responses) < 1 {
		return errors.New("no responses, nothing is generated")
	}

	// Sort column names
	contextCols := rp.contextColNames
	sort.Strings(contextCols)
	responseCols := rp.responseColNames
	sort.Strings(responseCols)

	// Prepare csv header
	header := []string{
		"participantID",
		"version",
		"submitted",
	}
	header = append(header, contextCols...)
	header = append(header, responseCols...)

	// Init writer
	w := csv.NewWriter(writer)

	// Write header
	err := w.Write(header)
	if err != nil {
		return err
	}

	// Write responses
	for _, resp := range rp.responses {
		line := []string{
			resp.ParticipantID,
			resp.Version,
			fmt.Sprint(resp.SubmittedAt),
		}

		for _, colName := range contextCols {
			v, ok := resp.Context[colName]
			if !ok {
				line = append(line, "")
				continue
			}
			line = append(line, v)
		}

		for _, colName := range responseCols {
			v, ok := resp.Responses[colName]
			if !ok {
				line = append(line, "")
				continue
			}
			line = append(line, v)
		}

		err := w.Write(line)
		if err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}

func (rp ResponseParser) GetMetaCSV() string {
	return ""
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
