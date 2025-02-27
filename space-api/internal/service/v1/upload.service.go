package service

import (
	"crypto/md5"
	"fmt"
	"image"

	// 通过副作用注册支持的图片编码器

	_ "image/jpeg"
	_ "image/png"
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
	"strconv"
	"strings"

	"github.com/chai2010/webp"
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
	case i == 0: // 无配置, 使用默认的全局上传大小限制
	case i == 1:
		// 重置文件上传大小
		inbound.ResetUploadFileLimitSize(ctx)
		// 额外的设置上传文件大小限制
		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, int64(maxSize[0]))
	default:
		err = util.CreateBizErr("仅支持一个可选的文件大小项", fmt.Errorf("just support 1 file size unit, but got :%d", i))
		return
	}

	formPartFile, e := ctx.FormFile("file")
	if e != nil {
		err = util.CreateBizErr("获取文件失败:"+e.Error(), e)
		return
	}
	formMD5, _ := ctx.GetPostForm("md5")

	fileOp := biz.FileRecord
	// 是否已经存在文件, 存在的情况下直接返回存储的相对路径
	if formMD5 != "" {
		f, e := fileOp.WithContext(ctx).
			Where(fileOp.Checksum.Eq(formMD5)).
			Take()
		if e == nil {
			resp = f
			return
		}
	}

	formFile, e := formPartFile.Open()
	if e != nil {
		err = util.CreateBizErr("读取文件流失败", e)
		return
	} else {
		defer formFile.Close()
	}

	cf := conf.ProjectConf.GetAppConf()
	// 在数据库中的存储的相对位置名称
	var location = strings.ReplaceAll(uuid.NewString(), "-", "") + path.Ext(formPartFile.Filename)
	// 实际写入文件系统的位置
	var fsWritePath = path.Join(cf.StaticDir, location)
	// 创建文件
	writeFile, e := os.Create(fsWritePath)
	if e != nil {
		err = util.CreateBizErr("上传文件失败", e)
		return
	} else {
		defer writeFile.Close()
	}

	hasher := md5.New()
	// 写入和哈希计算同时进行
	realSize, e := io.Copy(io.MultiWriter(writeFile, hasher), formFile)
	// 写入或者计算失败
	if e != nil {
		os.Remove(fsWritePath)
		err = util.CreateBizErr("上传文件失败", e)
		return
	}
	md5Count := fmt.Sprintf("%x", hasher.Sum(nil))
	// 和表单声明的 MD5 值不一致
	if formMD5 != "" && md5Count != formMD5 {
		// clean up
		os.Remove(fsWritePath)
		err = util.CreateBizErr("文件校验值不匹配", fmt.Errorf("md5 not matched: client present: %s, but got: %s", formMD5, md5Count))
		return
	}

	// 代表要创建的新文件
	newFile := &model.FileRecord{
		BaseColumn: model.BaseColumn{
			ID: id.GetSnowFlakeNode().Generate().Int64(),
		},
		FileName:      formPartFile.Filename,
		LocalLocation: location,
		Extension:     path.Ext(formPartFile.Filename),
		FileSize:      realSize,
		Category:      strings.TrimSpace(ctx.Request.FormValue("cat")),
		ChecksumType:  "md5",
		Checksum:      md5Count,
		PubAvailable:  1,
	}

	// 判断是否需要进行创建
	txErr := biz.Q.Transaction(func(tx *biz.Query) error {
		fileTx := tx.FileRecord
		// 再次进行匹配, 如果文件系统中已经存在, 直接返回即可, 省略重复创建
		f1, e2 := fileTx.WithContext(ctx).Where(
			fileTx.Checksum.Eq(md5Count),
			// fileTx.Extension.Eq(newFile.Extension), // 文件拓展也一致的情况
		).Take()
		// 存在相同文件
		if e2 == nil {
			// 设置相对位置
			resp = f1
			return nil
		}

		// 需要进行创建的情况下
		e = fileTx.WithContext(ctx).Create(newFile)
		if e != nil {
			return e
		}
		// 返回相对存储位置, 提供实际拼接使用
		f2, e := fileTx.WithContext(ctx).Where(fileTx.ID.Eq(newFile.ID)).Take()
		if e != nil {
			return fmt.Errorf("数据同步失败")
		}
		resp = f2
		return nil
	})

	if txErr != nil {
		os.Remove(fsWritePath)
		err = util.CreateBizErr("上传文件失败", txErr)
		return
	}

	return
}

