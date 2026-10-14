package core

import (
	"encoding/json"
	"fmt"
	"strings"
)

type State struct {
	Goal    string
	History []Step

	Transcript []json.RawMessage
	Risks      map[string]*Risk
	nextRiskID int
	TurnStarts []int
}

func (s *State) AddRisk(hypothesis string, contract string) *Risk {
	s.nextRiskID++
	ID := fmt.Sprintf("R%d", s.nextRiskID)

	risk := Risk{
		ID:          ID,
		Hypothesis:  hypothesis,
		Contract:    contract,
		Status:      RiskSuspected,
		EvidenceIDs: []string{},
	}
	s.Risks[ID] = &risk
	return &risk
}

func (s *State) SubmitVerdict(riskID string, status string, evidenceIDs []string, reason string) error {
	_, ok := s.Risks[riskID]
	if !ok {
		return fmt.Errorf("risk %s does not exist", riskID)
	}

	switch status {
	case string(RiskConfirmed), string(RiskInconclusive), string(RiskRejected):

	default:
		return fmt.Errorf("invalid verdict status: %s", status)
	}

	if strings.TrimSpace(reason) == "" {
		return fmt.Errorf("reason must not be empty")
	}

	if (status == string(RiskConfirmed) || status == string(RiskRejected)) && len(evidenceIDs) == 0 {
		return fmt.Errorf("verdict requires evidence")
	}

	for _, id := range evidenceIDs {
		if !s.HasVerificationResult(id) {
			return fmt.Errorf("verification evidence %s not found", id)
		}
	}

	s.Risks[riskID].Status = RiskStatus(status)
	s.Risks[riskID].EvidenceIDs = append([]string{}, evidenceIDs...)
	s.Risks[riskID].VerdictReason = reason

	return nil
}

func (s *State) HasVerificationResult(id string) bool {
	for _, h := range s.History {
		for _, t := range h.Turn.ToolCalls {
			if t.CallID == id && t.Name == ToolRunTest {
				for _, result := range h.ToolResults {
					if result.CallID == id {
						return true
					}
				}
			}
		}
	}
	return false
}

func (s *State) HasSuspectedRisk() bool {
	for _, risk := range s.Risks {
		if risk.Status == RiskSuspected {
			return true
		}
	}
	return false
}
