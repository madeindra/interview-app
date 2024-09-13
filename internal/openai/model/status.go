package model

type Status string

const (
	STATUS_OPERATIONAL          Status = "operational"
	STATUS_DEGRADED_PERFORMANCE Status = "degraded_performance"
	STATUS_PARTIAL_OUTAGE       Status = "partial_outage"
	STATUS_MAJOR_OUTAGE         Status = "major_outage"
	STATUS_UNKNOWN              Status = "unknown"
)

type ComponentStatusResponse struct {
	Components []Component `json:"components"`
}

type Component struct {
	Name   string `json:"name"`
	Status Status `json:"status"`
}
