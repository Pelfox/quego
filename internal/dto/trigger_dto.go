package dto

type CreateTriggerDTO struct {
	FunctionName string `json:"function_name"`
	Payload      string `json:"payload"`
}
