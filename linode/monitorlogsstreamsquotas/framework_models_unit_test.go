//go:build unit

package monitorlogsstreamsquotas

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
)

func TestLogStreamQuotaListModel_parseQuotas_FromSampleAPIResponse(t *testing.T) {
	// This matches the payload shape returned by the API for monitor/streams/quotas.
	quotas := []linodego.LogStreamQuota{
		{
			QuotaID:     "aclp_audit_logs_streams",
			QuotaName:   "Number of Audit Logs Streams",
			Description: "Current number of audit logs streams per account",
			QuotaLimit:  5,
			QuotaType:   "logs_streams_aclp_audit_logs_streams",
		},
		{
			QuotaID:     "aclp_lke_audit_logs_streams",
			QuotaName:   "Number of LKE Audit Logs Streams",
			Description: "Current number of LKE audit logs streams per account",
			QuotaLimit:  10,
			QuotaType:   "logs_streams_aclp_lke_audit_logs_streams",
		},
	}

	model := &LogStreamQuotaListModel{}
	model.parseQuotas(quotas)

	if assert.Len(t, model.Quotas, 2) {
		assert.Equal(t, types.StringValue("aclp_audit_logs_streams"), model.Quotas[0].QuotaID)
		assert.Equal(t, types.StringValue("Number of Audit Logs Streams"), model.Quotas[0].QuotaName)
		assert.Equal(t, types.StringValue("Current number of audit logs streams per account"), model.Quotas[0].Description)
		assert.Equal(t, types.Int64Value(5), model.Quotas[0].QuotaLimit)
		assert.Equal(t, types.StringValue("logs_streams_aclp_audit_logs_streams"), model.Quotas[0].QuotaType)

		assert.Equal(t, types.StringValue("aclp_lke_audit_logs_streams"), model.Quotas[1].QuotaID)
		assert.Equal(t, types.StringValue("Number of LKE Audit Logs Streams"), model.Quotas[1].QuotaName)
		assert.Equal(t, types.StringValue("Current number of LKE audit logs streams per account"), model.Quotas[1].Description)
		assert.Equal(t, types.Int64Value(10), model.Quotas[1].QuotaLimit)
		assert.Equal(t, types.StringValue("logs_streams_aclp_lke_audit_logs_streams"), model.Quotas[1].QuotaType)
	}
}
