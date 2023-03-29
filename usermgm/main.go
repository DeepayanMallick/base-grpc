package main

import (
	"fmt"
	"log"
	"net"

	"github.com/DeepayanMallick/base-grpc/gunk/v1/usermgm/user"
	"github.com/DeepayanMallick/base-grpc/usermgm/storage/postgres"
	"github.com/DeepayanMallick/base-grpc/utility"
	"github.com/DeepayanMallick/base-grpc/utility/logging"
	"github.com/DeepayanMallick/base-grpc/utility/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	userC "github.com/DeepayanMallick/base-grpc/usermgm/core/user"
	userS "github.com/DeepayanMallick/base-grpc/usermgm/service/user"
)

var (
	svcName = "usermgm"
	version = "development"
)

func main() {
	cfg, err := utility.NewConfig("env/config")
	if err != nil {
		log.Fatal(err)
	}

	logger := logging.NewLogger(cfg).WithFields(logrus.Fields{
		"service": svcName,
		"version": version,
	})

	dbString := utility.NewDBString(cfg)
	db, err := postgres.NewStorage(dbString, logger)
	if err != nil {
		logger.WithError(err).Error("unable to connect DB")
		return
	}

	if err := db.RunMigration(cfg.GetString("database.migrationDir")); err != nil {
		logger.WithError(err).Error("unable to run DB migrations")
		return
	}

	grpcServer, err := setupGRPCServer(db, cfg, logger)
	if err != nil {
		logger.WithError(err).Error("unable to setup grpc service")
		return
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GetInt("server.port")))
	if err != nil {
		logger.WithError(err).Error("unable to listen port")
		return
	}

	log.Printf("server %s listening at: %+v", svcName, lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		logger.WithError(err).Error("unable to serve the GRPC server")
	}

	log.Println("server stoped")
}

func setupGRPCServer(store *postgres.Storage, config *viper.Viper, logger *logrus.Entry) (*grpc.Server, error) {
	//	Authentication
	/* userID from cms */
	// 	Authorization
	mw := middleware.New(
		config.GetString("runtime.environment"),
		logger,
		middleware.Config{},
	)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(mw),
	)

	coreUsr := userC.New(store, logger)
	userSvc := userS.New(coreUsr, logger)
	user.RegisterUserServiceServer(grpcServer, userSvc)

	return grpcServer, nil
}
