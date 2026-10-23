package embedding

type Model string

const (
	JinaClipV2 Model = "jina-clip-v2"
)

func (m Model) Valid() bool {
	switch m {
	case JinaClipV2:
		return true
	}

	return false
}

func (m Model) String() string {
	if !m.Valid() {
		return "unknown(" + string(m) + ")"
	}
	return string(m)
}
