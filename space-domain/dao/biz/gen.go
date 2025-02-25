// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package biz

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"gorm.io/gen"

	"gorm.io/plugin/dbresolver"
)

var (
	Q                 = new(Query)
	BlockIPRecord     *blockIPRecord
	Category          *category
	CloudFn           *cloudFn
	Comment           *comment
	CronJob           *cronJob
	FileRecord        *fileRecord
	FriendLink        *friendLink
	LocalUser         *localUser
	MenuGroup         *menuGroup
	MenuLink          *menuLink
	OAuth2User        *oAuth2User
	Post              *post
	PostTagRelation   *postTagRelation
	PubSocialMedia    *pubSocialMedia
	S3ObjectRecord    *s3ObjectRecord
	ServiceConf       *serviceConf
	Sqlite3KeywordDoc *sqlite3KeywordDoc
	Tag               *tag
	UVStatistic       *uVStatistic
	UserLoginSession  *userLoginSession
)

func SetDefault(db *gorm.DB, opts ...gen.DOOption) {
	*Q = *Use(db, opts...)
	BlockIPRecord = &Q.BlockIPRecord
	Category = &Q.Category
	CloudFn = &Q.CloudFn
	Comment = &Q.Comment
	CronJob = &Q.CronJob
	FileRecord = &Q.FileRecord
	FriendLink = &Q.FriendLink
	LocalUser = &Q.LocalUser
	MenuGroup = &Q.MenuGroup
	MenuLink = &Q.MenuLink
	OAuth2User = &Q.OAuth2User
	Post = &Q.Post
	PostTagRelation = &Q.PostTagRelation
	PubSocialMedia = &Q.PubSocialMedia
	S3ObjectRecord = &Q.S3ObjectRecord
	ServiceConf = &Q.ServiceConf
	Sqlite3KeywordDoc = &Q.Sqlite3KeywordDoc
	Tag = &Q.Tag
	UVStatistic = &Q.UVStatistic
	UserLoginSession = &Q.UserLoginSession
}

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:                db,
		BlockIPRecord:     newBlockIPRecord(db, opts...),
		Category:          newCategory(db, opts...),
		CloudFn:           newCloudFn(db, opts...),
		Comment:           newComment(db, opts...),
		CronJob:           newCronJob(db, opts...),
		FileRecord:        newFileRecord(db, opts...),
		FriendLink:        newFriendLink(db, opts...),
		LocalUser:         newLocalUser(db, opts...),
		MenuGroup:         newMenuGroup(db, opts...),
		MenuLink:          newMenuLink(db, opts...),
		OAuth2User:        newOAuth2User(db, opts...),
		Post:              newPost(db, opts...),
		PostTagRelation:   newPostTagRelation(db, opts...),
		PubSocialMedia:    newPubSocialMedia(db, opts...),
		S3ObjectRecord:    newS3ObjectRecord(db, opts...),
		ServiceConf:       newServiceConf(db, opts...),
		Sqlite3KeywordDoc: newSqlite3KeywordDoc(db, opts...),
		Tag:               newTag(db, opts...),
		UVStatistic:       newUVStatistic(db, opts...),
		UserLoginSession:  newUserLoginSession(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	BlockIPRecord     blockIPRecord
	Category          category
	CloudFn           cloudFn
	Comment           comment
	CronJob           cronJob
	FileRecord        fileRecord
	FriendLink        friendLink
	LocalUser         localUser
	MenuGroup         menuGroup
	MenuLink          menuLink
	OAuth2User        oAuth2User
	Post              post
	PostTagRelation   postTagRelation
	PubSocialMedia    pubSocialMedia
	S3ObjectRecord    s3ObjectRecord
	ServiceConf       serviceConf
	Sqlite3KeywordDoc sqlite3KeywordDoc
	Tag               tag
	UVStatistic       uVStatistic
	UserLoginSession  userLoginSession
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:                db,
		BlockIPRecord:     q.BlockIPRecord.clone(db),
		Category:          q.Category.clone(db),
		CloudFn:           q.CloudFn.clone(db),
		Comment:           q.Comment.clone(db),
		CronJob:           q.CronJob.clone(db),
		FileRecord:        q.FileRecord.clone(db),
		FriendLink:        q.FriendLink.clone(db),
		LocalUser:         q.LocalUser.clone(db),
		MenuGroup:         q.MenuGroup.clone(db),
		MenuLink:          q.MenuLink.clone(db),
		OAuth2User:        q.OAuth2User.clone(db),
		Post:              q.Post.clone(db),
		PostTagRelation:   q.PostTagRelation.clone(db),
		PubSocialMedia:    q.PubSocialMedia.clone(db),
		S3ObjectRecord:    q.S3ObjectRecord.clone(db),
		ServiceConf:       q.ServiceConf.clone(db),
		Sqlite3KeywordDoc: q.Sqlite3KeywordDoc.clone(db),
		Tag:               q.Tag.clone(db),
		UVStatistic:       q.UVStatistic.clone(db),
		UserLoginSession:  q.UserLoginSession.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:                db,
		BlockIPRecord:     q.BlockIPRecord.replaceDB(db),
		Category:          q.Category.replaceDB(db),
		CloudFn:           q.CloudFn.replaceDB(db),
		Comment:           q.Comment.replaceDB(db),
		CronJob:           q.CronJob.replaceDB(db),
		FileRecord:        q.FileRecord.replaceDB(db),
		FriendLink:        q.FriendLink.replaceDB(db),
		LocalUser:         q.LocalUser.replaceDB(db),
		MenuGroup:         q.MenuGroup.replaceDB(db),
		MenuLink:          q.MenuLink.replaceDB(db),
		OAuth2User:        q.OAuth2User.replaceDB(db),
		Post:              q.Post.replaceDB(db),
		PostTagRelation:   q.PostTagRelation.replaceDB(db),
		PubSocialMedia:    q.PubSocialMedia.replaceDB(db),
		S3ObjectRecord:    q.S3ObjectRecord.replaceDB(db),
		ServiceConf:       q.ServiceConf.replaceDB(db),
		Sqlite3KeywordDoc: q.Sqlite3KeywordDoc.replaceDB(db),
		Tag:               q.Tag.replaceDB(db),
		UVStatistic:       q.UVStatistic.replaceDB(db),
		UserLoginSession:  q.UserLoginSession.replaceDB(db),
	}
}

