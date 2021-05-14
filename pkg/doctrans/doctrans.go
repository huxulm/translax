package doctrans

type DocType int

const (
	DOC_PDF DocType = iota
	DOC_WORD
	DOC_TXT
)

type Doc interface {
	GetType() DocType
}

func IsSupported(t DocType) bool {
	return t == DOC_PDF || t == DOC_TXT || t == DOC_WORD
}

func TransDoc(src Doc, dest Doc) {

}
