package repository

import (
	"context"
	"database/sql"

	"github.com/k88t76/GraphQL-gRPC-demo/article/pb"
	_ "github.com/mattn/go-sqlite3"
)

type Repository interface {
	InsertArticle(ctx context.Context, input *pb.ArticleInput) (int64, error)
	SelectArticleByID(ctx context.Context, id int64) (*pb.Article, error)
	UpdateArticle(ctx context.Context, id int64, input *pb.ArticleInput) error
	DeleteArticle(ctx context.Context, id int64) error
	SelectAllArticles() (*sql.Rows, error)
}

type sqliteRepo struct {
	db *sql.DB
}

func NewsqliteRepo() (Repository, error) {
	db, err := sql.Open("sqlite3", "./article/article.sql")
	if err != nil {
		return nil, err
	}

	// articlesテーブルを作成
	cmd := `CREATE TABLE IF NOT EXISTS articles(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		display_app_name STRING,
		icon_src STRING,
		cv_condition STRING,
		base_point INTEGER)`

	_, err = db.Exec(cmd)
	if err != nil {
		return nil, err
	}
	return &sqliteRepo{db}, nil
}

func (r *sqliteRepo) InsertArticle(ctx context.Context, input *pb.ArticleInput) (int64, error) {
	// Inputの内容(DisplayAppName,IconSrc,CvCondition,BasePoint)をarticlesテーブルにINSERT
	cmd := "INSERT INTO articles(display_app_name, icon_src, cv_condition, base_point) VALUES (?, ?, ?, ?)"
	result, err := r.db.Exec(cmd, input.DisplayAppName, input.IconSrc, input.CvCondition, input.BasePoint)
	if err != nil {
		return 0, err
	}

	// INSERTした広告のIDを取得
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// INSERTした広告のIDを返す
	return id, nil
}

func (r *sqliteRepo) SelectArticleByID(ctx context.Context, id int64) (*pb.Article, error) {
	// 該当IDの広告をSELECT
	cmd := "SELECT * FROM articles WHERE id = ?"
	row := r.db.QueryRow(cmd, id)
	var a pb.Article

	// SELECTした広告の内容を読み取る
	err := row.Scan(&a.Id, &a.DisplayAppName, &a.IconSrc, &a.CvCondition, &a.BasePoint)
	if err != nil {
		return nil, err
	}

	// SELECTした広告を返す
	return &pb.Article{
		Id:             a.Id,
		DisplayAppName: a.DisplayAppName,
		IconSrc:        a.IconSrc,
		CvCondition:    a.CvCondition,
		BasePoint:      a.BasePoint,
	}, nil
}

func (r *sqliteRepo) UpdateArticle(ctx context.Context, id int64, input *pb.ArticleInput) error {
	// 該当IDのDisplayAppName, IconSrc, CvCondition, BasePointをUPDATE
	cmd := "UPDATE articles SET display_app_name = ?, icon_src = ?, cv_condition = ?, base_point = ? WHERE id = ?"
	_, err := r.db.Exec(cmd, input.DisplayAppName, input.IconSrc, input.CvCondition, input.BasePoint, id)
	if err != nil {
		return err
	}

	// errorがなければ返り値なし
	return nil
}

func (r *sqliteRepo) DeleteArticle(ctx context.Context, id int64) error {
	// 該当IDの広告をDELETE
	cmd := "DELETE FROM articles WHERE id = ?"
	_, err := r.db.Exec(cmd, id)
	if err != nil {
		return err
	}

	// errorがなければ返り値なし
	return nil
}

func (r *sqliteRepo) SelectAllArticles() (*sql.Rows, error) {
	// articlesテーブルの広告を全取得
	cmd := "SELECT * FROM articles"
	rows, err := r.db.Query(cmd)
	if err != nil {
		return nil, err
	}

	// 全取得した広告を*sql.Rowsの形で返す
	return rows, nil
}
