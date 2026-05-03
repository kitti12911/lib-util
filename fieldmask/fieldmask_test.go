package fieldmask

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)
}

func TestValidateMaskRejectsMessageFieldWithLeafPathMessage(t *testing.T) {
	err := ValidateMask(&fieldmaskpb.FieldMask{Paths: []string{"profile"}}, testUserMessage(t), nil)
	require.Error(t, err)
	assert.ErrorContains(t, err, "specify a leaf field path")
}

func TestValidateMaskRejectsMissingMask(t *testing.T) {
	tests := []struct {
		name string
		mask *fieldmaskpb.FieldMask
	}{
		{
			name: "nil mask",
			mask: nil,
		},
		{
			name: "empty paths",
			mask: &fieldmaskpb.FieldMask{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMask(tt.mask, testUserMessage(t), nil)

			require.Error(t, err)
			assert.ErrorContains(t, err, "update_mask is required")
		})
	}
}

func TestValidateMaskRejectsNilMessage(t *testing.T) {
	err := ValidateMask(&fieldmaskpb.FieldMask{Paths: []string{"email"}}, nil, nil)

	require.Error(t, err)
	assert.ErrorContains(t, err, "message is required")
}

func TestValidateMaskRejectsUnsupportedPaths(t *testing.T) {
	tests := []string{
		"",
		"id",
		"created_at",
		"updated_at",
		"profile.address.unknown",
		"profile.unknown.value",
		"email.value",
	}

	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			err := ValidateMask(&fieldmaskpb.FieldMask{Paths: []string{path}}, testUserMessage(t), map[string]bool{
				"id":         true,
				"created_at": true,
				"updated_at": true,
			})
			require.Error(t, err)
		})
	}
}

func TestExtractChangesReturnsNilForEmptyPath(t *testing.T) {
	changes := ExtractChanges(&fieldmaskpb.FieldMask{Paths: []string{""}}, testUserMessage(t))

	assertChange(t, changes, "", nil)
}

func TestExtractChangesReturnsNilForInvalidPaths(t *testing.T) {
	user := testUserMessage(t)
	profile := dynamicpb.NewMessage(field(user.Descriptor(), "profile").Message())
	user.Set(field(user.Descriptor(), "profile"), protoreflect.ValueOfMessage(profile))

	changes := ExtractChanges(&fieldmaskpb.FieldMask{Paths: []string{
		"unknown",
		"profile.unknown.value",
		"email.value",
		"profile.first_name",
	}}, user)

	assertChange(t, changes, "unknown", nil)
	assertChange(t, changes, "profile.unknown.value", nil)
	assertChange(t, changes, "email.value", nil)
	assertChange(t, changes, "profile.first_name", nil)
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

	got, ok := changes["profile.address.line2"]
	require.True(t, ok)
	assert.Nil(t, got)
}

func TestExtractChangesWithNilInputsReturnsEmptyChanges(t *testing.T) {
	assert.Empty(t, ExtractChanges(nil, testUserMessage(t)))
	assert.Empty(t, ExtractChanges(&fieldmaskpb.FieldMask{Paths: []string{"email"}}, nil))
}

func TestExtractNestedChangesReturnsRootFields(t *testing.T) {
	changes := map[string]any{
		"":                             "ignored",
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
	assert.Len(t, got, 2)
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
	assert.Len(t, got, 2)
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
	assert.Len(t, got, 2)
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

func TestExtractNestedChangesIgnoresUnmappedFields(t *testing.T) {
	changes := map[string]any{
		"profile.first_name": "Alice",
	}

	got := ExtractNestedChanges(changes, map[string]string{
		"last_name": "last_name",
	}, "profile")

	assert.Empty(t, got)
}

func testUserMessage(t *testing.T) *dynamicpb.Message {
	t.Helper()

	file := &descriptorpb.FileDescriptorProto{
		Syntax:  new("proto2"),
		Name:    new("test.proto"),
		Package: new("test.v1"),
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: new("User"),
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
				Name: new("Profile"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField("first_name", 1),
					stringField("last_name", 2),
					stringField("phone_number", 3),
					messageField("address", 4, ".test.v1.Address"),
				},
			},
			{
				Name: new("Address"),
				Field: []*descriptorpb.FieldDescriptorProto{
					stringField("line2", 1),
					stringField("city", 2),
					stringField("country_code", 3),
				},
			},
		},
	}

	desc, err := protodesc.NewFile(file, nil)
	require.NoError(t, err)

	return dynamicpb.NewMessage(desc.Messages().ByName("User"))
}

func stringField(name string, number int32) *descriptorpb.FieldDescriptorProto {
	return &descriptorpb.FieldDescriptorProto{
		Name:   new(name),
		Number: new(number),
		Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
		Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
	}
}

func messageField(name string, number int32, typeName string) *descriptorpb.FieldDescriptorProto {
	return &descriptorpb.FieldDescriptorProto{
		Name:     new(name),
		Number:   new(number),
		Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
		Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
		TypeName: new(typeName),
	}
}

func setString(msg *dynamicpb.Message, name string, value string) {
	msg.Set(field(msg.Descriptor(), name), protoreflect.ValueOfString(value))
}

func field(desc protoreflect.MessageDescriptor, name string) protoreflect.FieldDescriptor {
	return desc.Fields().ByName(protoreflect.Name(name))
}

func assertChange(t *testing.T, changes map[string]any, key string, want any) {
	t.Helper()

	got, ok := changes[key]
	require.True(t, ok, "changes[%s] is missing", key)
	assert.Equal(t, want, got)
}
