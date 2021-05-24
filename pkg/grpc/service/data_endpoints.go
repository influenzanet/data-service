package service

import (
	"bytes"
	"context"
	"io"
	"log"

	"github.com/influenzanet/data-service/pkg/api"
	"github.com/influenzanet/data-service/pkg/response_parser"
	"github.com/influenzanet/go-utils/pkg/token_checks"
	studyAPI "github.com/influenzanet/study-service/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const chunkSize = 64 * 1024 // 64 KiB

func (s *dataServiceServer) GetResponsesCSV(req *api.ResponseQuery, stream api.DataServiceApi_GetResponsesCSVServer) error {
	if req == nil || token_checks.IsTokenEmpty(req.Token) || req.StudyKey == "" {
		return status.Error(codes.InvalidArgument, "missing argument")
	}

	surveyDef, err := s.clients.StudyService.GetSurveyDefForStudy(context.Background(), &studyAPI.SurveyReferenceRequest{
		Token:     req.Token,
		StudyKey:  req.StudyKey,
		SurveyKey: req.SurveyKey,
	})
	if err != nil {
		return nil
	}

	rp, err := response_parser.NewResponseParser(surveyDef, "ignored", req.ShortQuestionKeys, req.Separator)
	if err != nil {
		return nil
	}

	respStream, err := s.clients.StudyService.StreamStudyResponses(context.Background(), &studyAPI.SurveyResponseQuery{
		Token:     req.Token,
		StudyKey:  req.StudyKey,
		SurveyKey: req.SurveyKey,
		From:      req.From,
		Until:     req.Until,
	})
	if err != nil {
		return err
	}
	for {
		r, err := respStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("GetResponsesCSV(_) = _, %v", err)
			break
		}
		err = rp.AddResponse(r)
		if err != nil {
			log.Printf("GetResponsesCSV.AddResponse(_) = _, %v", err)
		}
	}

	buf := new(bytes.Buffer)
	err = rp.GetResponsesCSV(buf, req.IncludeMeta)
	if err != nil {
		log.Printf("GetResponsesCSV: %v", err)
		return err
	}

	chnk := &api.Chunk{}

	for currentByte := 0; currentByte < len(buf.Bytes()); currentByte += chunkSize {
		if currentByte+chunkSize > len(buf.Bytes()) {
			chnk.Chunk = buf.Bytes()[currentByte:len(buf.Bytes())]
		} else {
			chnk.Chunk = buf.Bytes()[currentByte : currentByte+chunkSize]
		}

		if err := stream.Send(chnk); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}

func (s *dataServiceServer) GetSurveyInfoCSV(req *api.SurveyInfoQuery, stream api.DataServiceApi_GetSurveyInfoCSVServer) error {
	return nil
}

func (s *dataServiceServer) GetSurveyInfo(ctx context.Context, req *api.SurveyInfoQuery) (*api.SurveyInfo, error) {
	return nil, nil
}
