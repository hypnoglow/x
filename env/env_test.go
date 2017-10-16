package env

import (
	"os"
	"testing"
)

func TestMust(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("ENV_VAR", "hello")

		value := Must("ENV_VAR")
		if value != "hello" {
			t.Fatalf("Expected value to be %q but got %q", "hello", value)
		}
	})

	t.Run("ok for prefixed with $", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("ENV_VAR", "hello")

		value := Must("$ENV_VAR")
		if value != "hello" {
			t.Fatalf("Expected value to be %q but got %q", "hello", value)
		}
	})

	t.Run("panics on non-existent env var", func(t *testing.T) {
		os.Clearenv()
		defer func() {
			r := recover()
			if r == nil {
				t.Fatalf("Expected panic")
			}
		}()

		_ = Must("ENV_VAR")
	})
}

func TestMustBool(t *testing.T) {
	t.Run("ok for true", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("ENV_VAR", "true")

		value := MustBool("ENV_VAR")
		if !value {
			t.Fatalf("Expected value to be %v but got %v", true, value)
		}
	})

	t.Run("ok for false", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("ENV_VAR", "false")

		value := MustBool("ENV_VAR")
		if value {
			t.Fatalf("Expected value to be %v but got %v", false, value)
		}
	})

	t.Run("panics on values that are not in [true, false]", func(t *testing.T) {
		os.Clearenv()
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("Expected panic")
			}
		}()

		os.Setenv("ENV_VAR", "some")
		_ = MustBool("ENV_VAR")
	})
}

func TestGet(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("ENV_VAR", "hello")

		value := Get("ENV_VAR", "world")
		if value != "hello" {
			t.Fatalf("Expected value to be %v but got %v", "hello", value)
		}
	})

	t.Run("ok with default", func(t *testing.T) {
		os.Clearenv()
		value := Get("ENV_VAR", "world")
		if value != "world" {
			t.Fatalf("Expected value to be %v but got %v", "world", value)
		}
	})
}

func TestBool(t *testing.T) {
	t.Run("ok for true", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("ENV_VAR", "true")

		value := Bool("ENV_VAR", false)
		if !value {
			t.Fatalf("Expected value to be %v but got %v", true, value)
		}
	})

	t.Run("ok for false", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("ENV_VAR", "false")

		value := Bool("ENV_VAR", true)
		if value {
			t.Fatalf("Expected value to be %v but got %v", false, value)
		}
	})

	t.Run("ok with default", func(t *testing.T) {
	    os.Clearenv()
	    os.Setenv("ENV_VAR", "some")

		value := Bool("ENV_VAR", true)
		if value != true {
			t.Fatalf("Expected value to be %v but got %v", true, value)
		}
	})
}