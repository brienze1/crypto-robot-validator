package dto

import (
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/summary"
)

type AnalysisDto struct {
	Summary   summary.Summary `json:"summary"`
	Timestamp string          `json:"timestamp"`
}
