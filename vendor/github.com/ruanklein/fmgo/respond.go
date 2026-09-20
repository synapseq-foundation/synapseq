// SPDX-License-Identifier: Apache-2.0

package fmgo

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

// Request configures a Foundation Models response.
type Request struct {
	Prompt         string
	Instructions   string
	Model          Model
	Text           []string
	Images         []string
	ImageLabels    []string
	Tools          []Tool
	Schema         Schema
	SchemaFile     string
	Resume         string
	SaveTranscript string
	Greedy         bool
	Verbose        bool
	UseCase        UseCase
	Guardrails     GuardrailLevel
}

// Response is a non-streaming Foundation Models response.
type Response struct {
	Text string
}

// Respond generates a complete non-streaming response.
func (c *Client) Respond(ctx context.Context, request Request) (Response, error) {
	request, cleanup, err := request.withSchemaFile()
	if err != nil {
		return Response{}, err
	}
	defer cleanup()

	args, err := request.args(false)
	if err != nil {
		return Response{}, err
	}
	stdout, _, err := c.run(ctx, args...)
	if err != nil {
		return Response{}, err
	}
	return Response{Text: string(stdout)}, nil
}

func (r Request) withSchemaFile() (Request, func(), error) {
	if len(r.Schema) == 0 {
		return r, func() {}, nil
	}
	if r.SchemaFile != "" {
		return Request{}, nil, fmt.Errorf("fmgo: schema and schema file cannot both be set")
	}

	file, err := os.CreateTemp("", "fmgo-schema-*.json")
	if err != nil {
		return Request{}, nil, fmt.Errorf("fmgo: create schema file: %w", err)
	}
	path := file.Name()
	if _, err := file.Write(r.Schema); err != nil {
		_ = file.Close()
		_ = os.Remove(path)
		return Request{}, nil, fmt.Errorf("fmgo: write schema file: %w", err)
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(path)
		return Request{}, nil, fmt.Errorf("fmgo: close schema file: %w", err)
	}

	r.Schema = nil
	r.SchemaFile = path
	return r, func() { _ = os.Remove(path) }, nil
}

// RespondAs generates JSON for T and decodes it into T.
func RespondAs[T any](ctx context.Context, client *Client, request Request) (T, error) {
	var result T
	schema, err := SchemaFor[T]()
	if err != nil {
		return result, err
	}
	request.Schema = schema
	request.SchemaFile = ""
	response, err := client.Respond(ctx, request)
	if err != nil {
		return result, err
	}
	if err := json.Unmarshal([]byte(response.Text), &result); err != nil {
		return result, fmt.Errorf("fmgo: decode structured response: %w", err)
	}
	return result, nil
}

func (r Request) args(stream bool) ([]string, error) {
	if len(r.ImageLabels) != 0 && len(r.ImageLabels) != len(r.Images) {
		return nil, fmt.Errorf("fmgo: image labels must match image count")
	}
	if len(r.Schema) != 0 && r.SchemaFile != "" {
		return nil, fmt.Errorf("fmgo: schema and schema file cannot both be set")
	}

	args := []string{"respond"}
	if stream {
		args = append(args, "--stream")
	} else {
		args = append(args, "--no-stream")
	}
	if r.Model != "" {
		args = append(args, "--model", string(r.Model))
	}
	if r.Instructions != "" {
		args = append(args, "--instructions", r.Instructions)
	}
	if len(r.Schema) != 0 {
		args = append(args, "--schema", string(r.Schema))
	}
	if r.SchemaFile != "" {
		args = append(args, "--schema", r.SchemaFile)
	}
	for _, text := range r.Text {
		args = append(args, "--text", text)
	}
	for index, image := range r.Images {
		args = append(args, "--image", image)
		if len(r.ImageLabels) != 0 {
			args = append(args, "--label", r.ImageLabels[index])
		}
	}
	for _, tool := range r.Tools {
		args = append(args, "--tool", string(tool))
	}
	if r.Resume != "" {
		args = append(args, "--resume", r.Resume)
	}
	if r.SaveTranscript != "" {
		args = append(args, "--save-transcript", r.SaveTranscript)
	}
	if r.Greedy {
		args = append(args, "--greedy")
	}
	if r.Verbose {
		args = append(args, "--verbose")
	}
	if r.UseCase != "" {
		args = append(args, "--use-case", string(r.UseCase))
	}
	if r.Guardrails != "" {
		args = append(args, "--guardrails", string(r.Guardrails))
	}
	if r.Prompt != "" {
		args = append(args, r.Prompt)
	}
	return args, nil
}
