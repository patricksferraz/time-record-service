package grpc_test

// type mockTimeRecordService_ListTimeRecordsServer struct {
// 	gogrpc.ServerStream
// 	Results []*pb.TimeRecord
// }

// func (_m *mockTimeRecordService_ListTimeRecordsServer) Send(timeRecord *pb.TimeRecord) error {
// 	_m.Results = append(_m.Results, timeRecord)
// 	return nil
// }

// func TestGrpc_Register(t *testing.T) {

// 	interceptor := &grpc.AuthInterceptor{
// 		EmployeeClaims: &model.EmployeeClaims{
// 			ID: uuid.NewV4().String(),
// 		},
// 	}

// 	ctx := new(context.Context)
// 	uri := utils.GetEnv("DB_URI", "mongodb://localhost")
// 	dbName := utils.GetEnv("DB_NAME", "test")
// 	db, err := db.NewMongo(*ctx, uri, dbName)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close(*ctx)
// 	defer db.Database.Drop(*ctx)

// 	timeRecordRepository := repository.NewTimeRecordRepository(db)
// 	timeRecordService := service.NewTimeRecordService(timeRecordRepository)
// 	timeRecordGrpcService := grpc.NewTimeRecordGrpcService(timeRecordService, interceptor)

// 	reqRegister := &pb.RegisterTimeRecordRequest{
// 		Time:        timestamppb.Now(),
// 		Description: faker.Lorem().Sentence(10),
// 	}

// 	tr, err := timeRecordGrpcService.RegisterTimeRecord(*ctx, reqRegister)
// 	require.NotEmpty(t, uuid.FromStringOrNil(tr.Id))
// 	require.NotEmpty(t, uuid.FromStringOrNil(tr.EmployeeId))
// 	require.Equal(t, tr.Time, reqRegister.Time)
// 	require.Equal(t, tr.Description, reqRegister.Description)
// 	require.NotEmpty(t, tr.RegularTime)
// 	require.NotEmpty(t, tr.Status)
// 	require.Nil(t, err)

// 	interceptor.EmployeeClaims.ID = ""
// 	_, err = timeRecordGrpcService.RegisterTimeRecord(*ctx, reqRegister)
// 	require.NotNil(t, err)
// }

// func TestGrpc_Approve(t *testing.T) {

// 	interceptor := &grpc.AuthInterceptor{
// 		EmployeeClaims: &model.EmployeeClaims{
// 			ID: uuid.NewV4().String(),
// 		},
// 	}

// 	ctx := new(context.Context)
// 	uri := utils.GetEnv("DB_URI", "mongodb://localhost")
// 	dbName := utils.GetEnv("DB_NAME", "test")
// 	db, err := db.NewMongo(*ctx, uri, dbName)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close(*ctx)
// 	defer db.Database.Drop(*ctx)

// 	timeRecordRepository := repository.NewTimeRecordRepository(db)
// 	timeRecordService := service.NewTimeRecordService(timeRecordRepository)
// 	timeRecordGrpcService := grpc.NewTimeRecordGrpcService(timeRecordService, interceptor)

// 	reqRegister := &pb.RegisterTimeRecordRequest{
// 		Time:        timestamppb.New(time.Now().AddDate(0, 0, -1)),
// 		Description: faker.Lorem().Sentence(10),
// 	}
// 	tr, _ := timeRecordGrpcService.RegisterTimeRecord(*ctx, reqRegister)

// 	reqApprove := &pb.ApproveTimeRecordRequest{Id: tr.Id}
// 	_, err = timeRecordGrpcService.ApproveTimeRecord(*ctx, reqApprove)
// 	require.NotNil(t, err)

// 	interceptor.EmployeeClaims.ID = uuid.NewV4().String()
// 	resApprove, err := timeRecordGrpcService.ApproveTimeRecord(*ctx, reqApprove)
// 	require.NotEmpty(t, resApprove.Status)
// 	require.Nil(t, err)
// }

// func TestGrpc_Find(t *testing.T) {

// 	interceptor := &grpc.AuthInterceptor{
// 		EmployeeClaims: &model.EmployeeClaims{
// 			ID: uuid.NewV4().String(),
// 		},
// 	}

