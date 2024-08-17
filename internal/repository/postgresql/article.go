package postgresql

import (
	"app/pkg/domain/entity"
	"database/sql"
)

type ArticleRepo struct {
	db *sql.DB
}

func NewArticleRepo(db *sql.DB) *ArticleRepo {
	return &ArticleRepo{
		db: db,
	}
}

func (repo ArticleRepo) Create(article *entity.Article) error {
	panic("IMPLEMENT ME")
}

func (repo ArticleRepo) Update(article *entity.Article) error {
	panic("IMPLEMENT ME")
}

func (repo ArticleRepo) Delete(articleId string) error {
	panic("IMPLEMENT ME")
}

func (repo ArticleRepo) GetById(articleId string) (*entity.Article, error) {
	panic("IMPLEMENT ME")
}

func (repo ArticleRepo) GetList(listType string) (*[]entity.Article, error) {
	panic("IMPLEMENT ME")
}

func (repo ArticleRepo) ChangeStatus(articleId string) error {
	panic("IMPLEMENT ME")
}

func (repo ArticleRepo) CreateComment(comment *entity.ArticleComment) error {
	panic("IMPLEMENT ME")
}

func (repo ArticleRepo) GetCommentList(articleId string) (*[]entity.ArticleComment, error) {
	panic("IMPLEMENT ME")
}

func (repo ArticleRepo) UpdateComment(comment *entity.ArticleComment) error {
	panic("IMPLEMENT ME")
}

func (repo ArticleRepo) DeleteComment(commentId string) error {
	panic("IMPLEMENT ME")
}
