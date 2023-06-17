package companiespb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (s *Server) AddService(
	ctx context.Context,
	request *AddServiceRequest,
) (*emptypb.Empty, error) {
	companyID, err := primitive.ObjectIDFromHex(request.GetCompanyId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = verifyString(request.Name, 30)
	if err != nil {
		return nil, err
	}
	err = verifyInteger(request.Price, 0, 1000000)
	if err != nil {
		return nil, err
	}
	err = verifyInteger(request.Duration, 0, 480)
	if err != nil {
		return nil, err
	}
	err = verifyString(request.Description, 300)
	if err != nil {
		return nil, err
	}
	// check for nil
	newSerivce := models.Service{
		Name:        request.GetName(),
		Price:       request.GetPrice(),
		Duration:    request.GetDuration(),
		Description: request.GetDescription(),
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
	serviceID, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	companyID, err := primitive.ObjectIDFromHex(request.GetCompanyId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = verifyString(request.Name, 30)
	if err != nil {
		return nil, err
	}
	err = verifyInteger(request.Price, 0, 1000000)
	if err != nil {
		return nil, err
	}
	err = verifyInteger(request.Duration, 0, 480)
	if err != nil {
		return nil, err
	}
	err = verifyString(request.Description, 300)
	if err != nil {
		return nil, err
	}
	serviceUpdate := models.ServiceUpdate{
		Name:        request.Name,
		Price:       request.Price,
		Duration:    request.Duration,
		Description: request.Description,
	}
	db := s.Client.Database(database.DBName)
	result, err := serviceUpdate.UpdateOne(ctx, db, companyID, serviceID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if result.MatchedCount == 0 {
		return nil, status.Error(
			codes.NotFound,
			"Service with that companyID and serviceID was not found",
		)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteService(
	ctx context.Context,
	request *DeleteServiceRequest,
) (*emptypb.Empty, error) {
	serviceID, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	companyID, err := primitive.ObjectIDFromHex(request.GetCompanyId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	db := s.Client.Database(database.DBName)
	result, err := models.DeleteOneService(ctx, db, companyID, serviceID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if result.MatchedCount == 0 {
		return nil, status.Error(
			codes.NotFound,
			"Service with that companyID and serviceID was not found",
		)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) FindManyServices(
	ctx context.Context,
	request *ServicesRequest,
) (*ServicesReply, error) {
	companyID, err := primitive.ObjectIDFromHex(request.GetCompanyId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	startValue := primitive.NilObjectID
	if request.StartValue != nil {
		startValue, err = primitive.ObjectIDFromHex(request.GetStartValue())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	var nPerPage int64 = 30
	if request.NPerPage != nil {
		nPerPage = request.GetNPerPage()
	}
	db := s.Client.Database(database.DBName)
	cursor, err := models.FindManyServices(ctx, db, companyID, startValue, nPerPage)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
    defer cursor.Close(ctx)
	reply := &ServicesReply{}
	for cursor.Next(ctx) {
		var serviceModel models.Service
		if err := cursor.Decode(&serviceModel); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		serviceID := serviceModel.ID.Hex()
		serviceProto := &Service{
			Id:          &serviceID,
			Name:        &serviceModel.Name,
			Price:       &serviceModel.Price,
			Duration:    &serviceModel.Duration,
			Description: &serviceModel.Description,
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
) (*AddCompanyReply, error) {
	if request.Name == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"name should be set",
		)
	}
	err := verifyString(request.Name, 30)
	if err != nil {
		return nil, err
	}
	err = verifyString(request.Type, 30)
	if err != nil {
		return nil, err
	}
	err = verifyString(request.Localisation, 60)
	if err != nil {
		return nil, err
	}
	err = verifyString(request.ShortDescription, 150)
	if err != nil {
		return nil, err
	}
	err = verifyString(request.LongDescription, 300)
	if err != nil {
		return nil, err
	}
	newCompany := models.Company{
		Name:             request.GetName(),
		Type:             request.GetType(),
		Localisation:     request.GetLocalisation(),
		ShortDescription: request.GetShortDescription(),
		LongDescription:  request.GetLongDescription(),
	}
	db := s.Client.Database(database.DBName)
	result, err := newCompany.InsertOne(ctx, db)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	insertedID := result.InsertedID.(primitive.ObjectID).Hex()
	return &AddCompanyReply{
		Id: &insertedID,
	}, nil
}

func (s *Server) UpdateCompany(
	ctx context.Context,
	request *UpdateCompanyRequest,
) (*emptypb.Empty, error) {
	companyID, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = verifyString(request.Name, 30)
	if err != nil {
		return nil, err
	}
	err = verifyString(request.Type, 30)
	if err != nil {
		return nil, err
	}
	err = verifyString(request.Localisation, 60)
	if err != nil {
		return nil, err
	}
	err = verifyString(request.ShortDescription, 150)
	if err != nil {
		return nil, err
	}
	err = verifyString(request.LongDescription, 300)
	if err != nil {
		return nil, err
	}
	companyUpdate := models.CompanyUpdate{
		Name:             request.Name,
		Type:             request.Type,
		Localisation:     request.Localisation,
		ShortDescription: request.ShortDescription,
		LongDescription:  request.LongDescription,
	}
	db := s.Client.Database(database.DBName)
	result, err := companyUpdate.UpdateOne(ctx, db, companyID)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
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
	companyID, err := primitive.ObjectIDFromHex(request.GetId())
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
	companyID, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	db := s.Client.Database(database.DBName)
	var companyModel models.Company
	err = models.FindOneCompany(ctx, db, companyID).Decode(&companyModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	companyProto := &CompanyReply{
		Name:            &companyModel.Name,
		Type:            &companyModel.Type,
		Localisation:    &companyModel.Localisation,
		LongDescription: &companyModel.LongDescription,
	}
	for idx := range companyModel.Services {
		serviceID := companyModel.Services[idx].ID.Hex()
		serviceProto := Service{
			Id:          &serviceID,
			Name:        &companyModel.Services[idx].Name,
			Price:       &companyModel.Services[idx].Price,
			Duration:    &companyModel.Services[idx].Duration,
			Description: &companyModel.Services[idx].Description,
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
		startValue, err = primitive.ObjectIDFromHex(request.GetStartValue())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	var nPerPage int64 = 30
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
    defer cursor.Close(ctx)
	reply = &CompaniesReply{}
	for cursor.Next(ctx) {
		var companyModel models.Company
		if err := cursor.Decode(&companyModel); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		companyID := companyModel.ID.Hex()
		companyProto := &CompanyShort{
			Id:               &companyID,
			Name:             &companyModel.Name,
			Type:             &companyModel.Type,
			Localisation:     &companyModel.Localisation,
			ShortDescription: &companyModel.ShortDescription,
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

func (s *Server) FindManyCompaniesByIds(
	ctx context.Context,
	request *CompaniesByIdsRequest,
) (reply *CompaniesReply, err error) {
	var companiesIDS []primitive.ObjectID
	for _, hex := range request.GetIds() {
		companyID, err := primitive.ObjectIDFromHex(hex)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		companiesIDS = append(companiesIDS, companyID)
	}
	if len(companiesIDS) == 0 {
		return nil, status.Error(
			codes.InvalidArgument, 
			"At least one id should be provided in the request",
		)
	}
	startValue := primitive.NilObjectID
	if request.StartValue != nil {
		startValue, err = primitive.ObjectIDFromHex(request.GetStartValue())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	var nPerPage int64 = 30
	if request.NPerPage != nil {
		nPerPage = *request.NPerPage
	}

	db := s.Client.Database(database.DBName)
	cursor, err := models.FindManyCompaniesByIds(ctx, db, companiesIDS, startValue, nPerPage)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    defer cursor.Close(ctx)
	reply = &CompaniesReply{}
	for cursor.Next(ctx) {
		var companyModel models.Company
		if err := cursor.Decode(&companyModel); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		companyID := companyModel.ID.Hex()
		companyProto := &CompanyShort{
			Id:               &companyID,
			Name:             &companyModel.Name,
			Type:             &companyModel.Type,
			Localisation:     &companyModel.Localisation,
			ShortDescription: &companyModel.ShortDescription,
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
