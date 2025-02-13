package dto

import "space-domain/model"

type (
	// 全文检索
	SearchHighlight struct {
		Title                    string   `json:"title" yaml:"title" xml:"title" toml:"title"`
		TitleHighLightIndex      []int    `json:"titleHighLightIndex" yaml:"titleHighLightIndex" xml:"titleHighLightIndex" toml:"titleHighLightIndex"`                     // 标题需要进行高亮的位置, 无论是否存在, 都完整返回
		SubContent               string   `json:"subContent" yaml:"subContent" xml:"subContent" toml:"subContent"`                                                         // 截取的文章部分, 如果不存在, 直接为空字符串, 存在的话则截取一定长度字符串
		SubContentHighLightIndex []int    `json:"subContentHighLightIndex" yaml:"subContentHighLightIndex" xml:"subContentHighLightIndex" toml:"subContentHighLightIndex"` // 截取内容需要高亮的位置
		Category                 *string  `json:"category" yaml:"category" xml:"category" toml:"category"`
		Tags                     []string `json:"tags" yaml:"tags" xml:"tags" toml:"tags"`
		CreatedAt                int64    `json:"createdAt" yaml:"createdAt" xml:"createdAt" toml:"createdAt"` // 文章创建时间
		PubAt                    *int64   `json:"pubAt" yaml:"pubAt" xml:"pubAt" toml:"pubAt"`                 // 文章发表时间
		Weight                   *int     `json:"weight" yaml:"weight" xml:"weight" toml:"weight"`
	}
	GlobalSearchReq struct {
		Keyword string `form:"keyword" json:"keyword" yaml:"keyword" xml:"keyword" toml:"keyword"`
		*BasePageParam
	}
	GlobalSearchResp struct {
		Keyword string                           `json:"keyword" yaml:"keyword" xml:"keyword" toml:"keyword"`
		List    model.PageList[*SearchHighlight] `json:"list" yaml:"list" xml:"list" toml:"list"`
	}
	GetSearchIndexPagesReq struct {
		BasePageParam
	}
	GetSearchIndexPagesResp struct {
		model.PageList[*model.Sqlite3KeywordDoc]
	}
)
