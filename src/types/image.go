package types

type ImagesDetails struct {
	Url string
}

type Image struct {
	Details []ImagesDetails
}

type Headers = map[string]string
