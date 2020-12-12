package model

import (
	"context"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"mime/multipart"
	"qa/util"
)

//var AccessKey = util.AccessKey
//var SecretKey = util.SecretKey
//var Bucket = util.Bucket
//var ImgUrl = util.QiniuSever

// 上传文件  -->  要重启服务器才能加载配置文件，而不是热加载的重启，否则会读取不到配置文件的信息
func UpLoadFile(file multipart.File, fileSize int64) (string, util.MyCode) {
	putPolicy := storage.PutPolicy{
		Scope: util.Bucket,
	}

	mac := qbox.NewMac(util.AccessKey, util.SecretKey)
	upToken := putPolicy.UploadToken(mac)


	cfg := storage.Config{
		Zone:          &storage.ZoneHuadong,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}

	putExtra := storage.PutExtra{}

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	err := formUploader.PutWithoutKey(context.Background(), &ret, upToken, file, fileSize, &putExtra)
	if err != nil {
		return "", util.UploadFail
	}
	url := util.QiniuSever + ret.Key
	return url, util.CodeSuccess

}