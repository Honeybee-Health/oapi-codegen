package xml_marshal

import (
	"encoding/json"
	"encoding/xml"
	"testing"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ===== Simple Object Tests =====

func TestSimpleObject_XMLMarshalUnmarshal(t *testing.T) {
	obj := SimpleObject{
		Name:  "alice",
		Id:    42,
		Email: strPtr("alice@example.com"),
	}
	data, err := xml.Marshal(obj)
	require.NoError(t, err)

	var got SimpleObject
	err = xml.Unmarshal(data, &got)
	require.NoError(t, err)
	assert.Equal(t, obj.Name, got.Name)
	assert.Equal(t, obj.Id, got.Id)
	assert.Equal(t, *obj.Email, *got.Email)
}

func TestSimpleObject_XMLUnmarshalFromString(t *testing.T) {
	xmlStr := `<SimpleObject><name>bob</name><id>7</id></SimpleObject>`
	var got SimpleObject
	err := xml.Unmarshal([]byte(xmlStr), &got)
	require.NoError(t, err)
	assert.Equal(t, "bob", got.Name)
	assert.Equal(t, 7, got.Id)
	assert.Nil(t, got.Email)
}

// ===== XML Named Object Tests =====

func TestXmlNamedObject_MarshalUsesXmlNames(t *testing.T) {
	count := 5
	obj := XmlNamedObject{
		Title: "test",
		Count: &count,
	}
	data, err := xml.Marshal(obj)
	require.NoError(t, err)
	xmlStr := string(data)

	// Should use XML names from spec
	assert.Contains(t, xmlStr, "<obj_title>test</obj_title>")
	assert.Contains(t, xmlStr, "<obj_count>5</obj_count>")
}

func TestXmlNamedObject_UnmarshalUsesXmlNames(t *testing.T) {
	xmlStr := `<named_object><obj_title>hello</obj_title><obj_count>10</obj_count></named_object>`
	var got XmlNamedObject
	err := xml.Unmarshal([]byte(xmlStr), &got)
	require.NoError(t, err)
	assert.Equal(t, "hello", got.Title)
	require.NotNil(t, got.Count)
	assert.Equal(t, 10, *got.Count)
}

// ===== Additional Properties Tests (string type) =====

func TestAdditionalPropsString_XMLMarshalUnmarshal(t *testing.T) {
	obj := AdditionalPropsString{
		Name: "test",
		AdditionalProperties: map[string]string{
			"extra1": "val1",
			"extra2": "val2",
		},
	}
	data, err := xml.Marshal(obj)
	require.NoError(t, err)

	var got AdditionalPropsString
	err = xml.Unmarshal(data, &got)
	require.NoError(t, err)
	assert.Equal(t, "test", got.Name)
	assert.Equal(t, "val1", got.AdditionalProperties["extra1"])
	assert.Equal(t, "val2", got.AdditionalProperties["extra2"])
}

func TestAdditionalPropsString_XMLUnmarshal(t *testing.T) {
	xmlStr := `<AdditionalPropsString><name>foo</name><custom>bar</custom></AdditionalPropsString>`
	var got AdditionalPropsString
	err := xml.Unmarshal([]byte(xmlStr), &got)
	require.NoError(t, err)
	assert.Equal(t, "foo", got.Name)
	assert.Equal(t, "bar", got.AdditionalProperties["custom"])
}

func TestAdditionalPropsString_XMLMarshal(t *testing.T) {
	obj := AdditionalPropsString{
		Name: "hello",
		AdditionalProperties: map[string]string{
			"key1": "value1",
		},
	}
	data, err := xml.Marshal(obj)
	require.NoError(t, err)
	xmlStr := string(data)
	assert.Contains(t, xmlStr, "<name>hello</name>")
	assert.Contains(t, xmlStr, "<key1>value1</key1>")
}

// ===== Additional Properties Tests (int type) =====

func TestAdditionalPropsInt_XMLMarshalUnmarshal(t *testing.T) {
	label := "counter"
	obj := AdditionalPropsInt{
		Id:    1,
		Label: &label,
		AdditionalProperties: map[string]int{
			"score": 100,
		},
	}
	data, err := xml.Marshal(obj)
	require.NoError(t, err)

	var got AdditionalPropsInt
	err = xml.Unmarshal(data, &got)
	require.NoError(t, err)
	assert.Equal(t, 1, got.Id)
	require.NotNil(t, got.Label)
	assert.Equal(t, "counter", *got.Label)
	assert.Equal(t, 100, got.AdditionalProperties["score"])
}

// ===== Simple OneOf Union Tests =====

func TestSimpleUnion_XMLFromVariantA(t *testing.T) {
	var u SimpleUnion
	valA := 42
	err := u.FromOneOfVariantA(OneOfVariantA{
		VariantA: "hello",
		ValueA:   &valA,
	})
	require.NoError(t, err)

	// Marshal to XML
	data, err := xml.Marshal(u)
	require.NoError(t, err)
	assert.Contains(t, string(data), "hello")

	// Unmarshal back
	var got SimpleUnion
	err = xml.Unmarshal(data, &got)
	require.NoError(t, err)

	// Extract variant
	va, err := got.AsOneOfVariantA()
	require.NoError(t, err)
	assert.Equal(t, "hello", va.VariantA)
	require.NotNil(t, va.ValueA)
	assert.Equal(t, 42, *va.ValueA)
}

func TestSimpleUnion_XMLFromVariantB(t *testing.T) {
	var u SimpleUnion
	valB := true
	err := u.FromOneOfVariantB(OneOfVariantB{
		VariantB: "world",
		ValueB:   &valB,
	})
	require.NoError(t, err)

	// Marshal to XML
	data, err := xml.Marshal(u)
	require.NoError(t, err)
	assert.Contains(t, string(data), "world")

	// Unmarshal back
	var got SimpleUnion
	err = xml.Unmarshal(data, &got)
	require.NoError(t, err)

	vb, err := got.AsOneOfVariantB()
	require.NoError(t, err)
	assert.Equal(t, "world", vb.VariantB)
	require.NotNil(t, vb.ValueB)
	assert.True(t, *vb.ValueB)
}

// ===== OneOf with Discriminator Tests =====

func TestUnionWithDiscriminator_XMLFromVariant1(t *testing.T) {
	var u UnionWithDiscriminator
	err := u.FromDiscriminatorVariant1(DiscriminatorVariant1{
		Name: "test-name",
	})
	require.NoError(t, err)

	// Marshal to XML
	data, err := xml.Marshal(u)
	require.NoError(t, err)

	// Unmarshal back
	var got UnionWithDiscriminator
	err = xml.Unmarshal(data, &got)
	require.NoError(t, err)

	disc, err := got.Discriminator()
	require.NoError(t, err)
	assert.Equal(t, "variant1", disc)

	v1, err := got.AsDiscriminatorVariant1()
	require.NoError(t, err)
	assert.Equal(t, "test-name", v1.Name)
}

func TestUnionWithDiscriminator_XMLFromVariant2(t *testing.T) {
	var u UnionWithDiscriminator
	err := u.FromDiscriminatorVariant2(DiscriminatorVariant2{
		Count: 99,
	})
	require.NoError(t, err)

	data, err := xml.Marshal(u)
	require.NoError(t, err)

	var got UnionWithDiscriminator
	err = xml.Unmarshal(data, &got)
	require.NoError(t, err)

	disc, err := got.Discriminator()
	require.NoError(t, err)
	assert.Equal(t, "variant2", disc)

	v2, err := got.AsDiscriminatorVariant2()
	require.NoError(t, err)
	assert.Equal(t, 99, v2.Count)
}

func TestUnionWithDiscriminator_ValueByDiscriminator(t *testing.T) {
	var u UnionWithDiscriminator
	err := u.FromDiscriminatorVariant1(DiscriminatorVariant1{
		Name: "disc-test",
	})
	require.NoError(t, err)

	data, err := xml.Marshal(u)
	require.NoError(t, err)

	var got UnionWithDiscriminator
	err = xml.Unmarshal(data, &got)
	require.NoError(t, err)

	val, err := got.ValueByDiscriminator()
	require.NoError(t, err)
	v1, ok := val.(DiscriminatorVariant1)
	assert.True(t, ok)
	assert.Equal(t, "disc-test", v1.Name)
}

// ===== OneOf with Fixed Properties Tests =====

func TestUnionWithFixedProps_XMLRoundTrip(t *testing.T) {
	var u UnionWithFixedProps
	err := u.FromOneOfVariantA(OneOfVariantA{
		VariantA: "fixed-test",
	})
	require.NoError(t, err)
	assert.Equal(t, "a", u.Type) // discriminator should be set

	data, err := xml.Marshal(u)
	require.NoError(t, err)

	var got UnionWithFixedProps
	err = xml.Unmarshal(data, &got)
	require.NoError(t, err)
	assert.Equal(t, "a", got.Type)

	va, err := got.AsOneOfVariantA()
	require.NoError(t, err)
	assert.Equal(t, "fixed-test", va.VariantA)
}

// ===== AnyOf Union Tests =====

func TestAnyOfUnion_XMLRoundTrip(t *testing.T) {
	var u AnyOfUnion
	valA := 10
	err := u.FromOneOfVariantA(OneOfVariantA{
		VariantA: "any-test",
		ValueA:   &valA,
	})
	require.NoError(t, err)

	data, err := xml.Marshal(u)
	require.NoError(t, err)

	var got AnyOfUnion
	err = xml.Unmarshal(data, &got)
	require.NoError(t, err)

	va, err := got.AsOneOfVariantA()
	require.NoError(t, err)
	assert.Equal(t, "any-test", va.VariantA)
	require.NotNil(t, va.ValueA)
	assert.Equal(t, 10, *va.ValueA)
}

// ===== AllOf Merged Object Tests =====

func TestExtendedObject_XMLMarshalUnmarshal(t *testing.T) {
	desc := "a description"
	obj := ExtendedObject{
		Id:          1,
		Name:        "extended",
		Description: &desc,
	}
	data, err := xml.Marshal(obj)
	require.NoError(t, err)

	var got ExtendedObject
	err = xml.Unmarshal(data, &got)
	require.NoError(t, err)
	assert.Equal(t, 1, got.Id)
	assert.Equal(t, "extended", got.Name)
	require.NotNil(t, got.Description)
	assert.Equal(t, "a description", *got.Description)
}

func TestExtendedObject_XMLUnmarshalFromString(t *testing.T) {
	xmlStr := `<ExtendedObject><id>5</id><name>test</name><description>desc</description></ExtendedObject>`
	var got ExtendedObject
	err := xml.Unmarshal([]byte(xmlStr), &got)
	require.NoError(t, err)
	assert.Equal(t, 5, got.Id)
	assert.Equal(t, "test", got.Name)
	require.NotNil(t, got.Description)
	assert.Equal(t, "desc", *got.Description)
}

// ===== Union with Additional Properties Tests =====

func TestUnionWithAdditionalProps_XMLRoundTrip(t *testing.T) {
	var u UnionWithAdditionalProps
	u.AdditionalProperties = map[string]string{"extra": "data"}
	err := u.MergeOneOfVariantA(OneOfVariantA{VariantA: "ap-test"})
	require.NoError(t, err)
	assert.Equal(t, "a", u.Type) // discriminator should be set

	// JSON should still work
	jdata, err := json.Marshal(u)
	require.NoError(t, err)
	assert.Contains(t, string(jdata), `"extra":"data"`)
	assert.Contains(t, string(jdata), `"variantA":"ap-test"`)
}

func TestUnionWithAdditionalProps_JSONRoundTrip(t *testing.T) {
	var u UnionWithAdditionalProps
	u.AdditionalProperties = map[string]string{"key": "val"}
	err := u.MergeOneOfVariantB(OneOfVariantB{VariantB: "json-test"})
	require.NoError(t, err)

	data, err := json.Marshal(u)
	require.NoError(t, err)

	var got UnionWithAdditionalProps
	err = json.Unmarshal(data, &got)
	require.NoError(t, err)
	assert.Equal(t, "b", got.Type)
	val, ok := got.Get("key")
	assert.True(t, ok)
	assert.Equal(t, "val", val)
}

// ===== Object with Optional Fields =====

func TestObjectWithOptionalFields_XMLMarshalUnmarshal(t *testing.T) {
	optStr := "optional"
	obj := ObjectWithOptionalFields{
		RequiredField:  "required",
		OptionalString: &optStr,
	}
	data, err := xml.Marshal(obj)
	require.NoError(t, err)

	var got ObjectWithOptionalFields
	err = xml.Unmarshal(data, &got)
	require.NoError(t, err)
	assert.Equal(t, "required", got.RequiredField)
	require.NotNil(t, got.OptionalString)
	assert.Equal(t, "optional", *got.OptionalString)
	assert.Nil(t, got.OptionalInt)
}

// ===== Nested Object Tests =====

func TestNestedObject_XMLMarshalUnmarshal(t *testing.T) {
	obj := OuterObject{
		Label: "outer",
		Inner: InnerObject{
			Value: "inner-value",
		},
	}
	data, err := xml.Marshal(obj)
	require.NoError(t, err)

	var got OuterObject
	err = xml.Unmarshal(data, &got)
	require.NoError(t, err)
	assert.Equal(t, "outer", got.Label)
	assert.Equal(t, "inner-value", got.Inner.Value)
}

// ===== JSON-XML Interoperability Tests =====

func TestSimpleUnion_JSONAndXMLInterop(t *testing.T) {
	// Set via From (which populates both union and xunion)
	var u SimpleUnion
	err := u.FromOneOfVariantA(OneOfVariantA{VariantA: "interop"})
	require.NoError(t, err)

	// Verify JSON marshal works
	jdata, err := json.Marshal(u)
	require.NoError(t, err)
	assert.Contains(t, string(jdata), "interop")

	// Verify XML marshal works
	xdata, err := xml.Marshal(u)
	require.NoError(t, err)
	assert.Contains(t, string(xdata), "interop")

	// Unmarshal from JSON
	var fromJson SimpleUnion
	err = json.Unmarshal(jdata, &fromJson)
	require.NoError(t, err)
	va, err := fromJson.AsOneOfVariantA()
	require.NoError(t, err)
	assert.Equal(t, "interop", va.VariantA)

	// Unmarshal from XML
	var fromXml SimpleUnion
	err = xml.Unmarshal(xdata, &fromXml)
	require.NoError(t, err)
	va2, err := fromXml.AsOneOfVariantA()
	require.NoError(t, err)
	assert.Equal(t, "interop", va2.VariantA)
}

func TestAdditionalPropsString_JSONAndXMLInterop(t *testing.T) {
	obj := AdditionalPropsString{
		Name: "interop",
		AdditionalProperties: map[string]string{
			"dynamic": "field",
		},
	}

	// JSON round trip
	jdata, err := json.Marshal(obj)
	require.NoError(t, err)
	var fromJson AdditionalPropsString
	err = json.Unmarshal(jdata, &fromJson)
	require.NoError(t, err)
	assert.Equal(t, "interop", fromJson.Name)
	assert.Equal(t, "field", fromJson.AdditionalProperties["dynamic"])

	// XML round trip
	xdata, err := xml.Marshal(obj)
	require.NoError(t, err)
	var fromXml AdditionalPropsString
	err = xml.Unmarshal(xdata, &fromXml)
	require.NoError(t, err)
	assert.Equal(t, "interop", fromXml.Name)
	assert.Equal(t, "field", fromXml.AdditionalProperties["dynamic"])
}

// ===== Edge Case Tests =====

func TestSimpleUnion_EmptyUnionMarshalXML(t *testing.T) {
	var u SimpleUnion
	// Empty union should marshal without errors
	data, err := xml.Marshal(u)
	require.NoError(t, err)
	assert.NotEmpty(t, data)
}

func TestAdditionalPropsString_NoAdditionalProps(t *testing.T) {
	obj := AdditionalPropsString{Name: "only-name"}
	data, err := xml.Marshal(obj)
	require.NoError(t, err)

	var got AdditionalPropsString
	err = xml.Unmarshal(data, &got)
	require.NoError(t, err)
	assert.Equal(t, "only-name", got.Name)
	assert.Nil(t, got.AdditionalProperties)
}

func TestDiscriminator_XMLEmptyUnion(t *testing.T) {
	var u UnionWithDiscriminator
	disc, err := u.Discriminator()
	require.NoError(t, err)
	assert.Equal(t, "", disc)
}

// ===== Helper functions =====

func strPtr(s string) *string {
	return &s
}

// ===== Date Tests =====
//
// `format: date` must serialize as CCYY-MM-DD in XML. The type it used to map
// to (openapi_types.Date) defines UnmarshalText but no MarshalText, so XML
// marshaling fell through to the embedded time.Time and emitted RFC3339 --
// invalid for an xsd:date element. See PZR-254.

func mustDate(t *testing.T, s string) XMLDate {
	t.Helper()
	parsed, err := time.Parse("2006-01-02", s)
	require.NoError(t, err)
	return XMLDate{Time: parsed}
}

func TestDateWrapper_XMLMarshalsAsPlainDate(t *testing.T) {
	obj := DateWrapper{Date: mustDate(t, "2010-12-01")}

	data, err := xml.Marshal(obj)
	require.NoError(t, err)
	assert.Equal(t, "<DateWrapper><Date>2010-12-01</Date></DateWrapper>", string(data))
	assert.NotContains(t, string(data), "T00:00:00Z", "date must not render as an RFC3339 timestamp")
}

func TestDateWrapper_XMLRoundTrip(t *testing.T) {
	input := "<DateWrapper><Date>2010-12-01</Date></DateWrapper>"

	var obj DateWrapper
	require.NoError(t, xml.Unmarshal([]byte(input), &obj))
	assert.Equal(t, "2010-12-01", obj.Date.String())

	out, err := xml.Marshal(obj)
	require.NoError(t, err)
	assert.Equal(t, input, string(out), "re-marshaling must reproduce the identical plain-date form")
}

func TestDateTimeWrapper_StaysRFC3339(t *testing.T) {
	obj := DateTimeWrapper{DateTime: time.Date(2010, 12, 1, 13, 45, 30, 0, time.UTC)}

	data, err := xml.Marshal(obj)
	require.NoError(t, err)
	assert.Equal(t, "<DateTimeWrapper><DateTime>2010-12-01T13:45:30Z</DateTime></DateTimeWrapper>", string(data))
}

func TestDateHolder_AllPositionsMarshalAsPlainDate(t *testing.T) {
	attr := mustDate(t, "2003-03-03")
	opt := mustDate(t, "2002-02-02")
	list := []XMLDate{mustDate(t, "2004-04-04"), mustDate(t, "2005-05-05")}
	obj := DateHolder{
		RequiredDate: mustDate(t, "2001-01-01"),
		OptionalDate: &opt,
		DateList:     &list,
		AttrDate:     &attr,
		Stamp:        time.Date(2006, 6, 6, 7, 8, 9, 0, time.UTC),
	}

	data, err := xml.Marshal(obj)
	require.NoError(t, err)
	out := string(data)

	// Attributes are the case a MarshalXML method could not have fixed:
	// encoding/xml ignores xml.Marshaler for attributes and uses TextMarshaler.
	assert.Contains(t, out, `attr_date="2003-03-03"`)
	assert.Contains(t, out, "<required_date>2001-01-01</required_date>")
	assert.Contains(t, out, "<optional_date>2002-02-02</optional_date>")
	assert.Contains(t, out, "<date_list>2004-04-04</date_list>")
	assert.Contains(t, out, "<date_list>2005-05-05</date_list>")
	assert.Contains(t, out, "<stamp>2006-06-06T07:08:09Z</stamp>", "date-time must remain RFC3339")
	assert.NotContains(t, out, "2001-01-01T", "dates must not carry a time component")
}

func TestDateHolder_XMLRoundTrip(t *testing.T) {
	attr := mustDate(t, "2003-03-03")
	opt := mustDate(t, "2002-02-02")
	list := []XMLDate{mustDate(t, "2004-04-04")}
	obj := DateHolder{
		RequiredDate: mustDate(t, "2001-01-01"),
		OptionalDate: &opt,
		DateList:     &list,
		AttrDate:     &attr,
		Stamp:        time.Date(2006, 6, 6, 7, 8, 9, 0, time.UTC),
	}

	data, err := xml.Marshal(obj)
	require.NoError(t, err)

	var got DateHolder
	require.NoError(t, xml.Unmarshal(data, &got))

	assert.Equal(t, "2001-01-01", got.RequiredDate.String())
	require.NotNil(t, got.OptionalDate)
	assert.Equal(t, "2002-02-02", got.OptionalDate.String())
	require.NotNil(t, got.AttrDate)
	assert.Equal(t, "2003-03-03", got.AttrDate.String())
	require.NotNil(t, got.DateList)
	require.Len(t, *got.DateList, 1)
	assert.Equal(t, "2004-04-04", (*got.DateList)[0].String())
	assert.True(t, obj.Stamp.Equal(got.Stamp))
}

func TestDate_JSONUnchanged(t *testing.T) {
	obj := DateWrapper{Date: mustDate(t, "2010-12-01")}

	data, err := json.Marshal(obj)
	require.NoError(t, err)
	assert.JSONEq(t, `{"Date":"2010-12-01"}`, string(data))

	var got DateWrapper
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Equal(t, "2010-12-01", got.Date.String())
}

func TestDate_AliasedSchema(t *testing.T) {
	// BirthDate is a bare `format: date` schema, generated via alias.
	var bd BirthDate = mustDate(t, "1985-07-04")

	data, err := xml.Marshal(struct {
		XMLName xml.Name  `xml:"b"`
		Value   BirthDate `xml:"v"`
	}{Value: bd})
	require.NoError(t, err)
	assert.Equal(t, "<b><v>1985-07-04</v></b>", string(data))
}

func TestXMLDate_ConvertibleToRuntimeDate(t *testing.T) {
	// The runtime's parameter binding detects dates by reflect-converting to
	// types.Date. XMLDate must stay convertible or deepObject binding silently
	// yields the zero date.
	d := mustDate(t, "2010-12-01")
	converted := openapi_types.Date(d)
	assert.Equal(t, "2010-12-01", converted.String())
	assert.Equal(t, d.Time, converted.Time)
}
