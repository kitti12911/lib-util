package fieldmask

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestValidateMaskAllowsNestedLeafFields(t *testing.T) {
	err := ValidateMask(&fieldmaskpb.FieldMask{Paths: []string{
		"email",
		"profile.first_name",
		"profile.address.country_code",
	}}, testUserMessage(t), map[string]bool{
		"id":         true,
		"created_at": true,
		"updated_at": true,
	})
	if err != nil {
		t.Fatalf("ValidateMask() error = %v", err)
	}
}

func TestValidateMaskRejectsMessageFieldWithLeafPathMessage(t *testing.T) {
	err := ValidateMask(&fieldmaskpb.FieldMask{Paths: []string{"profile"}}, testUserMessage(t), nil)
	if err == nil {
		t.Fatal("ValidateMask() error = nil")
	}
	if !strings.Contains(err.Error(), "specify a leaf field path") {
		t.Fatalf("ValidateMask() error = %v, want leaf path message", err)
	}
}

func TestValidateMaskRejectsUnsupportedPaths(t *testing.T) {
	tests := []string{
		"id",
		"created_at",
		"updated_at",
		"profile.address.unknown",
	}

	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			err := ValidateMask(&fieldmaskpb.FieldMask{Paths: []string{path}}, testUserMessage(t), map[string]bool{
				"id":         true,
				"created_at": true,
				"updated_at": true,
			})
			if err == nil {
				t.Fatal("ValidateMask() error = nil")
			}
		})
	}
}

func TestExtractChangesSupportsNestedFields(t *testing.T) {
	user := testUserMessage(t)
	profile := dynamicpb.NewMessage(field(user.Descriptor(), "profile").Message())
	address := dynamicpb.NewMessage(field(profile.Descriptor(), "address").Message())

	setString(user, "display_name", "Display Name")
	setString(profile, "first_name", "First")
	setString(address, "line2", "Floor 2")
	setString(address, "country_code", "TH")
	profile.Set(field(profile.Descriptor(), "address"), protoreflect.ValueOfMessage(address))
	user.Set(field(user.Descriptor(), "profile"), protoreflect.ValueOfMessage(profile))

	changes := ExtractChanges(&fieldmaskpb.FieldMask{Paths: []string{
		"display_name",
		"profile.first_name",
		"profile.address.line2",
		"profile.address.country_code",
	}}, user)

	assertChange(t, changes, "display_name", "Display Name")
	assertChange(t, changes, "profile.first_name", "First")
	assertChange(t, changes, "profile.address.line2", "Floor 2")
	assertChange(t, changes, "profile.address.country_code", "TH")
}

func TestExtractChangesCanClearOptionalLeaf(t *testing.T) {
	user := testUserMessage(t)
	profile := dynamicpb.NewMessage(field(user.Descriptor(), "profile").Message())
	address := dynamicpb.NewMessage(field(profile.Descriptor(), "address").Message())
	profile.Set(field(profile.Descriptor(), "address"), protoreflect.ValueOfMessage(address))
	user.Set(field(user.Descriptor(), "profile"), protoreflect.ValueOfMessage(profile))

	changes := ExtractChanges(&fieldmaskpb.FieldMask{Paths: []string{
		"profile.address.line2",
	}}, user)

	if got, ok := changes["profile.address.line2"]; !ok || got != nil {
		t.Fatalf("changes[profile.address.line2] = %#v, %t; want nil, true", got, ok)
	}
}

func TestExtractChangesWithNilInputsReturnsEmptyChanges(t *testing.T) {
	if got := ExtractChanges(nil, testUserMessage(t)); len(got) != 0 {
		t.Fatalf("ExtractChanges(nil, msg) len = %d, want 0", len(got))
	}
	if got := ExtractChanges(&fieldmaskpb.FieldMask{Paths: []string{"email"}}, nil); len(got) != 0 {
		t.Fatalf("ExtractChanges(mask, nil) len = %d, want 0", len(got))
	}
}

func TestExtractNestedChangesReturnsRootFields(t *testing.T) {
	changes := map[string]any{
		"email":                        "a@example.com",
		"username":                     "alice",
		"profile.first_name":           "Alice",
		"profile.address.country_code": "TH",
	}

	got := ExtractNestedChanges(changes, map[string]string{
		"email":    "email",
		"username": "username",
	}, RootNestedName)

	assertChange(t, got, "email", "a@example.com")
	assertChange(t, got, "username", "alice")
	if len(got) != 2 {
		t.Fatalf("ExtractNestedChanges() len = %d, want 2", len(got))
	}
}

