package clients

var DistrQueryClientInstance = &DistributionQueryClient{}
var BankQueryClientInstance = &BankQueryClient{}

func init() {
	DistrQueryClientInstance.New()
	BankQueryClientInstance.New()
}