type queryCtx struct {
	BlockIPRecord     IBlockIPRecordDo
	Category          ICategoryDo
	CloudFn           ICloudFnDo
	Comment           ICommentDo
	CronJob           ICronJobDo
	FileRecord        IFileRecordDo
	FriendLink        IFriendLinkDo
	LocalUser         ILocalUserDo
	MenuGroup         IMenuGroupDo
	MenuLink          IMenuLinkDo
	OAuth2User        IOAuth2UserDo
	Post              IPostDo
	PostTagRelation   IPostTagRelationDo
	PubSocialMedia    IPubSocialMediaDo
	S3ObjectRecord    IS3ObjectRecordDo
	ServiceConf       IServiceConfDo
	Sqlite3KeywordDoc ISqlite3KeywordDocDo
	Tag               ITagDo
	UVStatistic       IUVStatisticDo
	UserLoginSession  IUserLoginSessionDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		BlockIPRecord:     q.BlockIPRecord.WithContext(ctx),
		Category:          q.Category.WithContext(ctx),
		CloudFn:           q.CloudFn.WithContext(ctx),
		Comment:           q.Comment.WithContext(ctx),
		CronJob:           q.CronJob.WithContext(ctx),
		FileRecord:        q.FileRecord.WithContext(ctx),
		FriendLink:        q.FriendLink.WithContext(ctx),
		LocalUser:         q.LocalUser.WithContext(ctx),
		MenuGroup:         q.MenuGroup.WithContext(ctx),
		MenuLink:          q.MenuLink.WithContext(ctx),
		OAuth2User:        q.OAuth2User.WithContext(ctx),
		Post:              q.Post.WithContext(ctx),
		PostTagRelation:   q.PostTagRelation.WithContext(ctx),
		PubSocialMedia:    q.PubSocialMedia.WithContext(ctx),
		S3ObjectRecord:    q.S3ObjectRecord.WithContext(ctx),
		ServiceConf:       q.ServiceConf.WithContext(ctx),
		Sqlite3KeywordDoc: q.Sqlite3KeywordDoc.WithContext(ctx),
		Tag:               q.Tag.WithContext(ctx),
		UVStatistic:       q.UVStatistic.WithContext(ctx),
		UserLoginSession:  q.UserLoginSession.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
