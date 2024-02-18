package models

type To struct {
	Email string `json:"email"`
}

type Personalizations struct {
	To []To `json:"to"`
}

type From struct {
	Email string `json:"email"`
}

type Content struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type SendGridRequest struct {
	Personalizations Personalizations `json:"personalizations"`
	From             From             `json:"from"`
	Subject          string           `json:"subject"`
	Content          []Content        `json:"content"`
}

type SendGridInternal struct {
	Recipients []string `json:"recipients"`
	Subject    string   `json:"subject"`
	Content    string   `json:"content"`
}
