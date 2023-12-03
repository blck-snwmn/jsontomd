package jsontomd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/xerrors"
)

type pair struct {
	key   string
	value any
}

type jsonObject []pair

type jsonArray []jsonObject

func main() {
	var f string
	flag.StringVar(&f, "f", "", "json file path")
	flag.Parse()

	file, err := os.Open(f)
	if err != nil {
		panic(fmt.Sprintf("failed to open file: %+v", err))
	}
	decoder := json.NewDecoder(file)
	array, err := DecodeArray(decoder)
	if err != nil {
		panic(fmt.Sprintf("failed to decode: %+v", err))
	}
	md, err := EncodeMarkdown(array)
	if err != nil {
		panic(fmt.Sprintf("failed to encode: %+v", err))
	}
	fmt.Println(md)
}

func EncodeMarkdown(j jsonArray) (string, error) {
	header, err := encodeMarkdownHeader(j)
	if err != nil {
		return "", xerrors.Errorf("failed to encode header: %w", err)
	}
	d := encodeMarkdownDelimiter(j)
	body := encodeMarkdownBody(j)
	return header + d + body, nil
}

func encodeMarkdownHeader(j jsonArray) (string, error) {
	if len(j) == 0 {
		return "", xerrors.New("empty array")
	}
	var header []string
	for _, pair := range j[0] {
		header = append(header, pair.key)
	}
	return strings.Join(header, "|") + "\n", nil
}

func encodeMarkdownDelimiter(j jsonArray) string {
	var delimiter []string
	for range j[0] {
		delimiter = append(delimiter, "---")
	}
	return strings.Join(delimiter, "|") + "\n"
}

func encodeMarkdownBody(j jsonArray) string {
	var builder strings.Builder
	for _, obj := range j {
		for _, pair := range obj {
			builder.WriteString(fmt.Sprintf("%v|", pair.value))
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

func DecodeArray(decoder *json.Decoder) (jsonArray, error) {
	t, err := decoder.Token()
	if err != nil {
		return nil, xerrors.Errorf("failed to get token: %w", err)
	}
	switch td := t.(type) {
	case json.Delim:
		if td != '[' {
			return nil, xerrors.Errorf("expected [, got %v", td)
		}
	default:
		return nil, xerrors.Errorf("expected delimiter, got %T", t)
	}

	var result jsonArray
	for {
		item, next, err := decodeObject(decoder)
		if err != nil {
			return nil, xerrors.Errorf("failed to decode elem: %w", err)
		}
		if item != nil {
			result = append(result, item)
		}
		if !next {
			break
		}
	}

	return result, nil
}

func decodeObject(decoder *json.Decoder) (jsonObject, bool, error) {
	t, err := decoder.Token()
	if err != nil {
		return nil, false, xerrors.Errorf("failed to get token: %w", err)
	}
	switch td := t.(type) {
	case json.Delim:
		if td == ']' {
			return nil, false, nil
		}
		if td != '{' {
			return nil, false, xerrors.Errorf("expected [, got %v", td)
		}
	default:
		return nil, false, xerrors.Errorf("expected delimiter, got %T", t)
	}
	var result jsonObject
	for {
		key, err := decoder.Token()
		if err != nil {
			return nil, false, xerrors.Errorf("failed to get token: %w", err)
		}
		if td, ok := key.(json.Delim); ok && td == '}' {
			// end of object
			break
		}
		value, err := decoder.Token()
		if err != nil {
			return nil, false, xerrors.Errorf("failed to get token: %w", err)
		}
		result = append(result, pair{key.(string), value})
	}
	return result, true, nil
}
