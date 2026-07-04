package medicine

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
)

type ParsedMedicine struct {
	Name string
}

func ParseFile(filePath string) ([]ParsedMedicine, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".csv":
		return parseCSV(filePath)
	case ".pdf":
		return parsePDF(filePath)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}
}

func parseCSV(filePath string) ([]ParsedMedicine, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open csv file: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv file: %w", err)
	}

	var result []ParsedMedicine
	seen := make(map[string]bool)

	for _, record := range records {
		if len(record) == 0 {
			continue
		}
		name := strings.TrimSpace(record[0])
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		result = append(result, ParsedMedicine{Name: name})
	}

	return result, nil
}

func parsePDF(filePath string) ([]ParsedMedicine, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open pdf file: %w", err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat pdf file: %w", err)
	}

	r, err := pdf.NewReader(f, fi.Size())
	if err != nil {
		return nil, fmt.Errorf("failed to read pdf: %w", err)
	}

	var text strings.Builder
	for i := 1; i <= r.NumPage(); i++ {
		page := r.Page(i)

		fonts := make(map[string]*pdf.Font)
		for _, name := range page.Fonts() {
			f := page.Font(name)
			fonts[name] = &f
		}

		content, err := page.GetPlainText(fonts)
		if err != nil {
			continue
		}
		text.WriteString(content)
		text.WriteString("\n")
	}

	var result []ParsedMedicine
	seen := make(map[string]bool)

	for line := range strings.SplitSeq(text.String(), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || seen[line] {
			continue
		}
		seen[line] = true
		result = append(result, ParsedMedicine{Name: line})
	}

	return result, nil
}
