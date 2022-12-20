package repositories

import (
	"bufio"
	"encoding/json"
	"mfuss/internal/entity"
	"os"
)

type FileStorage struct {
	file    *os.File
	writer  *bufio.Writer
	scanner *bufio.Scanner
}

func NewFileStorage(filename string) (*FileStorage, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &FileStorage{
		file:    file,
		writer:  bufio.NewWriter(file),
		scanner: bufio.NewScanner(file),
	}, nil
}

func (fs *FileStorage) SaveData(ms map[string]entity.ShortURL) error {
	for _, value := range ms {

		err := fs.WriteURL(&value)
		if err != nil {
			return err
		}
	}
	return nil

}

func (fs *FileStorage) LoadData(ms map[string]entity.ShortURL) error {

	fs.scanner.Split(bufio.ScanLines)

	for fs.scanner.Scan() {
		data := fs.scanner.Bytes()
		URL := entity.ShortURL{}
		err := json.Unmarshal(data, &URL)
		if err != nil {
			return err
		}
		ms[URL.ID] = URL

	}
	if err := fs.scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (fs *FileStorage) WriteURL(URL *entity.ShortURL) error {
	data, err := json.Marshal(&URL)
	if err != nil {
		return err
	}
	if _, err := fs.writer.Write(data); err != nil {
		return err
	}
	if err := fs.writer.WriteByte('\n'); err != nil {
		return err
	}
	return fs.writer.Flush()

}

func (fs *FileStorage) ReadURL() (*entity.ShortURL, error) {

	if !fs.scanner.Scan() {
		return nil, fs.scanner.Err()
	}

	data := fs.scanner.Bytes()

	URL := entity.ShortURL{}
	err := json.Unmarshal(data, &URL)
	if err != nil {
		return nil, err
	}

	return &URL, nil
}

func (fs *FileStorage) Close() error {

	return fs.file.Close()
}
