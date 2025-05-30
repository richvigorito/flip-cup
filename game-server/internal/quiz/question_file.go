package quiz

import (
    "fmt"
    "gopkg.in/yaml.v2"
	  "io/ioutil"
	  "path/filepath"
	  "math/rand"
)

type QuestionFile struct {
    ID        string        `yaml:"id"`
    Name      string        `yaml:"name"`
    Filename  string        `yaml:"filename"`
    Questions []*Question   `yaml:"questions"`
}

func NewQuestionFile(filename string) (*QuestionFile, error) {
    fullPath := filepath.Join("questions", filename)
    fmt.Println("fullpath", fullPath);
    data, err := ioutil.ReadFile(fullPath) // Read the YAML file
    if err != nil {
        return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
    }

    var qf QuestionFile
    qf.ID = fmt.Sprintf("%x", rand.Intn(999999))
    if err := yaml.Unmarshal(data, &qf); err != nil {
		    return nil, fmt.Errorf("failed to parse YAML: %v", err)
    }
    qf.Filename = filename

    if len(qf.Questions) == 0 {
		    return nil, fmt.Errorf("no questions found in %s", filename)
    }
    qf.ShuffleQuestions()
    return &qf, nil
}

func (qf *QuestionFile) ShuffleQuestions() {
    rand.Shuffle(len(qf.Questions), func(i, j int) {
        qf.Questions[i], qf.Questions[j] = qf.Questions[j], qf.Questions[i]
    })
}
