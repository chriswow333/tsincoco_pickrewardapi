package application

import (
	"context"
	"fmt"
	"os"
	"path"
	pb "pickrewardapi/internal/application/image/v1/proto/generated"

	"go.uber.org/dig"
	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
)

type server struct {
	dig.In

	pb.UnimplementedImageV1Server
}

func NewImageServer(
	s *grpc.Server,
) {
	log.WithFields(log.Fields{
		"pos": "[image.api][NewImageServer]",
	}).Info("Init")
	pb.RegisterImageV1Server(s, &server{})
}

func (s *server) DownloadImage(ctx context.Context, in *pb.ImageReq) (*pb.ImageReply, error) {
	logPos := "[image.api][DownloadImage]"

	log.WithFields(log.Fields{
		"pos": logPos,
		"req": in,
	}).Info("Request")

	folderPath := path.Join("./assets/images/", in.Type, in.Id)
	fmt.Println(folderPath)

	data, err := os.ReadFile(folderPath)
	if err != nil {
		log.WithFields(log.Fields{
			"pos": logPos,
		}).Error("os.ReadFile failed: ", err)

		return &pb.ImageReply{
			Reply: &pb.Reply{
				Status: 1,
				Error: &pb.Error{
					ErrorCode:    100,
					ErrorMessage: "os.ReadFile failed",
				},
			},
		}, nil
	}

	log.WithFields(log.Fields{
		"pos":  logPos,
		"resp": "downloaded",
	}).Info("Response")

	return &pb.ImageReply{
		Reply: &pb.Reply{
			Status: 0,
		},
		Image: &pb.ImageReply_Image{
			Data: data,
		},
	}, nil
}
