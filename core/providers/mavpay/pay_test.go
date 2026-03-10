package mavpay

import (
	"encoding/json"

	"github.com/mavryk-network/mavpay/common"
)

const (
	mavpayHome = "/mnt/c/Users/SD/Projects/mavrykdynamics/peak/mavpay"
)

type generatePayoutResult struct {
	Cycles               []int64                     `json:"cycles"`
	CyclePayoutBlueprint common.CyclePayoutBlueprint `json:"cycle_payout_blueprint"`
}

func getGeneratePayoutsResult(message string) (*generatePayoutResult, error) {
	var result generatePayoutResult
	err := json.Unmarshal([]byte(message), &result)
	return &result, err
}

// func TestGeneratePayouts(t *testing.T) {
// 	assert := assert.New(t)

// 	config := configuration.Runtime{
// 		MavpayHome: mavpayHome,
// 	}

// 	mavpayProvider, err := NewMavpayProvider(&config)
// 	assert.Nil(err)

// 	outputChannel := make(chan string)
// 	go mavpayProvider.GeneratePayouts(-1, outputChannel)

// 	phases := []string{}
// 	type message struct {
// 		Phase string `json:"phase"`
// 	}

// 	for output := range outputChannel {
// 		var m message
// 		_ = json.Unmarshal([]byte(output), &m)

// 		phases = append(phases, m.Phase)
// 	}

// 	assert.Contains(phases, "check_conditions_and_prepare")
// 	assert.Contains(phases, "generate_payout_candidates")
// 	assert.Contains(phases, "collect_transaction_fees")
// 	assert.Contains(phases, "validate_simulated_payouts")
// 	assert.Contains(phases, "finalize_payouts")
// 	assert.Contains(phases, "create_blueprint")
// 	assert.Contains(phases, "result")
//  assert.Contains(phases, "execution_finished")
// }

// func TestPay(t *testing.T) {
// 	assert := assert.New(t)

// 	config := configuration.Runtime{
// 		MavpayHome: mavpayHome,
// 	}

// 	assert.Nil(os.RemoveAll(path.Join(mavpayHome, "reports")))

// 	mavpayProvider, err := NewMavpayProvider(&config)
// 	assert.Nil(err)

// 	outputChannel := make(chan string)
// 	go mavpayProvider.GeneratePayouts(-1, outputChannel)

// 	phases := []string{}
// 	type message struct {
// 		Phase string `json:"phase"`
// 	}

// 	var resultMessage *generatePayoutResult

// 	for output := range outputChannel {
// 		var m message
// 		_ = json.Unmarshal([]byte(output), &m)

// 		phases = append(phases, m.Phase)
// 		if m.Phase == "result" {
// 			resultMessage, err = getGeneratePayoutsResult(output)
// 			assert.Nil(err)
// 		}
// 	}

// 	assert.NotNil(resultMessage)
// 	assert.Contains(phases, "check_conditions_and_prepare")
// 	assert.Contains(phases, "generate_payout_candidates")
// 	assert.Contains(phases, "collect_transaction_fees")
// 	assert.Contains(phases, "validate_simulated_payouts")
// 	assert.Contains(phases, "finalize_payouts")
// 	assert.Contains(phases, "create_blueprint")
// 	assert.Contains(phases, "result")
// 	assert.Contains(phases, "execution_finished")

// 	outputChannel = make(chan string)
// 	go mavpayProvider.Pay(&resultMessage.CyclePayoutBlueprint, outputChannel)

// 	phases = []string{}
// 	for output := range outputChannel {
// 		var m message
// 		_ = json.Unmarshal([]byte(output), &m)

// 		phases = append(phases, m.Phase)
// 	}

// 	assert.Contains(phases, "split_into_batches")
// 	assert.Contains(phases, "execute_payouts")
// 	assert.Contains(phases, "result")
// 	assert.Contains(phases, "execution_finished")
// }

// func TestVersion(t *testing.T) {
// 	assert := assert.New(t)

// 	config := configuration.Runtime{
// 		Modules: map[string]json.RawMessage{
// 			"mavpay": json.RawMessage(`{
// 				"app_path": "mavpay"
// 			}`),
// 		},
// 	}

// 	err := setupMavpayProvider(&config, nil)
// 	assert.Nil(err)

// 	ver, err := mavpayProvider.Version()
// 	assert.Nil(err)
// 	assert.NotEmpty(ver)
// 	assert.Contains(ver, "mavpay:")
// 	assert.Contains(ver, "ami-mavpay:")

// 	// reports, err := mavpayProvider.ListDryReports()
// 	// assert.Nil(err)
// 	// assert.NotEmpty(reports)

// 	// fmt.Println(mavpayProvider.GetReport("754", true))
// }
