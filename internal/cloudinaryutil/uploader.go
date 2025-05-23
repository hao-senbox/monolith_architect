package cloudinaryutil

import (
	"context"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryUploader struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryUploader(cloudinaryURL string) (*CloudinaryUploader, error) {
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, err
	}
	return &CloudinaryUploader{cld: cld}, nil
}

func (u *CloudinaryUploader) UploadImage(ctx context.Context, image string, folderName string) (string, error) {
	result, err := u.cld.Upload.Upload(ctx, image, uploader.UploadParams{
		Folder: folderName,
	})
	if err != nil {
		return "", err
	}
	return result.SecureURL, nil
}