package csv_writer

import (
	"io"
	"log"
	"sync"

	"github.com/ivankwongtszfung/engineering-blog-web-crawler/entity/blog"

	"github.com/mohae/struct2csv"
)

type CSVWriter struct {
	w           *struct2csv.Writer
	headerLock  sync.Once
	writeHeader bool
}

func writeHeader(csvWriter *struct2csv.Writer, article *blog.Article) error {
	if err := csvWriter.WriteColNames(*article); err != nil {
		log.Println(err)
		return err
	}
	csvWriter.Flush()
	return nil
}

func writeLine(csvWriter *struct2csv.Writer, article *blog.Article) error {
	if err := csvWriter.WriteStruct(*article); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func NewCSVWriter(w io.Writer, writeHeader bool) *CSVWriter {
	return &CSVWriter{w: struct2csv.NewWriter(w), writeHeader: writeHeader}
}

func (cw *CSVWriter) Write(article *blog.Article) {
	// things in sync.Once only run once, so no performance issue
	cw.headerLock.Do(func() {
		if cw.writeHeader {
			writeHeader(cw.w, article)
			cw.writeHeader = false
		}
	})
	writeLine(cw.w, article)
}

func (cw *CSVWriter) Flush() {
	cw.w.Flush()
}
