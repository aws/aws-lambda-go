package events

// Amazon Lex V2 event structures for Lambda functions.
//
// Lex V2 provides a conversational interface for building chatbots and voice assistants.
// These structures represent the input events sent from Lex V2 to Lambda functions
// and the expected response format.
//
// Key differences from Lex V1:
//   - sessionId instead of userId
//   - interpretations array with confidence scores
//   - enhanced slot structure with shape and values
//   - sessionState replaces separate sessionAttributes
//   - dialogAction moved into sessionState
//
// Note: This library uses plain string types for simplicity. For type-safe constants
// and validation, refer to the AWS documentation for valid values:
// https://docs.aws.amazon.com/lexv2/latest/dg/lambda-input-format.html
// https://docs.aws.amazon.com/lexv2/latest/dg/lambda-response-format.html
//
// For more information, see:
// https://docs.aws.amazon.com/lexv2/latest/dg/lambda.html

// LexV2Event represents the input event from Amazon Lex V2 to Lambda.
// https://docs.aws.amazon.com/lexv2/latest/dg/lambda-input-format.html
type LexV2Event struct {
	MessageVersion      string                  `json:"messageVersion"`
	InvocationSource    string                  `json:"invocationSource"`
	InputMode           string                  `json:"inputMode"`
	ResponseContentType string                  `json:"responseContentType"`
	SessionID           string                  `json:"sessionId"`
	InputTranscript     string                  `json:"inputTranscript,omitempty"`
	Bot                 LexV2Bot                `json:"bot"`
	Interpretations     []LexV2Interpretation   `json:"interpretations"`
	ProposedNextState   *LexV2ProposedNextState `json:"proposedNextState,omitempty"`
	RequestAttributes   map[string]string       `json:"requestAttributes"`
	SessionState        LexV2SessionState       `json:"sessionState"`
	Transcriptions      []LexV2Transcription    `json:"transcriptions,omitempty"`
}

type LexV2Bot struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AliasID   string `json:"aliasId"`
	AliasName string `json:"aliasName,omitempty"`
	LocaleID  string `json:"localeId"`
	Version   string `json:"version"`
}

type LexV2Interpretation struct {
	Intent               LexV2Intent             `json:"intent"`
	InterpretationSource string                  `json:"interpretationSource,omitempty"`
	NLUConfidence        *LexV2NLUConfidence     `json:"nluConfidence,omitempty"`
	SentimentResponse    *LexV2SentimentResponse `json:"sentimentResponse,omitempty"`
}

type LexV2NLUConfidence struct {
	Score float64 `json:"score,omitempty"`
}

type LexV2Intent struct {
	ConfirmationState string               `json:"confirmationState"`
	Name              string               `json:"name"`
	Slots             map[string]LexV2Slot `json:"slots"`
	State             string               `json:"state"`
	KendraResponse    *LexV2KendraResponse `json:"kendraResponse,omitempty"`
}

type LexV2Slot struct {
	Shape  string           `json:"shape,omitempty"`
	Value  *LexV2SlotValue  `json:"value,omitempty"`
	Values []LexV2SlotValue `json:"values,omitempty"`
}

type LexV2SlotValue struct {
	OriginalValue    string   `json:"originalValue"`
	InterpretedValue string   `json:"interpretedValue"`
	ResolvedValues   []string `json:"resolvedValues,omitempty"`
}

// LexV2KendraResponse contains information about the results of a Kendra search query.
// This field only appears if the intent is a KendraSearchIntent.
// For the complete structure, see:
// https://docs.aws.amazon.com/kendra/latest/dg/API_Query.html#API_Query_ResponseSyntax
type LexV2KendraResponse struct {
	// The response structure matches the Kendra Query API response.
	// Using interface{} to accommodate the full Kendra response without
	// duplicating the entire Kendra API structure here.
	// Users should unmarshal this to the appropriate Kendra types if needed.
}