// 	ctx := new(context.Context)
// 	uri := utils.GetEnv("DB_URI", "mongodb://localhost")
// 	dbName := utils.GetEnv("DB_NAME", "test")
// 	db, err := db.NewMongo(*ctx, uri, dbName)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close(*ctx)
// 	defer db.Database.Drop(*ctx)

// 	timeRecordRepository := repository.NewTimeRecordRepository(db)
// 	timeRecordService := service.NewTimeRecordService(timeRecordRepository)
// 	timeRecordGrpcService := grpc.NewTimeRecordGrpcService(timeRecordService, interceptor)

// 	reqRegister := &pb.RegisterTimeRecordRequest{
// 		Time:        timestamppb.Now(),
// 		Description: faker.Lorem().Sentence(10),
// 	}
// 	tr, _ := timeRecordGrpcService.RegisterTimeRecord(*ctx, reqRegister)

// 	findRequest := &pb.FindTimeRecordRequest{
// 		Id: tr.Id,
// 	}

// 	resFind, err := timeRecordGrpcService.FindTimeRecord(*ctx, findRequest)
// 	require.Equal(t, tr.Id, resFind.Id)
// 	require.Equal(t, tr.EmployeeId, resFind.EmployeeId)
// 	require.Equal(t, tr.Time.Seconds, resFind.Time.Seconds)
// 	require.Equal(t, tr.Description, resFind.Description)
// 	require.Equal(t, tr.RegularTime, resFind.RegularTime)
// 	require.Equal(t, tr.Status, resFind.Status)
// 	require.Equal(t, tr.ApprovedBy, resFind.ApprovedBy)
// 	require.Nil(t, err)

// 	findRequest.Id = ""
// 	_, err = timeRecordGrpcService.FindTimeRecord(*ctx, findRequest)
// 	require.NotNil(t, err)
// }

// FIXME: runtime error in Mock context
// func TestGrpc_FindAllByEmployeeID(t *testing.T) {

// 	id := uuid.NewV4().String()
// 	interceptor := &grpc.AuthInterceptor{
// 		EmployeeClaims: &model.EmployeeClaims{
// 			ID: id,
// 		},
// 	}

// 	db, err := db.ConnectMongoDB()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Session.Close()
// 	defer db.DropDatabase()

// 	timeRecordRepository := repository.NewTimeRecordRepository(db)
// 	timeRecordService := service.NewTimeRecordService(timeRecordRepository)
// 	timeRecordGrpcService := grpc.NewTimeRecordGrpcService(timeRecordService, interceptor)

// 	mock := &mockTimeRecordService_ListTimeRecordsServer{}
// 	reqRegister := &pb.RegisterTimeRecordRequest{
// 		Time:        timestamppb.Now(),
// 		Description: faker.Lorem().Sentence(10),
// 	}
// 	timeRecordGrpcService.RegisterTimeRecord(mock.Context(), reqRegister)

// reqList := &pb.ListTimeRecordsRequest{
// 	EmployeeId: id,
// 	FromDate:   timestamppb.New(time.Now().AddDate(0, 0, -1)),
// 	ToDate:     timestamppb.New(time.Now()),
// }

// timeRecordGrpcService.ListTimeRecords(reqList, mock)
// log.Println(mock)
// resFind := mock.Results[0]
// require.Equal(t, 1, len(mock.Results))
// require.Equal(t, tr.Id, resFind.Id)
// require.Equal(t, tr.EmployeeId, resFind.EmployeeId)
// require.Equal(t, tr.Time.Seconds, resFind.Time.Seconds)
// require.Equal(t, tr.Description, resFind.Description)
// require.Equal(t, tr.RegularTime, resFind.RegularTime)
// require.Equal(t, tr.Status, resFind.Status)
// require.Equal(t, tr.ApprovedBy, resFind.ApprovedBy)
// require.Nil(t, err)

// listTimeRecordsRequest.EmployeeId = ""
// err = timeRecordGrpcService.ListTimeRecords(listTimeRecordsRequest, mock)
// require.NotNil(t, err)
// }
