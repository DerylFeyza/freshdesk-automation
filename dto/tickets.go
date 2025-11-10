package dto

import "time"

type Ticket struct {
	CCEmails              []string      `json:"cc_emails"`
	FwdEmails             []string      `json:"fwd_emails"`
	ReplyCCEmails         []string      `json:"reply_cc_emails"`
	TicketCCEmails        []string      `json:"ticket_cc_emails"`
	TicketBCCEmails       []string      `json:"ticket_bcc_emails"`
	FREscalated           bool          `json:"fr_escalated"`
	Spam                  bool          `json:"spam"`
	EmailConfigID         int64         `json:"email_config_id"`
	GroupID               *int64        `json:"group_id"`
	Priority              int           `json:"priority"`
	RequesterID           int64         `json:"requester_id"`
	ResponderID           *int64        `json:"responder_id"`
	Source                int           `json:"source"`
	CompanyID             *int64        `json:"company_id"`
	Status                int           `json:"status"`
	Subject               string        `json:"subject"`
	AssociationType       *string       `json:"association_type"`
	SupportEmail          string        `json:"support_email"`
	ToEmails              []string      `json:"to_emails"`
	ProductID             *int64        `json:"product_id"`
	ID                    int64         `json:"id"`
	Type                  *string       `json:"type"`
	DueBy                 time.Time     `json:"due_by"`
	FRDueBy               time.Time     `json:"fr_due_by"`
	IsEscalated           bool          `json:"is_escalated"`
	Description           string        `json:"description"`
	DescriptionText       string        `json:"description_text"`
	CustomFields          CustomFields  `json:"custom_fields"`
	CreatedAt             time.Time     `json:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at"`
	Tags                  []string      `json:"tags"`
	Attachments           []interface{} `json:"attachments"`
	SourceAdditionalInfo  *string       `json:"source_additional_info"`
	StructuredDescription *string       `json:"structured_description"`
	NRDueBy               *string       `json:"nr_due_by"`
	NREscalated           bool          `json:"nr_escalated"`
}

type CustomFields struct {
	CFFSMContactName          *string `json:"cf_fsm_contact_name"`
	CFFSMPhoneNumber          *string `json:"cf_fsm_phone_number"`
	CFFSMServiceLocation      *string `json:"cf_fsm_service_location"`
	CFFSMAppointmentStartTime *string `json:"cf_fsm_appointment_start_time"`
	CFFSMAppointmentEndTime   *string `json:"cf_fsm_appointment_end_time"`
}
