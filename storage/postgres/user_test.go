package postgres

import "testing"

func CreateUser(t *testing.T) {
	
}

// func TestCreateInstallment(t *testing.T) {
// 	katmid := createRandomId(t)
// 	var tests = []struct {
// 	  give *installment_service.Installment
// 	  want error
// 	}{
// 	  {
// 		&installment_service.Installment{
// 		  MerchantBranchId: createMerchantBranch(t).Id,
// 		  CustomerId:       CreateCustomer(t).Id,
// 		  AgentId:          createAgent(t).Id,
// 		  StageId:          createStage(t).Id,
// 		  StatusId:         "aca171ab-e8ea-45eb-94e3-b3968b4ec225",
// 		  FirstPaymentDate: "2020-11-20",
// 		  TermMonth: &installment_service.Catalogue{
// 			Guid:  createRandomId(t),
// 			Label: fakeData.JobTitle(),
// 		  },
// 		  DownPayment:          fakeData.Rand.Float64(),
// 		  TotalAmount:          fakeData.Rand.Float64(),
// 		  MonthlyPaymentAmount: fakeData.Rand.Float64(),
// 		  PaymeUrl:             fakeData.URL(),
// 		  ContractNumber:       fakeData.CellPhoneNumber(),
// 		  KatmId:               katmid,
// 		},
// 		nil,
// 	  },
// 	}

// 	installmentRepo := NewInstallmentRepo(db)

// 	for i, tt := range tests {
// 	  testname := fmt.Sprintf("Test %d", i+1)
// 	  t.Run(testname, func(t *testing.T) {
// 		resp, err := installmentRepo.Create(context.Background(), &installment_service.CreateInstallmentRequest{
// 		  Installment: tt.give,
// 		})

// 		assert.NoError(t, err)

// 		fmt.Print("Create Installment------->")

// 		b, err := json.MarshalIndent(resp, "", "  ")

// 		if !assert.Equal(t, tt.want, err) {
// 		  t.Errorf("got %s, want %s", err, tt.want)
// 		  return
// 		}
// 		fmt.Println(string(b))
// 	  })
// 	}
//   }
