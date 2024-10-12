package main

import (
    "encoding/json"
    "net/http"
    "time"
)

type Request struct {
    Model  string `json:"model"`
    Prompt string `json:"prompt"`
    Stream bool   `json:"stream"`
}

type Response struct {
    Model              string    `json:"model"`
    CreatedAt          time.Time `json:"created_at"`
    Response           string    `json:"response"`
    Done               bool      `json:"done"`
    Context            []int     `json:"context"`
    TotalDuration      int64     `json:"total_duration"`
    LoadDuration       int64     `json:"load_duration"`
    PromptEvalCount    int       `json:"prompt_eval_count"`
    PromptEvalDuration int64     `json:"prompt_eval_duration"`
    EvalCount          int       `json:"eval_count"`
    EvalDuration       int64     `json:"eval_duration"`
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Decode the request body
    var req Request
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    // Create a dummy response based on the request
    response := Response{
        Model:              req.Model,
        CreatedAt:          time.Now(),
        Response:           "This is a generated response to: " + req.Prompt,
        Done:               true,
        Context:            []int{1, 2, 3}, // Dummy data
        TotalDuration:      1234,
        LoadDuration:       234,
        PromptEvalCount:    1,
        PromptEvalDuration: 100,
        EvalCount:          1,
        EvalDuration:       200,
    }

    // Set the response header and write the response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func main() {
    http.HandleFunc("/api/generate", generateHandler)
    serverAddress := ":11434"

    // Start the server
    err := http.ListenAndServe(serverAddress, nil)
    if err != nil {
        panic(err)
    }
}
