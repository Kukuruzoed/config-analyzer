package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Kukuruzoed/config-analyzer/internal/analyzer"
	"github.com/Kukuruzoed/config-analyzer/internal/parser"
	"github.com/Kukuruzoed/config-analyzer/internal/rules"
	pb "github.com/Kukuruzoed/config-analyzer/proto"
)

type Server struct {
	pb.UnimplementedAnalyzerServiceServer
	analyzer *analyzer.Analyzer
	port     int
}

func New(a *analyzer.Analyzer, port int) *Server {
	return &Server{analyzer: a, port: port}
}

func (s *Server) Analyze(ctx context.Context, req *pb.AnalyzeRequest) (*pb.AnalyzeResponse, error) {
	if req.Content == "" {
		return nil, status.Error(codes.InvalidArgument, "content обязателен")
	}

	config, err := parser.ParseString(req.Content)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ошибка парсинга: %v", err)
	}

	issues := s.analyzer.Analyze(config)

	return &pb.AnalyzeResponse{
		Issues: toProto(issues),
		Total:  int32(len(issues)),
	}, nil
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("не удалось занять порт %d: %w", s.port, err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAnalyzerServiceServer(grpcServer, s)

	log.Printf("gRPC сервер запущен на :%d", s.port)
	return grpcServer.Serve(lis)
}

func toProto(issues []rules.Issue) []*pb.Issue {
	result := make([]*pb.Issue, len(issues))
	for i, issue := range issues {
		result[i] = &pb.Issue{
			Severity:       string(issue.Severity),
			Description:    issue.Description,
			Recommendation: issue.Recommendation,
		}
	}
	return result
}
