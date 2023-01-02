
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// InferenceRequest represents the structure of an incoming AI inference request.
type InferenceRequest struct {
	ModelName string        `json:"model_name"`
	InputData   []float64     `json:"input_data"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// InferenceResponse represents the structure of an outgoing AI inference response.
type InferenceResponse struct {
	ModelName string        `json:"model_name"`
	Prediction  []float64     `json:"prediction"`
	Timestamp   time.Time     `json:"timestamp"`
	LatencyMs   int64         `json:"latency_ms"`
	Status      string        `json:"status"`
	Error       string        `json:"error,omitempty"`
}

// simulateInference simulates an AI model inference process.
func simulateInference(req InferenceRequest) ([]float64, error) {
	// In a real scenario, this would involve calling an actual AI model (e.g., TensorFlow Serving, PyTorch).
	// For demonstration, we'll just return a dummy prediction based on input data.

	if len(req.InputData) == 0 {
		return nil, fmt.Errorf("input data cannot be empty")
	}

	// Simple simulation: sum of input data as a single prediction value
	prediction := make([]float64, 1)
	for _, val := range req.InputData {
		prediction[0] += val
	}

	// Add some noise or complexity based on model name
	switch req.ModelName {
	case "sentiment_analysis":
		prediction[0] = prediction[0] / float64(len(req.InputData)) // Average for sentiment
		if prediction[0] > 0.5 {
			prediction = []float64{0.9} // Positive
		} else {
			prediction = []float64{0.1} // Negative
		}
	case "image_classifier":
		prediction = []float64{0.1, 0.2, 0.7} // Example class probabilities
	case "fraud_detection":
		prediction = []float64{0.01} // Low fraud probability
	default:
		// Default to the sum
	}

	return prediction, nil
}

// inferenceHandler handles incoming AI inference requests.
func inferenceHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are accepted", http.StatusMethodNotAllowed)
		return
	}

	var req InferenceRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	prediction, err := simulateInference(req)
	resp := InferenceResponse{
		ModelName: req.ModelName,
		Timestamp: time.Now(),
		Status:    "success",
	}

	if err != nil {
		resp.Status = "error"
		resp.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		resp.Prediction = prediction
		w.WriteHeader(http.StatusOK)
	}

	duration := time.Since(start)
	resp.LatencyMs = duration.Milliseconds()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func main() {
	// Setup HTTP server
	http.HandleFunc("/infer", inferenceHandler)

	port := ":8080"
	fmt.Printf("AI Inference Service starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
