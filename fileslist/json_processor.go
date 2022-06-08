package fileslist

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

type jsonProcessor struct {
	Filename string
	files    *[]string
}

func NewJSONFilesListProcessor(filename string) *jsonProcessor {
	return &jsonProcessor{
		Filename: filename,
	}
}

func (j *jsonProcessor) GetFiles() ([]string, error) {
	if j.files != nil {
		return *j.files, nil
	}

	var files []string

	fd, err := os.Open(j.Filename)
	if err != nil {
		return files, errors.Wrap(err, "opening file")
	}
	defer fd.Close()

	if err := json.NewDecoder(fd).Decode(&files); err != nil {
		return files, errors.Wrap(err, "decoding file")
	}

	j.files = &files

	return *j.files, nil
}
