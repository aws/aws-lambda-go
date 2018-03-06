package events

type LexEvent struct {
	MessageVersion    string            `json:"messageVersion,omitempty"`
	InvocationSource  string            `json:"invocationSource,omitempty"`
	UserID            string            `json:"userId,omitempty"`
	InputTranscript   string            `json:"inputTranscript,omitempty"`
	SessionAttributes map[string]string `json:"sessionAttributes"`
	RequestAttributes map[string]string `json:"requestAttributes,omitempty"`
	Bot               *LexBot           `json:"bot,omitempty"`
	OutputDialogMode  string            `json:"outputDialogMode,omitempty"`
	CurrentIntent     *LexCurrentIntent `json:"currentIntent,omitempty"`
	DialogAction      *LexDialogAction  `json:"dialogAction,omitempty"`
}

type LexBot struct {
	Name    string `json:"name"`
	Alias   string `json:"alias"`
	Version string `json:"version"`
}

type LexCurrentIntent struct {
	Name               string                `json:"name"`
	Slots              Slots                 `json:"slots"`
	SlotDetails        map[string]SlotDetail `json:"slotDetails"`
	ConfirmationStatus string                `json:"confirmationStatus"`
}

type SlotDetail struct {
	Resolutions   []map[string]string `json:"resolutions"`
	OriginalValue string              `json:"originalValue"`
}

type LexDialogAction struct {
	Type             string            `json:"type"`
	FulfillmentState string            `json:"fulfillmentState"`
	Message          map[string]string `json:"message"`
	IntentName       string            `json:"intentName"`
	Slots            Slots             `json:"slots"`
	SlotToElicit     string            `json:"slotToElicit"`
	ResponseCard     LexResponseCard   `json:"responseCard"`
}

type Slots map[string]string

type LexResponseCard struct {
	Version            int64        `json:"version"`
	ContentType        string       `json:"contentType"`
	GenericAttachments []Attachment `json:"genericAttachments"`
}

type Attachment struct {
	Title             string              `json:"title"`
	SubTitle          string              `json:"subTitle"`
	ImageURL          string              `json:"imageUrl"`
	AttachmentLinkURL string              `json:"attachmentLinkUrl"`
	Buttons           []map[string]string `json:"buttons"`
}

func (h *LexEvent) Clear() {
	h.Bot = nil
	h.CurrentIntent = nil
	h.DialogAction = nil
}
