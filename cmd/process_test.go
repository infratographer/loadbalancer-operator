package cmd

// func (suite *cmdTestSuite) TestNewKubeAuth() {
// 	suite.T().Parallel()

// 	type testCase struct {
// 		name        string
// 		path        string
// 		config      *rest.Config
// 		expectError bool
// 	}

// 	pwd, _ := os.Getwd()

// 	testCases := []testCase{
// 		{
// 			name:        "invalid incluster",
// 			config:      &rest.Config{},
// 			path:        "",
// 			expectError: true,
// 		},
// 		{
// 			name:        "invalid path",
// 			config:      &rest.Config{},
// 			path:        "invalid",
// 			expectError: true,
// 		},
// 		{
// 			name:   "valid path",
// 			config: &rest.Config{},
// 			path:   pwd + "/../hack/ci/config",
// 		},
// 	}

// 	for _, tcase := range testCases {
// 		suite.T().Run(tcase.name, func(t *testing.T) {
// 			_, err := newKubeAuth(tcase.path)

// 			if tcase.expectError {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 		})
// 	}
// }

// func (suite *cmdTestSuite) TestLoadHelmChart() {
// 	suite.T().Parallel()

// 	type testCase struct {
// 		name        string
// 		createChart bool
// 		expectError bool
// 	}

// 	testDir, err := os.MkdirTemp("", "load-helm-chart")
// 	if err != nil {
// 		suite.T().Fatal(err)
// 	}

// 	defer os.RemoveAll(testDir)

// 	testCases := []testCase{
// 		{
// 			name:        "valid chart",
// 			expectError: false,
// 			createChart: true,
// 		},
// 		{
// 			name:        "missing chart",
// 			expectError: true,
// 			createChart: false,
// 		},
// 	}

// 	for _, tcase := range testCases {
// 		suite.T().Run(tcase.name, func(t *testing.T) {
// 			cpath := ""

// 			if tcase.createChart {
// 				cpath, err = utils.CreateTestChart(testDir)

// 				if err != nil {
// 					t.Fatal(err)
// 				}
// 			}

// 			ch, err := loadHelmChart(cpath)

// 			if tcase.expectError {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.NotNil(t, ch)
// 			}
// 		})
// 	}
// }

// func (suite *cmdTestSuite) TestValidateFlags() {
// 	type flagSet struct {
// 		key   string
// 		value string
// 	}

// 	type testCase struct {
// 		name        string
// 		flagSet     []flagSet
// 		errors      error
// 		expectError bool
// 	}

// 	testCases := []testCase{
// 		{
// 			name:        "valid flags",
// 			flagSet:     []flagSet{{"chart-path", "chart"}, {"nats.subject-prefix", "stream"}},
// 			errors:      nil,
// 			expectError: false,
// 		},
// 		{
// 			name:        "missing chart-path",
// 			flagSet:     []flagSet{{"nats.subject-prefix", "stream"}},
// 			errors:      errChartPath,
// 			expectError: true,
// 		},
// 		{
// 			name:        "missing nats.subject-prefix",
// 			flagSet:     []flagSet{{"chart-path", "chart"}},
// 			errors:      ErrNATSSubjectPrefix,
// 			expectError: true,
// 		},
// 	}

// 	for _, tcase := range testCases {
// 		viper.Reset()
// 		suite.T().Run(tcase.name, func(t *testing.T) {
// 			for _, f := range tcase.flagSet {
// 				viper.Set(f.key, f.value)
// 			}

// 			err := validateFlags()

// 			if tcase.expectError {
// 				assert.Error(t, err)
// 				assert.ErrorContains(t, err, tcase.errors.Error())
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 		})
// 	}
// }
