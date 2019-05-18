package user_data

import (
	"log"

	pb "../../api_struct"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

type GrpcProvider struct {
	ServerAddress string
}

func (provider GrpcProvider) GetUsers(ids []string) ([]*User, error) {
	conn, err := grpc.Dial(provider.ServerAddress, grpc.WithInsecure())

	if err != nil {
		log.Printf("Can't connect to  server: %v", err)
		return nil, err
	}
	defer conn.Close()
	c := pb.NewAPIClient(conn)

	response, err := c.GetUsers(context.Background(), &pb.GetUsersRequest{
		Uid: ids,
	})

	if err != nil {
		log.Printf("Can't get users: %v", err)
		return nil, err
	}
	var users []*User
	for _, data := range response.Data {
		user := User{
			Uid:    data.Uid,
			Login:  data.Login,
			Avatar: data.Avatar,
		}
		users = append(users, &user)
	}
	return users, nil
}
