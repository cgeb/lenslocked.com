package models

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

type Image struct {
	GalleryID uint
	Filename  string
}

// Path is used to build the absolute path used to reference this image
// via a web request.
func (i *Image) Path() string {
	temp := url.URL{
		Path: "/" + i.RelativePath(),
	}
	return temp.String()
}

// RelativePath is used to build the path to this image on our local
// disk, relative to where our Go application is run from.
// Note: ToSlash is likely only needed in Windows systems, as they use a backslash
// (\) instead of a forward slash (/) to separate directories and files. The
// ToSlash function works by replacing all forward slashes with backslashes on
// operating systems that use backslashes.
func (i *Image) RelativePath() string {
	galleryID := fmt.Sprintf("%v", i.GalleryID)
	return filepath.ToSlash(filepath.Join("images", "galleries", galleryID, i.Filename))
}

type ImageService interface {
	Create(galleryID uint, r io.Reader, filename string) error
	ByGalleryID(galleryID uint) ([]Image, error)
	Delete(i *Image) error
}

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct{}

func (is *imageService) mkImageDir(galleryID uint) (string, error) {
	galleryDir := is.imageDir(galleryID)
	err := os.MkdirAll(galleryDir, 0755)
	if err != nil {
		return "", err
	}
	return galleryDir, nil
}

func (is *imageService) Create(galleryID uint, r io.Reader, filename string) error {
	dir, err := is.mkImageDir(galleryID)
	if err != nil {
		return err
	}
	dst, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}
	return nil
}

func (is *imageService) ByGalleryID(galleryID uint) ([]Image, error) {
	dir := is.imageDir(galleryID)
	strings, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return nil, err
	}
	// Setup the Image slice we are returning
	ret := make([]Image, len(strings))
	for i, imgStr := range strings {
		ret[i] = Image{
			Filename:  filepath.Base(imgStr),
			GalleryID: galleryID,
		}
	}
	return ret, nil
}

func (is *imageService) imageDir(galleryID uint) string {
	return filepath.Join("images", "galleries",
		fmt.Sprintf("%v", galleryID))
}

func (is *imageService) Delete(i *Image) error {
	return os.Remove(i.RelativePath())
}
