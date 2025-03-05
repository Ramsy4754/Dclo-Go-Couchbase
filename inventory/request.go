package inventory

type Request struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenantId"`
	Provider  string `json:"provider"`
	AccountID string `json:"accountId"`
	Service   string `json:"service"`
	Resource  string `json:"resource"`

	RequestID string
}
