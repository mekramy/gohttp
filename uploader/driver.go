package uploader

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/gabriel-vasile/mimetype"
	"github.com/inhies/go-bytesize"
	"github.com/mekramy/goutils"
	"github.com/valyala/fasthttp"
)

type uploader struct {
	opt   Option
	file  *multipart.FileHeader
	name  string
	root  string
	saved bool
}

func (u *uploader) IsNil() bool {
	return u.file == nil
}

func (u *uploader) ValidateSize(min, max string) (bool, error) {
	// Invalidate nil file
	if u.IsNil() {
		return false, nil
	}

	// Parse min string
	minSize, err := bytesize.Parse(min)
	if err != nil {
		return false, err
	}

	// Parse max string
	maxSize, err := bytesize.Parse(max)
	if err != nil {
		return false, err
	}

	// Validate
	size := u.file.Size
	return size >= int64(minSize) && size <= int64(maxSize), nil
}

func (u *uploader) ValidateMime(mimes ...string) (bool, error) {
	// Invalidate nil file
	if u.IsNil() {
		return false, nil
	}

	// Read file content
	f, err := u.file.Open()
	if err != nil {
		return false, err
	}
	defer f.Close()

	// Validate mime
	mime, err := mimetype.DetectReader(f)
	if err != nil {
		return false, err
	}

	return mimetype.EqualsAny(mime.String(), mimes...), nil
}

func (u *uploader) Path() string {
	// Skip nil file
	if u.IsNil() {
		return ""
	}

	return goutils.NormalizePath(u.root, u.name)
}

func (u *uploader) URL() string {
	// Skip nil file
	if u.IsNil() {
		return ""
	}

	return goutils.AbsoluteURL(u.opt.prefix, u.Path())
}

func (u *uploader) Save() error {
	// Skip nil file or saved
	if u.IsNil() || u.saved {
		return nil
	}

	dest := u.Path()

	// Check if exists
	exists, err := goutils.FileExists(dest)
	if err != nil {
		return err
	} else if exists {
		return fmt.Errorf("%s file exists", dest)
	}

	// Save
	err = fasthttp.SaveMultipartFile(u.file, dest)
	if err != nil {
		return err
	}

	u.saved = true
	return nil
}

func (u *uploader) Delete() error {
	// Skip nil file or not saved
	if u.IsNil() || !u.saved {
		return nil
	}

	// Delete
	err := os.Remove(u.Path())
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return err
}

func (u *uploader) SafeDelete() {
	err := u.Delete()
	if u.opt.queue == nil {
		return
	}

	if err != nil {
		u.opt.queue.Push(u.Path())
	}
}
