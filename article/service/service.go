package service

import (
	"context"

	"github.com/k88t76/GraphQL-gRPC-demo/article/pb"
	"github.com/k88t76/GraphQL-gRPC-demo/article/repository"
)

type Service interface {
	CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleResponse, error)
	ReadArticle(ctx context.Context, req *pb.ReadArticleRequest) (*pb.ReadArticleResponse, error)
	UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*pb.UpdateArticleResponse, error)
	DeleteArticle(ctx context.Context, req *pb.DeleteArticleRequest) (*pb.DeleteArticleResponse, error)
	ListArticle(req *pb.ListArticleRequest, stream pb.ArticleService_ListArticleServer) error
}

type service struct {
	repository repository.Repository
}

func NewService(r repository.Repository) Service {
	return &service{r}
}

func (s *service) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleResponse, error) {
	// INSERTする広告のInput(DisplayAppName, IconSrc, CvCondition, BasePoint)を取得
	input := req.GetArticleInput()

	// 広告をDBにINSERTし、INSERTした広告のIDを返す
	id, err := s.repository.InsertArticle(ctx, input)
	if err != nil {
		return nil, err
	}

	// INSERTした広告をレスポンスとして返す
	return &pb.CreateArticleResponse{
		Article: &pb.Article{
			Id:             id,
			DisplayAppName: input.DisplayAppName,
			IconSrc:        input.IconSrc,
			CvCondition:    input.CvCondition,
			BasePoint:      input.BasePoint,
		},
	}, nil
}

func (s *service) ReadArticle(ctx context.Context, req *pb.ReadArticleRequest) (*pb.ReadArticleResponse, error) {
	// READする広告のIDを取得
	id := req.GetId()

	// DBから該当IDの広告を取得
	a, err := s.repository.SelectArticleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	//　取得した広告をレスポンスとして返す
	return &pb.ReadArticleResponse{
		Article: &pb.Article{
			Id:             id,
			DisplayAppName: a.DisplayAppName,
			IconSrc:        a.IconSrc,
			CvCondition:    a.CvCondition,
			BasePoint:      a.BasePoint,
		},
	}, nil

}

func (s *service) UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*pb.UpdateArticleResponse, error) {
	// UPDATEする広告のIDを取得
	id := req.GetId()

	// UPDATEする広告の変更内容(DisplayAppName, IconSrc, CvCondition)を取得
	input := req.GetArticleInput()

	//　該当IDの広告をUPDATE
	if err := s.repository.UpdateArticle(ctx, id, input); err != nil {
		return nil, err
	}

	// UPDATEした広告をレスポンスとして返す
	return &pb.UpdateArticleResponse{
		Article: &pb.Article{
			Id:             id,
			DisplayAppName: input.DisplayAppName,
			IconSrc:        input.IconSrc,
			CvCondition:    input.CvCondition,
			BasePoint:      input.BasePoint,
		},
	}, nil
}

func (s *service) DeleteArticle(ctx context.Context, req *pb.DeleteArticleRequest) (*pb.DeleteArticleResponse, error) {
	// DELETEする広告のIDを取得
	id := req.GetId()

	// 該当IDの広告をDELETE
	if err := s.repository.DeleteArticle(ctx, id); err != nil {
		return nil, err
	}

	// DELETEした広告のIDをレスポンスとして返す
	return &pb.DeleteArticleResponse{Id: id}, nil
}

func (s *service) ListArticle(req *pb.ListArticleRequest, stream pb.ArticleService_ListArticleServer) error {
	// 広告を全取得
	rows, err := s.repository.SelectAllArticles()
	if err != nil {
		return err
	}
	defer rows.Close()

	// 取得した広告を１つ１つレスポンスとしてServer Streamingで返す
	for rows.Next() {
		var a pb.Article
		err := rows.Scan(&a.Id, &a.DisplayAppName, &a.IconSrc, &a.CvCondition, &a.BasePoint)
		if err != nil {
			return err
		}
		stream.Send(&pb.ListArticleResponse{Article: &a})
	}
	return nil
}