// UploadImage2Webp 上传图片并转码为 webp 格式
func (d *_localUploadService) UploadImage2Webp(ctx *gin.Context, maxSize ...constants.MemoryByteSize) (resp *model.FileRecord, err error) {
	i := len(maxSize)
	switch {
	case i == 0: // 无配置, 使用默认的全局上传大小限制
	case i == 1:
		// 重置文件上传大小
		inbound.ResetUploadFileLimitSize(ctx)
		// 额外的设置上传文件大小限制
		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, int64(maxSize[0]))
	default:
		err = util.CreateBizErr(
			"仅支持一个可选的文件大小项",
			fmt.Errorf("just support 1 file size unit, but got :%d", i),
		)
		return
	}

	formPartFile, e := ctx.FormFile("file")
	if e != nil {
		err = util.CreateBizErr("获取文件失败:"+e.Error(), e)
		return
	}

	openedFile, e := formPartFile.Open()
	if e != nil {
		err = util.CreateBizErr("读取文件流失败", e)
	} else {
		defer openedFile.Close()
	}

	cf := conf.ProjectConf.GetAppConf()
	// 存储在数据库记录里的相对位置名称
	location := strings.ReplaceAll(uuid.NewString(), "-", "") + ".webp"
	// 实际写入文件系统的位置
	fsWritePath := path.Join(cf.StaticDir, location)

	// 创建文件
	outFile, err := os.Create(fsWritePath)
	if err != nil {
		err = util.CreateBizErr("转码文件失败", err)
		return
	} else {
		defer outFile.Close()
	}

	// 解析设置的转码质量
	qualityStr := ctx.Request.FormValue("quality")
	parsedQuality := 90
	switch {
	case qualityStr == "":
	default:
		v, e := strconv.Atoi(qualityStr)
		if e != nil || (v > 100 || v < 1) {
			err = util.CreateBizErr("图片压缩质量期待值为: 1-100", e)
			os.Remove(fsWritePath)
			return
		}
		parsedQuality = v
	}

	// 对于 outWriterWrapper 字节写入数包装
	outWriterWrapper, byteWriteCounter := util.NewByteWriteCountWrapper(outFile)
	hasher := md5.New()
	fuseWriter := io.MultiWriter(outWriterWrapper, hasher)

	img, _, e := image.Decode(openedFile)
	if e != nil {
		err = util.CreateBizErr("解码图片失败", e)
		os.Remove(fsWritePath)
		return
	}
	// 同时写入文件和计算 md5 值和计算实际文件大小
	err = webp.Encode(fuseWriter, img, &webp.Options{Quality: float32(parsedQuality), Exact: true})
	if err != nil {
		err = util.CreateBizErr("图片转码失败", e)
		os.Remove(fsWritePath)
		return
	}

	md5Val := fmt.Sprintf("%x", hasher.Sum(nil))
	newFile := &model.FileRecord{
		BaseColumn: model.BaseColumn{
			ID: id.GetSnowFlakeNode().Generate().Int64(),
		},
		FileName:      formPartFile.Filename[:len(formPartFile.Filename)-len(path.Ext(formPartFile.Filename))] + ".webp",
		LocalLocation: location,
		Extension:     ".webp",
		FileSize:      byteWriteCounter(), // 获取转码后的实际文件大小
		Category:      strings.TrimSpace(ctx.Request.FormValue("cat")),
		ChecksumType:  "md5",
		Checksum:      md5Val,
		PubAvailable:  1,
	}

	// 判断是否需要进行创建
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		fileTx := tx.FileRecord
		// 查找是否已经存在文件
		f1, e := fileTx.WithContext(ctx).
			Where(fileTx.Checksum.Eq(md5Val), fileTx.Extension.Eq(".webp")).
			Take()
		// 文件系统中已经存在相同文件
		if e == nil {
			resp = f1
			os.Remove(fsWritePath)
			return nil
		}

		// 需要进行创建的情况下
		e = fileTx.WithContext(ctx).Create(newFile)
		if e != nil {
			os.Remove(fsWritePath)
			return e
		}

		f2, e := fileTx.WithContext(ctx).
			Where(fileTx.ID.Eq(newFile.ID)).
			Take()
		if e != nil {
			os.Remove(fsWritePath)
			return fmt.Errorf("数据同步失败")
		}

		resp = f2
		return nil
	})

	if err != nil {
		err = util.CreateBizErr("上传文件失败", err)
		return
	}

	return
}
