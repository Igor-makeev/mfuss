package repositories

import (
	"bufio"
	"encoding/json"
	"mfuss/internal/entity"
	"os"
)

type Dump struct {
	file    *os.File
	writer  *bufio.Writer
	scanner *bufio.Scanner
}

func NewDump(filename string) (*Dump, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &Dump{
		file:    file,
		writer:  bufio.NewWriter(file),
		scanner: bufio.NewScanner(file),
	}, nil
}

func (d *Dump) SaveData(ms map[string]entity.ShortURL) error {
	for _, value := range ms {

		err := d.WriteURL(&value)
		if err != nil {
			return err
		}
	}
	return nil

}

func (d *Dump) LoadData(ms map[string]entity.ShortURL) error {

	d.scanner.Split(bufio.ScanLines)

	for d.scanner.Scan() {
		data := d.scanner.Bytes()
		URL := entity.ShortURL{}
		err := json.Unmarshal(data, &URL)
		if err != nil {
			return err
		}
		ms[URL.ID] = URL

	}
	if err := d.scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (d *Dump) WriteURL(URL *entity.ShortURL) error {
	data, err := json.Marshal(&URL)
	if err != nil {
		return err
	}
	if _, err := d.writer.Write(data); err != nil {
		return err
	}
	if err := d.writer.WriteByte('\n'); err != nil {
		return err
	}
	return d.writer.Flush()

}

func (d *Dump) Close() error {

	return d.file.Close()
}