func TestExtractNestedChangesReturnsDirectNestedFields(t *testing.T) {
	changes := map[string]any{
		"email":                        "a@example.com",
		"profile.first_name":           "Alice",
		"profile.last_name":            "Example",
		"profile.address.country_code": "TH",
	}

	got := ExtractNestedChanges(changes, map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
	}, "profile")

	assertChange(t, got, "first_name", "Alice")
	assertChange(t, got, "last_name", "Example")
	if len(got) != 2 {
		t.Fatalf("ExtractNestedChanges() len = %d, want 2", len(got))
	}
}

func TestExtractNestedChangesReturnsDeepNestedFields(t *testing.T) {
	changes := map[string]any{
		"profile.first_name":           "Alice",
		"profile.address.city":         "Bangkok",
		"profile.address.country_code": "TH",
	}

	got := ExtractNestedChanges(changes, map[string]string{
		"city":         "city",
		"country_code": "country_code",
	}, "profile.address")

	assertChange(t, got, "city", "Bangkok")
	assertChange(t, got, "country_code", "TH")
	if len(got) != 2 {
		t.Fatalf("ExtractNestedChanges() len = %d, want 2", len(got))
	}
}

func TestExtractNestedChangesCanRenameFields(t *testing.T) {
	changes := map[string]any{
		"profile.phone_number": "+66123456789",
	}

	got := ExtractNestedChanges(changes, map[string]string{
		"phone_number": "phone",
	}, "profile")

	assertChange(t, got, "phone", "+66123456789")
}

func testUserMessage(t *testing.T) *dynamicpb.Message {
	t.Helper()

	file := &descriptorpb.FileDescriptorProto{
		Syntax:  stringPtr("proto2"),
		Name:    stringPtr("test.proto"),
		Package: stringPtr("test.v1"),
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: stringPtr("User"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField("id", 1),
					stringField("email", 2),
					stringField("username", 3),
					stringField("display_name", 4),
					stringField("created_at", 5),
					stringField("updated_at", 6),
					messageField("profile", 7, ".test.v1.Profile"),
				},
			},
			{
				Name: stringPtr("Profile"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField("first_name", 1),
					stringField("last_name", 2),
					stringField("phone_number", 3),
					messageField("address", 4, ".test.v1.Address"),
				},
			},
			{
				Name: stringPtr("Address"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField("line2", 1),
					stringField("city", 2),
					stringField("country_code", 3),
				},
			},
		},
	}

	desc, err := protodesc.NewFile(file, nil)
	if err != nil {
		t.Fatalf("protodesc.NewFile() error = %v", err)
	}

	return dynamicpb.NewMessage(desc.Messages().ByName("User"))
}

func stringField(name string, number int32) *descriptorpb.FieldDescriptorProto {
	return &descriptorpb.FieldDescriptorProto{
		Name:   stringPtr(name),
		Number: int32Ptr(number),
		Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
		Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
	}
}

func messageField(name string, number int32, typeName string) *descriptorpb.FieldDescriptorProto {
	return &descriptorpb.FieldDescriptorProto{
		Name:     stringPtr(name),
		Number:   int32Ptr(number),
		Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
		Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
		TypeName: stringPtr(typeName),
	}
}

func setString(msg *dynamicpb.Message, name string, value string) {
	msg.Set(field(msg.Descriptor(), name), protoreflect.ValueOfString(value))
}

func field(desc protoreflect.MessageDescriptor, name string) protoreflect.FieldDescriptor {
	return desc.Fields().ByName(protoreflect.Name(name))
}

func stringPtr(value string) *string {
	return &value
}

func int32Ptr(value int32) *int32 {
	return &value
}

func assertChange(t *testing.T, changes map[string]any, key string, want any) {
	t.Helper()

	got, ok := changes[key]
	if !ok {
		t.Fatalf("changes[%s] is missing", key)
	}
	if got != want {
		t.Fatalf("changes[%s] = %#v, want %#v", key, got, want)
	}
}
