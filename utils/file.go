// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/mattermost/platform/model"
	s3 "github.com/minio/minio-go"
)

func ReadFile(path string) ([]byte, *model.AppError) {
	if Cfg.FileSettings.DriverName == model.IMAGE_DRIVER_S3 {
		endpoint := Cfg.FileSettings.AmazonS3Endpoint
		accessKey := Cfg.FileSettings.AmazonS3AccessKeyId
		secretKey := Cfg.FileSettings.AmazonS3SecretAccessKey
		secure := *Cfg.FileSettings.AmazonS3SSL
		s3Clnt, err := s3.New(endpoint, accessKey, secretKey, secure)
		if err != nil {
			return nil, model.NewLocAppError("ReadFile", "api.file.read_file.s3.app_error", nil, err.Error())
		}
		bucket := Cfg.FileSettings.AmazonS3Bucket
		minioObject, err := s3Clnt.GetObject(bucket, path)
		defer minioObject.Close()
		if err != nil {
			return nil, model.NewLocAppError("ReadFile", "api.file.read_file.s3.app_error", nil, err.Error())
		}
		if f, err := ioutil.ReadAll(minioObject); err != nil {
			return nil, model.NewLocAppError("ReadFile", "api.file.read_file.s3.app_error", nil, err.Error())
		} else {
			return f, nil
		}
	} else if Cfg.FileSettings.DriverName == model.IMAGE_DRIVER_LOCAL {
		if f, err := ioutil.ReadFile(Cfg.FileSettings.Directory + path); err != nil {
			return nil, model.NewLocAppError("ReadFile", "api.file.read_file.reading_local.app_error", nil, err.Error())
		} else {
			return f, nil
		}
	} else {
		return nil, model.NewAppError("ReadFile", "api.file.read_file.configured.app_error", nil, "", http.StatusNotImplemented)
	}
}

func MoveFile(oldPath, newPath string) *model.AppError {
	if Cfg.FileSettings.DriverName == model.IMAGE_DRIVER_S3 {
		endpoint := Cfg.FileSettings.AmazonS3Endpoint
		accessKey := Cfg.FileSettings.AmazonS3AccessKeyId
		secretKey := Cfg.FileSettings.AmazonS3SecretAccessKey
		secure := *Cfg.FileSettings.AmazonS3SSL
		s3Clnt, err := s3.New(endpoint, accessKey, secretKey, secure)
		if err != nil {
			return model.NewLocAppError("moveFile", "api.file.write_file.s3.app_error", nil, err.Error())
		}
		bucket := Cfg.FileSettings.AmazonS3Bucket

		var copyConds = s3.NewCopyConditions()
		if err = s3Clnt.CopyObject(bucket, newPath, "/"+path.Join(bucket, oldPath), copyConds); err != nil {
			return model.NewLocAppError("moveFile", "api.file.move_file.delete_from_s3.app_error", nil, err.Error())
		}
		if err = s3Clnt.RemoveObject(bucket, oldPath); err != nil {
			return model.NewLocAppError("moveFile", "api.file.move_file.delete_from_s3.app_error", nil, err.Error())
		}
	} else if Cfg.FileSettings.DriverName == model.IMAGE_DRIVER_LOCAL {
		if err := os.MkdirAll(filepath.Dir(Cfg.FileSettings.Directory+newPath), 0774); err != nil {
			return model.NewLocAppError("moveFile", "api.file.move_file.rename.app_error", nil, err.Error())
		}

		if err := os.Rename(Cfg.FileSettings.Directory+oldPath, Cfg.FileSettings.Directory+newPath); err != nil {
			return model.NewLocAppError("moveFile", "api.file.move_file.rename.app_error", nil, err.Error())
		}
	} else {
		return model.NewLocAppError("moveFile", "api.file.move_file.configured.app_error", nil, "")
	}

	return nil
}

