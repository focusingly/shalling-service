package service

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"space-api/conf"
	"space-api/constants"
	"space-api/middleware/inbound"
	"space-api/util"
	"space-api/util/id"
	"space-domain/dao/biz"
	"space-domain/model"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	_localUploadService struct{}
)

var DefaultUploadService = &_localUploadService{}

// Upload 上传文件到本地
// 表单参数规范: file: File(必须), md5: text(可选)
func (d *_localUploadService) Upload(ctx *gin.Context, maxSize ...constants.MemoryByteSize) (resp *model.FileRecord, err error) {
	i := len(maxSize)
	switch {
	case i == 0:
	case i == 1:
		// 重置文件上传大小
		inbound.ResetUploadFileLimitSize(ctx)
		// 额外的设置上传文件大小限制
		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, int64(maxSize[0]))
	default:
		err = util.CreateBizErr("仅支持一个可选的文件大小项", fmt.Errorf("just support 1 file size unit, but got :%d", i))
		return
	}

	formMD5 := ctx.Request.FormValue("md5")
	fileOp := biz.FileRecord
	// 是否已经存在文件, 存在的情况下直接返回存储路径
	if formMD5 != "" {
		f, e := fileOp.WithContext(ctx).Where(fileOp.Checksum.Eq(formMD5)).Take()
		if e == nil {
			resp = f
			return
		}
	}

	file, e := ctx.FormFile("file")
	if e != nil {
		err = util.CreateBizErr("获取文件失败:"+e.Error(), e)
		return
	}

	formFile, e := file.Open()
	if e != nil {
		err = util.CreateBizErr("读取文件流失败", e)
	}
	cf := conf.ProjectConf.GetAppConf()
	// 相对位置名称
	var location = uuid.NewString() + path.Ext(file.Filename)
	// 实际写入文件系统的位置
	var writePath = path.Join(cf.StaticDir, location)
	// 代表要创建的新文件
	newFile := &model.FileRecord{
		BaseColumn: model.BaseColumn{
			ID: id.GetSnowFlakeNode().Generate().Int64(),
		},
		FileName:      file.Filename,
		LocalLocation: location,
		Extension:     path.Ext(file.Filename),
		// FileSize:      file.Size,
		Category:     strings.TrimSpace(ctx.Request.FormValue("cat")),
		ChecksumType: "md5",
		// Checksum:     "",
		PubAvailable: 1,
	}

	fd, e := os.Create(writePath)
	if e != nil {
		err = util.CreateBizErr("上传文件失败", e)
		return
	}
	defer fd.Close()

	md5Hasher := md5.New()
	// 写入
	realSize, e := io.Copy(io.MultiWriter(fd, md5Hasher), formFile)
	if e != nil {
		os.Remove(writePath)
		err = util.CreateBizErr("上传文件失败", e)
		return
	}
	md5Count := fmt.Sprintf("%x", md5Hasher.Sum(nil))

	if formMD5 != "" && md5Count != formMD5 {
		os.Remove(writePath)
		err = util.CreateBizErr("文件校验值不匹配", fmt.Errorf("md5 not matched: client present: %s, but got: %s", formMD5, md5Count))
		return
	}

	newFile.Checksum = md5Count
	newFile.FileSize = realSize

	// 需要进行创建
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		fileTx := tx.FileRecord

		// 再次进行匹配, 如果文件系统中已经存在, 直接返回即可, 省略重复创建
		f, e2 := fileOp.WithContext(ctx).Where(fileOp.Checksum.Eq(md5Count)).Take()
		if e2 == nil {
			// 设置相对位置
			resp = f
			// 移除写入的重复文件
			os.Remove(writePath)
			return nil
		}

		// 需要进行创建的情况下
		e = fileTx.WithContext(ctx).Create(newFile)
		if e != nil {
			// 删除文件
			os.Remove(writePath)
			return e
		}

		// 返回相对存储位置, 提供实际拼接使用
		f, e := fileTx.WithContext(ctx).Where(fileTx.ID.Eq(newFile.ID)).Take()
		if e != nil {
			return fmt.Errorf("数据同步失败")
		}
		resp = f
		return nil
	})

	if err != nil {
		err = util.CreateBizErr("上传文件失败", err)
		return
	}

	return
}
