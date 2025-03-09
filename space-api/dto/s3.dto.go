package dto

import (
	"space-api/dto/query"
	"space-api/util/performance"
	"space-domain/model"
)

type (
	GetUploadObjectURLReq struct {
		Bucket   string             `json:"bucket" yaml:"bucket" xml:"bucket" toml:"bucket"`
		SHA256   *string            `json:"sha256" yaml:"sha256" xml:"sha256" toml:"sha256"`
		Filename string             `json:"filename" yaml:"filename" xml:"filename" toml:"filename"`
		Metadata map[string]*string `json:"metadata" yaml:"metadata" xml:"metadata" toml:"metadata"`
	}
	GetUploadObjectURLResp struct {
		ObjectKey      string `json:"objectKey" yaml:"objectKey" xml:"objectKey" toml:"objectKey"`
		OriginFileName string `json:"originFileName" yaml:"originFileName" xml:"originFileName" toml:"originFileName"`
		Method         string `json:"method" yaml:"method" xml:"method" toml:"method"`
		ChecksumType   string `json:"checksumType" yaml:"checksumType" xml:"checksumType" toml:"checksumType"`
		Checksum       string `json:"checksum" yaml:"checksum" xml:"checksum" toml:"checksum"`
		URL            string `json:"url" yaml:"url" xml:"url" toml:"url"`
		ExpiredAt      int64  `json:"expiredAt,string" yaml:"expiredAt" xml:"expiredAt" toml:"expiredAt"`
	}

	SyncS3RecordToDatabaseReq struct {
		ObjectKey      string `json:"objectKey" yaml:"objectKey" xml:"objectKey" toml:"objectKey"`
		BucketName     string `json:"bucketName" yaml:"bucketName" xml:"bucketName" toml:"bucketName"`
		OriginFileName string `json:"originFileName" yaml:"originFileName" xml:"originFileName" toml:"originFileName"`
		ChecksumType   string `json:"checksumType" yaml:"checksumType" xml:"checksumType" toml:"checksumType"`
		Checksum       string `json:"checksum" yaml:"checksum" xml:"checksum" toml:"checksum"`
		URL            string `json:"url" yaml:"url" xml:"url" toml:"url"`
	}
	SyncToDatabaseResp struct {
		FileSize     int64  `json:"fileSize" yaml:"fileSize" xml:"fileSize" toml:"fileSize"`
		LinkedDomain string `json:"linkedDomain" yaml:"linkedDomain" xml:"linkedDomain" toml:"linkedDomain"`
		ObjectKey    string `json:"objectKey" yaml:"objectKey" xml:"objectKey" toml:"objectKey"`
		FullVisitURL string `json:"fullVisitURL" yaml:"fullVisitURL" xml:"fullVisitURL" toml:"fullVisitURL"`
	}

	GetS3ObjectPagesReq struct {
		BasePageParam
		CondList  []*query.WhereCond   `json:"condList" yaml:"condList" xml:"condList" toml:"condList"`
		OrderList []*query.OrderColumn `json:"orderList" yaml:"orderList" xml:"orderList" toml:"orderList"`
	}
	GetS3ObjectPagesResp struct {
		model.PageList[*model.S3ObjectRecord]
	}

	DeleteS3ObjectPagesReq struct {
		CondList []*query.WhereCond `json:"condList" yaml:"condList" xml:"condList" toml:"condList"`
	}
	DeleteS3ObjectPagesResp performance.Empty
)
