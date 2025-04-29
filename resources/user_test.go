package resources

import (
	"testing"
)

func TestResourceAidboxUser(t *testing.T) {
	resource := ResourceAidboxUser()
	if resource == nil {
		t.Fatal("resource is nil")
	}

	// Test schema
	schema := resource.Schema
	if schema == nil {
		t.Fatal("schema is nil")
	}

	// Test required fields
	requiredFields := []string{"name"}
	for _, field := range requiredFields {
		if schema[field] == nil {
			t.Errorf("required field %s is missing", field)
		}
		if !schema[field].Required {
			t.Errorf("field %s should be required", field)
		}
	}

	// Test optional fields
	optionalFields := []string{"password", "resource_id", "extensions"}
	for _, field := range optionalFields {
		if schema[field] == nil {
			t.Errorf("optional field %s is missing", field)
		}
		if schema[field].Required {
			t.Errorf("field %s should not be required", field)
		}
	}

	// Test computed fields
	computedFields := []string{"id", "resource_type", "meta"}
	for _, field := range computedFields {
		if schema[field] == nil {
			t.Errorf("computed field %s is missing", field)
		}
		if !schema[field].Computed {
			t.Errorf("field %s should be computed", field)
		}
	}
}
