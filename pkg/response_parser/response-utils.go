package response_parser

import (
	"errors"
	"strconv"
	"strings"
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

func timestampsToStr(ts []int64) string {
	if len(ts) == 0 {
		return ""
	}

	b := make([]string, len(ts))
	for i, v := range ts {
		b[i] = strconv.Itoa(int(v))
	}
	return strings.Join(b, ",")
}
