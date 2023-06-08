package communication

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/constraints"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/msik-404/micro-appoint-companies/internal/database"
	"github.com/msik-404/micro-appoint-companies/internal/models"
)

type Server struct {
	UnimplementedApiServer
	Client mongo.Client
}

func verifyString(value *string, maxLength int) (*string, error) {
	if value != nil {
		if len(*value) > int(maxLength) {
			return nil, status.Errorf(
				codes.InvalidArgument,
				"value should be shorter than %d",
				maxLength,
			)
		}
	}
	return value, nil
}

func verifyInteger[T constraints.Integer](value *T, low T, high T) (*T, error) {
	if value != nil {
		if *value > high || *value <= low {
			return nil, status.Errorf(
				codes.InvalidArgument,
				"value should be smaller than %d and greater than %d",
				high,
				low,
			)
		}
	}
	return value, nil
}

func (s *Server) AddService(
	ctx context.Context,
	request *AddServiceRequest,
) (*emptypb.Empty, error) {
	companyID, err := primitive.ObjectIDFromHex(request.CompanyId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	name, err := verifyString(request.Name, 30)
	if err != nil {
		return nil, err
	}
	price, err := verifyInteger(request.Price, 0, 1000000)
	if err != nil {
		return nil, err
	}
	duration, err := verifyInteger(request.Duration, 0, 480)
	if err != nil {
		return nil, err
	}
	description, err := verifyString(request.Description, 300)
	if err != nil {
		return nil, err
	}
	// check for nil
	newSerivce := models.Service{
		Name:        name,
		Price:       price,
		Duration:    duration,
		Description: description,
	}
	db := s.Client.Database(database.DBName)
	result, err := newSerivce.InsertOne(ctx, db, companyID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if result.MatchedCount == 0 {
		return nil, status.Error(
			codes.NotFound,
			"Company with that id was not found",
		)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateService(
	ctx context.Context,
	request *UpdateServiceRequest,
) (*emptypb.Empty, error) {
	serviceID, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	name, err := verifyString(request.Name, 30)
	if err != nil {
		return nil, err
	}
	price, err := verifyInteger(request.Price, 0, 1000000)
	if err != nil {
		return nil, err
	}
	duration, err := verifyInteger(request.Duration, 0, 480)
	if err != nil {
		return nil, err
	}
	description, err := verifyString(request.Description, 300)
	if err != nil {
		return nil, err
	}
	serviceUpdate := models.Service{
		Name:        name,
		Price:       price,
		Duration:    duration,
		Description: description,
	}
	db := s.Client.Database(database.DBName)
	result, err := serviceUpdate.UpdateOne(ctx, db, serviceID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if result.MatchedCount == 0 {
		return nil, status.Error(
			codes.NotFound,
			"Service with that id was not found",
		)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteService(
	ctx context.Context,
	request *DeleteServiceRequest,
) (*emptypb.Empty, error) {
	serviceID, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	db := s.Client.Database(database.DBName)
	result, err := models.DeleteOneService(ctx, db, serviceID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if result.MatchedCount == 0 {
		return nil, status.Error(
			codes.NotFound,
			"Service with that id was not found",
		)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) FindManyServices(
	ctx context.Context,
	request *ManyServicesRequest,
) (*ManyServicesReply, error) {
	companyID, err := primitive.ObjectIDFromHex(request.CompanyId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	startValue := primitive.NilObjectID
	if request.StartValue != nil {
		startValue, err = primitive.ObjectIDFromHex(*request.StartValue)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	var nPerPage int64 = 10
	if request.NPerPage != nil {
		nPerPage = *request.NPerPage
	}
	db := s.Client.Database(database.DBName)
	cursor, err := models.FindManyServices(ctx, db, companyID, startValue, nPerPage)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	reply := &ManyServicesReply{}
	for cursor.Next(ctx) {
		var serviceModel = models.Service{}
		if err := cursor.Decode(&serviceModel); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		serviceProto := &Service{
			Id:          serviceModel.ID.Hex(),
			Name:        serviceModel.Name,
			Price:       serviceModel.Price,
			Duration:    serviceModel.Duration,
			Description: serviceModel.Description,
		}
		reply.Services = append(reply.Services, serviceProto)
	}
	if len(reply.Services) == 0 {
		return nil, status.Error(
			codes.NotFound,
			"This company does not have any services",
		)
	}
	return reply, nil
}

func (s *Server) AddCompany(
	ctx context.Context,
	request *AddCompanyRequest,
) (*emptypb.Empty, error) {
	name, err := verifyString(request.Name, 30)
	if err != nil {
		return nil, err
	}
	if name == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"name should be set",
		)

	}
	companyType, err := verifyString(request.Type, 30)
	if err != nil {
		return nil, err
	}
	localisation, err := verifyString(request.Localisation, 60)
	if err != nil {
		return nil, err
	}
	shortDescription, err := verifyString(request.ShortDescription, 150)
	if err != nil {
		return nil, err
	}
	longDescription, err := verifyString(request.LongDescription, 300)
	if err != nil {
		return nil, err
	}
	newCompany := models.Company{
		Name:             name,
		Type:             companyType,
		Localisation:     localisation,
		ShortDescription: shortDescription,
		LongDescription:  longDescription,
	}
	db := s.Client.Database(database.DBName)
	_, err = newCompany.InsertOne(ctx, db)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateCompany(
	ctx context.Context,
	request *UpdateCompanyRequest,
) (*emptypb.Empty, error) {
	companyID, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	name, err := verifyString(request.Name, 30)
	if err != nil {
		return nil, err
	}
	companyType, err := verifyString(request.Type, 30)
	if err != nil {
		return nil, err
	}
	localisation, err := verifyString(request.Localisation, 60)
	if err != nil {
		return nil, err
	}
	shortDescription, err := verifyString(request.ShortDescription, 150)
	if err != nil {
		return nil, err
	}
	longDescription, err := verifyString(request.LongDescription, 300)
	if err != nil {
		return nil, err
	}
	companyUpdate := models.Company{
		Name:             name,
		Type:             companyType,
		Localisation:     localisation,
		ShortDescription: shortDescription,
		LongDescription:  longDescription,
	}
	db := s.Client.Database(database.DBName)
	result, err := companyUpdate.UpdateOne(ctx, db, companyID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if result.MatchedCount == 0 {
		return nil, status.Error(
			codes.NotFound,
			"Company with that id was not found",
		)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteCompany(
	ctx context.Context,
	request *DeleteCompanyRequest,
) (*emptypb.Empty, error) {
	companyID, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	db := s.Client.Database(database.DBName)
	result, err := models.DeleteOneCompany(ctx, db, companyID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if result.DeletedCount == 0 {
		return nil, status.Error(
			codes.NotFound,
			"Company with that id was not found",
		)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) FindOneCompany(
	ctx context.Context,
	request *CompanyRequest,
) (*CompanyReply, error) {
	companyID, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	db := s.Client.Database(database.DBName)
	companyModel := models.Company{}
	err = models.FindOneCompany(ctx, db, companyID).Decode(&companyModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	companyProto := &CompanyReply{
		Name:            companyModel.Name,
		Type:            companyModel.Type,
		Localisation:    companyModel.Localisation,
		LongDescription: companyModel.LongDescription,
	}
	for _, service := range companyModel.Services {
		// Mandatory copy, to get proper pointers
		service := service
		serviceProto := Service{
			Id:          service.ID.Hex(),
			Name:        service.Name,
			Price:       service.Price,
			Duration:    service.Duration,
			Description: service.Description,
		}
		companyProto.Services = append(companyProto.Services, &serviceProto)
	}
	return companyProto, nil
}

func (s *Server) FindManyCompanies(
	ctx context.Context,
	request *CompaniesRequest,
) (reply *CompaniesReply, err error) {
	startValue := primitive.NilObjectID
	if request.StartValue != nil {
		startValue, err = primitive.ObjectIDFromHex(*request.StartValue)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	var nPerPage int64 = 10
	if request.NPerPage != nil {
		nPerPage = *request.NPerPage
	}
	db := s.Client.Database(database.DBName)
	cursor, err := models.FindManyCompanies(ctx, db, startValue, nPerPage)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	reply = &CompaniesReply{}
	for cursor.Next(ctx) {
		var companyModel models.Company
		if err := cursor.Decode(&companyModel); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		companyProto := &CompanyShort{
			Id:               companyModel.ID.Hex(),
			Name:             companyModel.Name,
			Type:             companyModel.Type,
			Localisation:     companyModel.Localisation,
			ShortDescription: companyModel.ShortDescription,
		}
		reply.Companies = append(reply.Companies, companyProto)
	}
	if len(reply.Companies) == 0 {
		return nil, status.Error(
			codes.NotFound,
			"There aren't any companies",
		)
	}
	return reply, nil
}
