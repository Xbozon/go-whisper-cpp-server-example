package whisper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type Config struct {
	Temperature    float32
	TemperatureInc float32
	Timeout        time.Duration
}

type Response struct {
	Text string `json:"text"`
}

type ServerApi struct {
	url    string
	config Config
}

func NewServerApi(url string, config Config) *ServerApi {
	return &ServerApi{url: url, config: config}
}

func (s *ServerApi) SendMultiPartForm(ctx context.Context, wavData []byte) (Response, error) {
	// Create a buffer to hold the multipart form data
	var b bytes.Buffer
	multipartWriter := multipart.NewWriter(&b)

	// Create a form file part
	part, err := multipartWriter.CreateFormFile("file", "example.wav")
	if err != nil {
		return Response{}, fmt.Errorf("creating multipart file form: %w", err)
	}

	// Write WAV data to the form file part
	_, err = part.Write(wavData)
	if err != nil {
		return Response{}, fmt.Errorf("write data to multipart writer: %w", err)
	}

	// Add a form field for the response format
	err = s.writeConfig(multipartWriter)
	if err != nil {
		return Response{}, fmt.Errorf("write whisper config: %w", err)
	}

	// Close the multipart writer to finalize the boundary
	err = multipartWriter.Close()
	if err != nil {
		return Response{}, fmt.Errorf("multipart writer close: %w", err)
	}

	// Create the HTTP send
	request, err := http.NewRequestWithContext(ctx, "POST", s.url, &b)
	if err != nil {
		return Response{}, fmt.Errorf("create request with context: %w", err)
	}
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	return s.send(request)
}

func (s *ServerApi) send(request *http.Request) (Response, error) {
	// Perform the request
	client := &http.Client{Timeout: s.config.Timeout}

	response, err := client.Do(request)
	if err != nil {
		return Response{}, err // Handle the error appropriately
	}
	defer response.Body.Close()

	// Check response status
	if response.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("server responded with status code: %d", response.StatusCode)
	}

	// read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Response{}, fmt.Errorf("%w", err)
	}

	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		return Response{}, fmt.Errorf("body unmarshal: %w", err)
	}

	return result, nil
}

func (s *ServerApi) writeConfig(mw *multipart.Writer) error {
	// Add a form field for the response format
	err := mw.WriteField("response_format", "json")
	if err != nil {
		return fmt.Errorf("add response_format: %w", err)
	}

	err = mw.WriteField("temperature", fmt.Sprintf("%.2f", s.config.Temperature))
	if err != nil {
		return fmt.Errorf("add temperature: %w", err)
	}

	err = mw.WriteField("temperature_inc", fmt.Sprintf("%.2f", s.config.TemperatureInc))
	if err != nil {
		return fmt.Errorf("add temperature_inc: %w", err)
	}

	return nil
}
