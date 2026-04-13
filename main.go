package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// InferenceRequest represents the structure of an incoming inference request.
type InferenceRequest struct {
	InputData json.RawMessage `json:"input_data"` // Flexible input data
}

// InferenceResponse represents the structure of an outgoing inference response.
type InferenceResponse struct {
	Prediction  json.RawMessage `json:"prediction"`
	ModelName   string          `json:"model_name"`
	InferenceTime string          `json:"inference_time"`
	Status      string          `json:"status"`
	Message     string          `json:"message,omitempty"`
}

// simulateModelInference simulates an AI model inference process.
// In a real application, this would involve loading a model and running prediction.
func simulateModelInference(inputData json.RawMessage) (json.RawMessage, error) {
	// For demonstration, we'll just echo the input as a 'prediction' after a delay.
	// In a real scenario, you'd load your ONNX, TensorFlow Lite, or other model here
	// and perform actual inference.

	// Simulate some processing time
	time.Sleep(50 * time.Millisecond)

	// Example: If input is {"features": [1, 2, 3]}, prediction could be {"class": "A"}
	// For simplicity, we'll just return a dummy prediction or process the input.
	var dummyPrediction map[string]interface{}
	json.Unmarshal(inputData, &dummyPrediction)
	
	// A very basic simulation: if input has a 'text' field, predict sentiment.
	// Otherwise, just return a generic success.
	if text, ok := dummyPrediction["text"]; ok {
		// Simple sentiment logic
		if t, isString := text.(string); isString && len(t) > 10 && t[0] == 'g' {
			dummyPrediction["sentiment"] = "positive"
		} else {
			dummyPrediction["sentiment"] = "neutral"
		}
	} else {
		dummyPrediction["result"] = "simulated_success"
	}

	predictedOutput, err := json.Marshal(dummyPrediction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal simulated prediction: %w", err)
	}
	return predictedOutput, nil
}

// inferenceHandler handles incoming AI inference requests.
func inferenceHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Set response header
	w.Header().Set("Content-Type", "application/json")

	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading request body: %v", err), http.StatusInternalServerError)
		return
	}

	// Parse the request
	var req InferenceRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing request JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Simulate model inference
	prediction, err := simulateModelInference(req.InputData)
	if err != nil {
		log.Printf("Error during simulated inference: %v", err)
		resp := InferenceResponse{
			ModelName:   "simulated-model",
			Status:      "error",
			Message:     fmt.Sprintf("Inference failed: %v", err),
			InferenceTime: time.Since(start).String(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Prepare response
	resp := InferenceResponse{
		Prediction:  prediction,
		ModelName:   "simulated-model",
		InferenceTime: time.Since(start).String(),
		Status:      "success",
	}

	// Send response
	json.NewEncoder(w).Encode(resp)
	log.Printf("Handled inference request in %s", time.Since(start))
}

func main() {
	log.Println("Starting Go AI Inference Service...")
	log.Println("Service will listen on :8080")

	// Register handler for the /predict endpoint
	http.HandleFunc("/predict", inferenceHandler)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
