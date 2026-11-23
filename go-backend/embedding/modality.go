package embedding

type Modality string

const (
	Text  Modality = "text"
	Image Modality = "image"
	Audio Modality = "audio"
)

func (a Modality) Valid() bool {
	switch a {
	case Text, Image, Audio:
		return true
	}

	return false
}

func (a Modality) String() string {
	if !a.Valid() {
		return "unknown(" + string(a) + ")"
	}
	return string(a)
}
