// Copyright 2020-2026 ONDEWO GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Everything in this file names the ONDEWO SURVEY API specifically: its two services, one of its
// messages, its enums. It is the ONLY product-specific file of the suite -
// generated_code_test.go and auth_test.go are product agnostic.
package tests

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	survey "github.com/ondewo/ondewo-survey-client-go/v2/api/ondewo/survey"
)

// protoFileCount is the number of .proto files below ondewo-survey-api/ondewo that the compiler
// consumed. Every one of them has to end up in the global descriptor registry when this package
// is linked; a proto that silently stopped being compiled is otherwise invisible until a
// consumer misses a type.
const protoFileCount = 2

// services is every gRPC service this product exposes, keyed by the fully qualified proto name
// the ServiceDesc must declare.
var services = map[string]*grpc.ServiceDesc{
	"ondewo.survey.Surveys": &survey.Surveys_ServiceDesc,
	"ondewo.survey.FHIR":    &survey.FHIR_ServiceDesc,
}

// clientConstructors is the generated New<Service>Client of every service above. A client SDK
// that compiles but whose constructors are missing is useless, and the two generators that
// produce them (protoc-gen-go, protoc-gen-go-grpc) can disagree - so both halves are listed.
var clientConstructors = map[string]func(grpc.ClientConnInterface) any{
	"ondewo.survey.Surveys": func(cc grpc.ClientConnInterface) any { return survey.NewSurveysClient(cc) },
	"ondewo.survey.FHIR":    func(cc grpc.ClientConnInterface) any { return survey.NewFHIRClient(cc) },
}

// expectedMethods pins RPCs by name. The descriptor cross-check in generated_code_test.go proves
// the two generators agree with each other; it cannot notice an RPC that was renamed upstream,
// because both halves would be renamed together. These are spelled out so that a rename is a
// failing test rather than a silently broken consumer. The names are the ones on the wire (the
// `MethodName` of the ServiceDesc), which for FHIR differ in case from the go identifiers.
var expectedMethods = map[string][]string{
	"ondewo.survey.Surveys": {
		"CreateSurvey", "GetSurvey", "UpdateSurvey", "DeleteSurvey", "ListSurveys",
		"GetSurveyAnswers", "GetAllSurveyAnswers",
		"CreateAgentSurvey", "UpdateAgentSurvey", "DeleteAgentSurvey",
	},
	"ondewo.survey.FHIR": {"CreateFHIRSurvey", "GetFHIRSurveyAnswers", "GetAllFHIRSurveyAnswers"},
}

// TestMessageRoundTripsThroughTheWire is the core assertion about generated message code: a value
// built in go, serialized and parsed back is the same value. It covers scalars, a nested message,
// a repeated message with a oneof inside it and a packed repeated enum, so a generator that
// mis-numbers a field or loses a nested type fails here.
func TestMessageRoundTripsThroughTheWire(t *testing.T) {
	t.Parallel()

	original := &survey.Survey{
		SurveyId:     "projects/4c9a1f2e/agent",
		DisplayName:  "Patient satisfaction",
		LanguageCode: "de",
		Questions: []*survey.Question{
			{Question: &survey.Question_OpenQuestion{
				OpenQuestion: &survey.OpenQuestion{QuestionText: "How was your visit?"},
			}},
			{Question: &survey.Question_ScaleQuestion{
				ScaleQuestion: &survey.ScaleQuestion{
					QuestionText: "How likely are you to recommend us?",
					MinValue:     &survey.ScaleQuestion_ScaleValue{Value: 1, Label: "not at all"},
					MaxValue:     &survey.ScaleQuestion_ScaleValue{Value: 5, Label: "very likely"},
				},
			}},
		},
		SurveyInfo: &survey.SurveyInfo{
			LegalEntity:      "ONDEWO GmbH",
			PostalAddress:    "Seilergasse 16, 1010 Vienna",
			EmailAddress:     "office@ondewo.com",
			PhoneNumber:      "+4312345678",
			PhoneHours:       "from nine to five",
			ExpectedDuration: "about three minutes",
			Purpose:          "to improve our service",
			Topic:            "satisfaction",
		},
		ExcludeSubflows: []survey.SubFlow{survey.SubFlow_PHONE_HOURS, survey.SubFlow_PURPOSE},
		Status:          survey.Survey_UPDATED,
	}

	wire, err := proto.Marshal(original)
	if err != nil {
		t.Fatalf("proto.Marshal(%T) failed: %v", original, err)
	}
	if len(wire) == 0 {
		t.Fatal("proto.Marshal produced 0 bytes for a fully populated message")
	}

	parsed := &survey.Survey{}
	if err := proto.Unmarshal(wire, parsed); err != nil {
		t.Fatalf("proto.Unmarshal failed: %v", err)
	}

	if !proto.Equal(original, parsed) {
		t.Fatalf("round trip changed the message:\n original = %v\n  parsed = %v", original, parsed)
	}
	if got, want := len(parsed.GetQuestions()), 2; got != want {
		t.Fatalf("parsed %d questions, want %d", got, want)
	}
	if got, want := parsed.GetQuestions()[0].GetOpenQuestion().GetQuestionText(), "How was your visit?"; got != want {
		t.Errorf("oneof branch after round trip = %q, want %q", got, want)
	}
	if got, want := parsed.GetSurveyInfo().GetLegalEntity(), "ONDEWO GmbH"; got != want {
		t.Errorf("nested message field after round trip = %q, want %q", got, want)
	}
	if got, want := parsed.GetExcludeSubflows()[1], survey.SubFlow_PURPOSE; got != want {
		t.Errorf("packed repeated enum after round trip = %v, want %v", got, want)
	}
}

