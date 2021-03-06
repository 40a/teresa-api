package team

import (
	context "golang.org/x/net/context"

	"github.com/luizalabs/teresa-api/models/storage"
	teampb "github.com/luizalabs/teresa-api/pkg/protobuf/team"
	"github.com/luizalabs/teresa-api/pkg/server/auth"
	"google.golang.org/grpc"
)

type Service struct {
	ops Operations
}

func (s *Service) Create(ctx context.Context, request *teampb.CreateRequest) (*teampb.Empty, error) {
	u := ctx.Value("user").(*storage.User)
	if !u.IsAdmin {
		return nil, auth.ErrPermissionDenied
	}
	if err := s.ops.Create(request.Name, request.Email, request.Url); err != nil {
		return nil, err
	}
	return &teampb.Empty{}, nil
}

func (s *Service) AddUser(ctx context.Context, request *teampb.AddUserRequest) (*teampb.Empty, error) {
	u := ctx.Value("user").(*storage.User)
	if !u.IsAdmin {
		return nil, auth.ErrPermissionDenied
	}
	if err := s.ops.AddUser(request.Name, request.User); err != nil {
		return nil, err
	}
	return &teampb.Empty{}, nil
}

func (s *Service) List(ctx context.Context, _ *teampb.Empty) (*teampb.ListResponse, error) {
	var (
		teams []*storage.Team
		err   error
	)

	u := ctx.Value("user").(*storage.User)
	if u.IsAdmin {
		teams, err = s.ops.List()
	} else {
		teams, err = s.ops.ListByUser(u.Email)
	}

	if err != nil {
		return nil, err
	}

	resp := &teampb.ListResponse{}
	for _, t := range teams {
		currentTeam := &teampb.ListResponse_Team{Name: t.Name, Email: t.Email, Url: t.URL}
		for _, user := range t.Users {
			currentUser := &teampb.ListResponse_User{Name: user.Name, Email: user.Email}
			currentTeam.Users = append(currentTeam.Users, currentUser)
		}
		resp.Teams = append(resp.Teams, currentTeam)
	}

	return resp, nil
}

func (s *Service) RegisterService(grpcServer *grpc.Server) {
	teampb.RegisterTeamServer(grpcServer, s)
}

func NewService(ops Operations) *Service {
	return &Service{ops: ops}
}
