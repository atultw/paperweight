package paperweight

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

type JSONItem json.RawMessage

func (j JSONItem) Data() []byte {
	return j
}

type MultiSource[T Source] interface {
	Get() (T, error)
	Next() bool
}

type Source interface {
	Data() []byte
}

type MultiFileSource struct {
	InPaths []string
	Index   int
}

func (f *MultiFileSource) Next() bool {
	if f.Index >= len(f.InPaths)-1 {
		return false
	}
	f.Index++
	return true
}

type File struct {
	Path     string
	Contents []byte
}

func (f File) Data() []byte {
	return f.Contents
}

func (f *MultiFileSource) Get() (File, error) {
	path := f.InPaths[f.Index]
	contents, err := ioutil.ReadFile(filepath.Join(path))
	if err != nil {
		return File{}, err
	}
	return File{
		Path:     path,
		Contents: contents,
	}, nil
}

func GlobFileSource(glob string) (MultiFileSource, error) {
	paths, err := filepath.Glob(glob)
	if err != nil {
		return MultiFileSource{}, err
	}
	return MultiFileSource{
		InPaths: paths,
		Index:   -1,
	}, nil
}

type MultiNetworkSource struct {
	Models []json.RawMessage
	Index  int
}

func (m *MultiNetworkSource) Get() (JSONItem, error) {
	return JSONItem(m.Models[m.Index]), nil
}

func (m *MultiNetworkSource) Next() bool {
	if m.Index >= len(m.Models)-1 {
		return false
	}
	m.Index++
	return true
}

func NewMultiNetworkSource(url string) (MultiNetworkSource, error) {
	var res MultiNetworkSource
	resp, err := http.Get(url)
	if err != nil {
		return res, err
	}
	err = json.NewDecoder(resp.Body).Decode(&res.Models)
	if err != nil {
		return res, err
	}
	return res, nil
}

func Passthrough[T any](file T) (T, error) {
	return file, nil
}

func JSONDecode[T any](message JSONItem) (T, error) {
	var res T
	err := json.Unmarshal(message, &res)
	return res, err
}
