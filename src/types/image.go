package types

type ImagesDetails struct {
	Url      string
	FileName string
}

type Image struct {
	Details []ImagesDetails
}
