// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying files LICENSE.txt

package strutil

import (
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HelloKittyWorld", "hello_kitty_world"},
		{"testCase", "test_case"},
		{"TestCase", "test_case"},
		{"Test Case", "test_case"},
		{" Test Case", "test_case"},
		{"Test Case ", "test_case"},
		{" Test Case ", "test_case"},
		{"test", "test"},
		{"test_case", "test_case"},
		{"Test", "test"},
		{"", ""},
		{"ManyManyWords", "many_many_words"},
		{"manyManyWords", "many_many_words"},
		{"AnyKind of_string", "any_kind_of_string"},
		{"numbers2and55with000", "numbers_2_and_55_with_000"},
		{"JSONData", "json_data"},
		{"userID", "user_id"},
		{"AAAbbb", "aa_abbb"},
		{"", ""},
		{" ", ""},
	}

	for i, test := range tests {
		output := ToSnakeCase(test.input)
		if test.expected != output {
			t.Fatalf("Test case %d is not successful, expect: %s, got: %s", i, test.expected, output)
		}
	}
}

func TestToScreamingSnake(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"helloWorld", "HELLO_WORLD"},
		{"hello2world", "HELLO_2_WORLD"},
		{"hello world_again", "HELLO_WORLD_AGAIN"},
		{"HELLO_WORLD", "HELLO_WORLD"},
		{"", ""},
	}

	for i, test := range tests {
		output := ToScreamingSnake(test.input)
		if test.expected != output {
			t.Fatalf("Test case %d is not successful, expect: %s, got: %s", i, test.expected, output)
		}
	}
}

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HelloKittyWorld", "hello-kitty-world"},
		{"hello", "hello"},
		{"Hello2World", "hello-2-world"},
		{"Hello_World", "hello-world"},
		{"Hello World", "hello-world"},
		{"", ""},
	}

	for i, test := range tests {
		output := ToKebabCase(test.input)
		if test.expected != output {
			t.Fatalf("Test case %d is not successful, expect: %s, got: %s", i, test.expected, output)
		}
	}
}

func TestToScreamingKebab(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"helloWorld", "HELLO-WORLD"},
		{"hello2world", "HELLO-2-WORLD"},
		{"hello world_again", "HELLO-WORLD-AGAIN"},
		{"HELLO-WORLD", "HELLO-WORLD"},
		{"Hello", "HELLO"},
		{"", ""},
	}

	for i, test := range tests {
		output := ToScreamingKebab(test.input)
		if test.expected != output {
			t.Fatalf("Test case %d is not successful, expect: %s, got: %s", i, test.expected, output)
		}
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"test_case", "TestCase"},
		{"test", "Test"},
		{"TestCase", "TestCase"},
		{" test  case ", "TestCase"},
		{"", ""},
		{"many_many_words", "ManyManyWords"},
		{"AnyKind of_string", "AnyKindOfString"},
		{"odd-fix", "OddFix"},
		{"numbers2And55with000", "Numbers2And55With000"},
		{"hello_kitty_world", "HelloKittyWorld"},
		{"", ""},
		{" ", ""},
	}

	for i, test := range tests {
		output := ToCamelCase(test.input)
		if test.expected != output {
			t.Fatalf("Test case %d is not successful, expect: %s, got: %s", i, test.expected, output)
		}
	}
}

func TestToLowerCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello", "hello"},
		{"hello", "hello"},
		{"helloWorld", "helloWorld"},
		{"hello2World", "hello2World"},
		{"hello_world", "helloWorld"},
		{"hello world", "helloWorld"},
		{" ", ""},
	}

	for i, test := range tests {
		output := ToLowerCamelCase(test.input)
		if test.expected != output {
			t.Fatalf("Test case %d is not successful, expect: %s, got: %s", i, test.expected, output)
		}
	}
}