type LexV2SentimentResponse struct {
	Sentiment      string              `json:"sentiment"`
	SentimentScore LexV2SentimentScore `json:"sentimentScore"`
}

type LexV2SentimentScore struct {
	Mixed    float64 `json:"mixed"`
	Negative float64 `json:"negative"`
	Neutral  float64 `json:"neutral"`
	Positive float64 `json:"positive"`
}

type LexV2ProposedNextState struct {
	DialogAction *LexV2DialogAction `json:"dialogAction,omitempty"`
	Intent       *LexV2Intent       `json:"intent,omitempty"`
}

type LexV2SessionState struct {
	ActiveContexts       []LexV2ActiveContext `json:"activeContexts,omitempty"`
	SessionAttributes    map[string]string    `json:"sessionAttributes,omitempty"`
	RuntimeHints         *LexV2RuntimeHints   `json:"runtimeHints,omitempty"`
	DialogAction         *LexV2DialogAction   `json:"dialogAction,omitempty"`
	Intent               *LexV2Intent         `json:"intent,omitempty"`
	OriginatingRequestID string               `json:"originatingRequestId,omitempty"`
}

type LexV2ActiveContext struct {
	Name              string            `json:"name"`
	ContextAttributes map[string]string `json:"contextAttributes"`
	TimeToLive        LexV2TimeToLive   `json:"timeToLive"`
}

type LexV2TimeToLive struct {
	TimeToLiveInSeconds int `json:"timeToLiveInSeconds"`
	TurnsToLive         int `json:"turnsToLive"`
}

type LexV2RuntimeHints struct {
	SlotHints map[string]map[string]LexV2SlotHint `json:"slotHints,omitempty"`
}

type LexV2SlotHint struct {
	RuntimeHintValues []LexV2RuntimeHintValue `json:"runtimeHintValues"`
}

type LexV2RuntimeHintValue struct {
	Phrase string `json:"phrase"`
}

type LexV2Transcription struct {
	Transcription           string                `json:"transcription"`
	TranscriptionConfidence float64               `json:"transcriptionConfidence,omitempty"`
	ResolvedContext         *LexV2ResolvedContext `json:"resolvedContext,omitempty"`
	ResolvedSlots           map[string]LexV2Slot  `json:"resolvedSlots,omitempty"`
}

type LexV2ResolvedContext struct {
	Intent string `json:"intent"`
}

// LexV2Response represents the response from Lambda to Lex V2.
// https://docs.aws.amazon.com/lexv2/latest/dg/lambda-response-format.html
type LexV2Response struct {
	SessionState      LexV2SessionState `json:"sessionState"`
	Messages          []LexV2Message    `json:"messages,omitempty"`
	RequestAttributes map[string]string `json:"requestAttributes,omitempty"`
}

type LexV2DialogAction struct {
	Type                 string                `json:"type"`
	SlotToElicit         string                `json:"slotToElicit,omitempty"`
	SlotElicitationStyle string                `json:"slotElicitationStyle,omitempty"`
	SubSlotToElicit      *LexV2SubSlotToElicit `json:"subSlotToElicit,omitempty"`
}

type LexV2SubSlotToElicit struct {
	Name            string                `json:"name"`
	SubSlotToElicit *LexV2SubSlotToElicit `json:"subSlotToElicit,omitempty"`
}

type LexV2Message struct {
	ContentType       string                  `json:"contentType"`
	Content           string                  `json:"content,omitempty"`
	ImageResponseCard *LexV2ImageResponseCard `json:"imageResponseCard,omitempty"`
}

type LexV2ImageResponseCard struct {
	Title    string        `json:"title"`
	Subtitle string        `json:"subtitle,omitempty"`
	ImageURL string        `json:"imageUrl,omitempty"`
	Buttons  []LexV2Button `json:"buttons,omitempty"`
}

type LexV2Button struct {
	Text  string `json:"text"`
	Value string `json:"value"`
}