func WriteFile(f []byte, path string) *model.AppError {
	if Cfg.FileSettings.DriverName == model.IMAGE_DRIVER_S3 {
		endpoint := Cfg.FileSettings.AmazonS3Endpoint
		accessKey := Cfg.FileSettings.AmazonS3AccessKeyId
		secretKey := Cfg.FileSettings.AmazonS3SecretAccessKey
		secure := *Cfg.FileSettings.AmazonS3SSL
		s3Clnt, err := s3.New(endpoint, accessKey, secretKey, secure)
		if err != nil {
			return model.NewLocAppError("WriteFile", "api.file.write_file.s3.app_error", nil, err.Error())
		}
		bucket := Cfg.FileSettings.AmazonS3Bucket
		ext := filepath.Ext(path)

		if model.IsFileExtImage(ext) {
			_, err = s3Clnt.PutObject(bucket, path, bytes.NewReader(f), model.GetImageMimeType(ext))
		} else {
			_, err = s3Clnt.PutObject(bucket, path, bytes.NewReader(f), "binary/octet-stream")
		}
		if err != nil {
			return model.NewLocAppError("WriteFile", "api.file.write_file.s3.app_error", nil, err.Error())
		}
	} else if Cfg.FileSettings.DriverName == model.IMAGE_DRIVER_LOCAL {
		if err := writeFileLocally(f, Cfg.FileSettings.Directory+path); err != nil {
			return err
		}
	} else {
		return model.NewLocAppError("WriteFile", "api.file.write_file.configured.app_error", nil, "")
	}

	return nil
}

func writeFileLocally(f []byte, path string) *model.AppError {
	if err := os.MkdirAll(filepath.Dir(path), 0774); err != nil {
		directory, _ := filepath.Abs(filepath.Dir(path))
		return model.NewLocAppError("WriteFile", "api.file.write_file_locally.create_dir.app_error", nil, "directory="+directory+", err="+err.Error())
	}

	if err := ioutil.WriteFile(path, f, 0644); err != nil {
		return model.NewLocAppError("WriteFile", "api.file.write_file_locally.writing.app_error", nil, err.Error())
	}

	return nil
}

func RemoveFile(path string) *model.AppError {
	if Cfg.FileSettings.DriverName == model.IMAGE_DRIVER_S3 {
		endpoint := Cfg.FileSettings.AmazonS3Endpoint
		accessKey := Cfg.FileSettings.AmazonS3AccessKeyId
		secretKey := Cfg.FileSettings.AmazonS3SecretAccessKey
		secure := *Cfg.FileSettings.AmazonS3SSL
		s3Clnt, err := s3.New(endpoint, accessKey, secretKey, secure)
		if err != nil {
			return model.NewLocAppError("RemoveFile", "api.file.remove_file.s3.app_error", nil, err.Error())
		}

		bucket := Cfg.FileSettings.AmazonS3Bucket
		if err := s3Clnt.RemoveObject(bucket, path); err != nil {
			return model.NewLocAppError("RemoveFile", "api.file.remove_file.s3.app_error", nil, err.Error())
		}
	} else if Cfg.FileSettings.DriverName == model.IMAGE_DRIVER_LOCAL {
		if err := os.Remove(Cfg.FileSettings.Directory + path); err != nil {
			return model.NewLocAppError("RemoveFile", "api.file.remove_file.local.app_error", nil, err.Error())
		}
	} else {
		return model.NewLocAppError("RemoveFile", "api.file.write_file.configured.app_error", nil, "")
	}

	return nil
}

func getPathsFromObjectInfos(in <-chan s3.ObjectInfo) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer close(out)

		for {
			info, done := <-in

			if !done {
				break
			}

			out <- info.Key
		}
	}()

	return out
}

func RemoveDirectory(path string) *model.AppError {
	if Cfg.FileSettings.DriverName == model.IMAGE_DRIVER_S3 {
		endpoint := Cfg.FileSettings.AmazonS3Endpoint
		accessKey := Cfg.FileSettings.AmazonS3AccessKeyId
		secretKey := Cfg.FileSettings.AmazonS3SecretAccessKey
		secure := *Cfg.FileSettings.AmazonS3SSL
		s3Clnt, err := s3.New(endpoint, accessKey, secretKey, secure)
		if err != nil {
			return model.NewLocAppError("RemoveDirectory", "api.file.remove_directory.s3.app_error", nil, err.Error())
		}

		doneCh := make(chan struct{})

		bucket := Cfg.FileSettings.AmazonS3Bucket
		for err := range s3Clnt.RemoveObjects(bucket, getPathsFromObjectInfos(s3Clnt.ListObjects(bucket, path, true, doneCh))) {
			if err.Err != nil {
				doneCh <- struct{}{}
				return model.NewLocAppError("RemoveDirectory", "api.file.remove_directory.s3.app_error", nil, err.Err.Error())
			}
		}

		close(doneCh)
	} else if Cfg.FileSettings.DriverName == model.IMAGE_DRIVER_LOCAL {
		if err := os.RemoveAll(Cfg.FileSettings.Directory + path); err != nil {
			return model.NewLocAppError("RemoveDirectory", "api.file.remove_directory.local.app_error", nil, err.Error())
		}
	} else {
		return model.NewLocAppError("RemoveDirectory", "api.file.write_file.configured.app_error", nil, "")
	}

	return nil
}
