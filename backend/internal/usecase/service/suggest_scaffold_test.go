package usecase

import (
	"context"
	"testing"
)

func TestSuggestScaffold(t *testing.T) {
	usecase := &serviceUsecase{}

	suggestion, err := usecase.SuggestScaffold(context.Background(), SuggestScaffoldInput{
		ServiceName:        "Payments API",
		ProjectName:        "Platform",
		ProjectDescription: "Payment platform using Postgres with structured logging.",
		RepoURL:            "git@gitea.devhub.local:platform/payments-api.git",
		Environments:       []string{"staging", "prod"},
	})
	if err != nil {
		t.Fatalf("SuggestScaffold returned error: %v", err)
	}

	if suggestion.Source != "local-heuristic-v1" {
		t.Fatalf("expected local heuristic source, got %q", suggestion.Source)
	}
	if suggestion.Environment != "staging" {
		t.Fatalf("expected first project environment, got %q", suggestion.Environment)
	}
	if suggestion.Variables.ServiceName != "payments-api" {
		t.Fatalf("expected normalized service name, got %q", suggestion.Variables.ServiceName)
	}
	if suggestion.Variables.ModulePath != "gitea.devhub.local/platform/payments-api" {
		t.Fatalf("expected inferred module path, got %q", suggestion.Variables.ModulePath)
	}
	if suggestion.Variables.Port == 0 {
		t.Fatal("expected suggested port")
	}
	if !suggestion.Variables.EnableLogging {
		t.Fatal("expected logging to be enabled")
	}
}

func TestSuggestScaffoldUsesSelectedEnvironment(t *testing.T) {
	usecase := &serviceUsecase{}

	suggestion, err := usecase.SuggestScaffold(context.Background(), SuggestScaffoldInput{
		ServiceName:  "Cache",
		Environment:  "prod",
		Environments: []string{"dev", "staging"},
	})
	if err != nil {
		t.Fatalf("SuggestScaffold returned error: %v", err)
	}

	if suggestion.Environment != "prod" {
		t.Fatalf("expected selected environment, got %q", suggestion.Environment)
	}
	if suggestion.Variables.Database != "redis" {
		t.Fatalf("expected cache service database hint, got %q", suggestion.Variables.Database)
	}
}

func TestSuggestScaffoldInfersFromProjectDescription(t *testing.T) {
	usecase := &serviceUsecase{}

	suggestion, err := usecase.SuggestScaffold(context.Background(), SuggestScaffoldInput{
		ProjectDescription: "Customer document search API backed by MongoDB with no logging.",
	})
	if err != nil {
		t.Fatalf("SuggestScaffold returned error: %v", err)
	}

	if suggestion.Variables.ServiceName != "customer-document-search" {
		t.Fatalf("expected service name from description, got %q", suggestion.Variables.ServiceName)
	}
	if suggestion.Variables.Database != "mongodb" {
		t.Fatalf("expected database from description, got %q", suggestion.Variables.Database)
	}
	if suggestion.Variables.EnableLogging {
		t.Fatal("expected logging disabled from description")
	}
}
