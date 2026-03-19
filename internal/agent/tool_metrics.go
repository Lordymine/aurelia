package agent

import "encoding/json"

const OversizedToolOutputThresholdChars = 8000

type ToolPayloadMetrics struct {
	Count           int
	SerializedBytes int
}

type ToolOutputMetrics struct {
	Chars     int
	Oversized bool
}

func MeasureToolPayload(tools []Tool) ToolPayloadMetrics {
	metrics := ToolPayloadMetrics{Count: len(tools)}
	if len(tools) == 0 {
		return metrics
	}

	payload, err := json.Marshal(tools)
	if err != nil {
		return metrics
	}

	metrics.SerializedBytes = len(payload)
	return metrics
}

func MeasureToolOutput(content string) ToolOutputMetrics {
	chars := len([]rune(content))
	return ToolOutputMetrics{
		Chars:     chars,
		Oversized: chars > OversizedToolOutputThresholdChars,
	}
}
