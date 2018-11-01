package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rwirdemann/restvoice/cha05/domain"
)

type InvoicePresenter interface {
	Present(i interface{})
}

type JSONInvoicePresenter struct {
	writer http.ResponseWriter
}

func NewJSONInvoicePresenter(w http.ResponseWriter) JSONInvoicePresenter {
	return JSONInvoicePresenter{writer: w}
}

func (p JSONInvoicePresenter) Present(i interface{}) {
	if b, err := json.Marshal(i); err == nil {
		p.writer.Write(b)
	}
}

type DefaultPresenter struct {
}

func (p DefaultPresenter) Present(i interface{}) {
}

func NewDefaultPresenter() DefaultPresenter {
	return DefaultPresenter{}
}

type PDFInvoicePresenter struct {
	w http.ResponseWriter
	r *http.Request
}

func NewPDFInvoicePresenter(w http.ResponseWriter, r *http.Request) PDFInvoicePresenter {
	return PDFInvoicePresenter{w: w, r: r}
}

func (p PDFInvoicePresenter) Present(i interface{}) {
	modTime := time.Now()
	invoice := i.(domain.Invoice)
	content := bytes.NewReader(invoice.ToPDF())
	http.ServeContent(p.w, p.r, "invoice.pdf", modTime, content)
}