// NOTE ON THE MISSING PRESENCE TEST. The sibling go clients assert that a proto3 `optional` scalar
// keeps the difference between "set to the zero value" and "not set" - the distinction the angular
// target of the same compiler once lost. The ONDEWO SURVEY API declares NO `optional` field at all
// (the `oneof` wrappers in survey.proto are real oneofs, a different field kind), so there is
// nothing here to assert it against and the test is deliberately absent rather than faked. Add it
// the moment an `optional` scalar appears upstream; `grep -rn 'proto3,oneof" json:' api/` finds one.

// TestEnumZeroValueIsTheUnspecifiedMember checks the member every proto3 enum must have at 0 and
// the name maps generated beside it. A zero value that is a real choice rather than
// "unspecified" is unrequestable in several of the other clients of this API - which is exactly
// what Survey_AgentStatus is, so it is pinned by name too rather than being wished away.
func TestEnumZeroValueIsTheUnspecifiedMember(t *testing.T) {
	t.Parallel()

	var zero survey.SubFlow

	if zero != survey.SubFlow_SUBFLOW_UNSPECIFIED {
		t.Errorf("zero value of SubFlow = %v, want SUBFLOW_UNSPECIFIED", zero)
	}
	if got, want := zero.String(), "SUBFLOW_UNSPECIFIED"; got != want {
		t.Errorf("SubFlow(0).String() = %q, want %q", got, want)
	}
	if got, want := survey.SubFlow_name[0], "SUBFLOW_UNSPECIFIED"; got != want {
		t.Errorf("SubFlow_name[0] = %q, want %q", got, want)
	}
	if got, want := survey.SubFlow_value["PURPOSE"], int32(survey.SubFlow_PURPOSE); got != want {
		t.Errorf("SubFlow_value[PURPOSE] = %d, want %d", got, want)
	}
	if got, want := int32(survey.SubFlow_PHONE_HOURS), int32(6); got != want {
		t.Errorf("PHONE_HOURS = %d, want %d", got, want)
	}

	// Survey.AgentStatus has no *_UNSPECIFIED member: its zero value is the real state
	// TO_BE_INITIALIZED. Pinned here so a member inserted at 0 upstream - which would silently
	// re-label every stored zero - fails the build instead.
	var status survey.Survey_AgentStatus
	if status != survey.Survey_TO_BE_INITIALIZED {
		t.Errorf("zero value of Survey_AgentStatus = %v, want TO_BE_INITIALIZED", status)
	}
	if got, want := survey.Survey_AgentStatus_name[0], "TO_BE_INITIALIZED"; got != want {
		t.Errorf("Survey_AgentStatus_name[0] = %q, want %q", got, want)
	}
}

// TestUnmarshalRejectsTruncatedInput asserts the generated message reports a parse error instead
// of accepting a malformed payload: field 1 (`survey_id`) is announced as 5 bytes long but only 1
// follows.
func TestUnmarshalRejectsTruncatedInput(t *testing.T) {
	t.Parallel()

	if err := proto.Unmarshal([]byte{0x0a, 0x05, 'a'}, &survey.Survey{}); err == nil {
		t.Fatal("proto.Unmarshal accepted a truncated payload, want an error")
	}
}

// surveysServer is a fake ONDEWO server: it answers GetSurvey and inherits the "unimplemented"
// behaviour of the generated base type for every other RPC of the service.
type surveysServer struct {
	survey.UnimplementedSurveysServer
}

func (surveysServer) GetSurvey(_ context.Context, req *survey.GetSurveyRequest) (*survey.Survey, error) {
	return &survey.Survey{
		SurveyId:    req.GetSurveyId(),
		DisplayName: "Patient satisfaction",
		Status:      survey.Survey_UPDATED,
	}, nil
}

// TestUnaryRPCRoundTripsOverAnInProcessServer drives the generated client stub, the generated
// server stub and the generated ServiceDesc against each other over a real gRPC connection - the
// request is marshalled, routed by the method name baked into the stub, and the response is
// parsed back. Nothing here is mocked except the transport, which is in memory.
func TestUnaryRPCRoundTripsOverAnInProcessServer(t *testing.T) {
	t.Parallel()

	conn := dialInProcess(t, nil, func(srv *grpc.Server) {
		survey.RegisterSurveysServer(srv, surveysServer{})
	})
	client := survey.NewSurveysClient(conn)

	const surveyID = "projects/4c9a1f2e/agent"
	response, err := client.GetSurvey(t.Context(), &survey.GetSurveyRequest{SurveyId: surveyID})
	if err != nil {
		t.Fatalf("GetSurvey failed: %v", err)
	}

	if got := response.GetSurveyId(); got != surveyID {
		t.Errorf("response survey id = %q, want %q", got, surveyID)
	}
	if got, want := response.GetStatus(), survey.Survey_UPDATED; got != want {
		t.Errorf("response status = %v, want %v", got, want)
	}
}

// TestUnimplementedMethodIsReportedAsUnimplemented pins the other half of the generated server
// contract: an RPC the server does not implement must come back as codes.Unimplemented, not as a
// routing failure or a panic. It also proves the method is routed at all - a method missing from
// the ServiceDesc would surface as a different code.
func TestUnimplementedMethodIsReportedAsUnimplemented(t *testing.T) {
	t.Parallel()

	conn := dialInProcess(t, nil, func(srv *grpc.Server) {
		survey.RegisterSurveysServer(srv, surveysServer{})
	})
	client := survey.NewSurveysClient(conn)

	_, err := client.DeleteSurvey(t.Context(), &survey.DeleteSurveyRequest{})
	if got := status.Code(err); got != codes.Unimplemented {
		t.Fatalf("DeleteSurvey returned code %v (err = %v), want %v", got, err, codes.Unimplemented)
	}
}
