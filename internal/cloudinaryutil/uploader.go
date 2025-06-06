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

func (u *CloudinaryUploader) UploadImage(ctx context.Context, image string, folderName string) (string, string, error) {
	result, err := u.cld.Upload.Upload(ctx, image, uploader.UploadParams{
		Folder: folderName,
	})
	if err != nil {
		return "", "", err
	}
	return result.SecureURL, result.PublicID, nil
}

func (u *CloudinaryUploader) DeleteImage(ctx context.Context, publicID string) error {
	_, err := u.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return err
	}
	return nil
}