package doctrans

import (
	"io"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

type pdfDoc struct {
	docType DocType
	read    *pdfcpu.Context
	write   *pdfcpu.Context
}

func (pd *pdfDoc) GetType() DocType {
	return pd.docType
}

func NewPdfDocFromReader(r io.ReadSeeker) (d Doc, err error) {
	ctx, err := pdfcpu.NewContext(r, pdfcpu.NewDefaultConfiguration())
	if err != nil {
		return nil, err
	}
	return &pdfDoc{
		docType: DOC_PDF,
		read:    ctx,
	}, nil
}
