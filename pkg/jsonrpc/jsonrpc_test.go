package jsonrpc

import (
	"context"
	"testing"
)

func TestParse_Valid(t *testing.T) {
	raw := `{"method":"query","parameters":["secret"]}`
	req, err := Parse(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Method != "query" {
		t.Errorf("got Method=%s, want query", req.Method)
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	_, err := Parse(`not json`)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParse_MissingMethod(t *testing.T) {
	_, err := Parse(`{"parameters":[]}`)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParse_EmptyMethod(t *testing.T) {
	_, err := Parse(`{"method":""}`)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParse_ParamsNotArray(t *testing.T) {
	_, err := Parse(`{"method":"query","parameters":"string"}`)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParse_UnknownMethod(t *testing.T) {
	raw := `{"method":"unknown","parameters":[]}`
	req, err := Parse(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Method != "unknown" {
		t.Errorf("got Method=%s, want unknown", req.Method)
	}
}

func TestParse_ExecuteMethod(t *testing.T) {
	raw := `{"method":"execute_secret","parameters":["32"]}`
	req, err := Parse(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Method != "execute_secret" {
		t.Errorf("got Method=%s, want execute_secret", req.Method)
	}
}

func TestDispatcher(t *testing.T) {
	d := New()
	d.Handle("ping", func(ctx context.Context, req Request) (any, *Error) {
		return "pong", nil
	})

	result, err := d.Dispatch(context.Background(), `{"method":"ping"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "pong" {
		t.Errorf("got %s, want pong", result)
	}
}

func TestDispatcher_UnknownMethod(t *testing.T) {
	d := New()
	_, err := d.Dispatch(context.Background(), `{"method":"unknown"}`)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDispatcher_JSONResponse(t *testing.T) {
	d := New()
	d.Handle("query", func(ctx context.Context, req Request) (any, *Error) {
		return map[string]any{"result": []map[string]any{
			{"Title": "Hello"},
		}}, nil
	})

	result, err := d.Dispatch(context.Background(), `{"method":"query","parameters":["test"]}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == "" {
		t.Fatal("expected non-empty response")
	}
	// Should be direct JSON, not wrapped in jsonrpc envelope
	if result != `{"result":[{"Title":"Hello"}]}` {
		t.Errorf("got %s", result)
	}
}