package SM04

type SignatureTypeItem struct {
	ID        string          `json:"id"`
	Code      string          `json:"code"`
	Name      string          `json:"name"`
	CreatedAt string          `json:"created_at"`
	Steps     []ApprovalStepI `json:"steps"`
}

type ApprovalStepI struct {
	ID        string          `json:"id"`
	Step      int             `json:"step"`
	Condition string          `json:"condition"`
	Signers   []ApprovalSignI `json:"signers"`
}

type ApprovalSignI struct {
	UserID string `json:"user_id"`
}

type TypeCreateRequest struct {
	Code  string          `json:"code" binding:"required"`
	Name  string          `json:"name" binding:"required"`
	Steps []StepCreateReq `json:"steps" binding:"required"`
}

type StepCreateReq struct {
	Step      int      `json:"step" binding:"required"`
	Condition string   `json:"condition" binding:"required"`
	UserIDs   []string `json:"user_ids" binding:"required"`
}

type UserListItem struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Fullname  string `json:"fullname"`
	CompanyID string `json:"company_id"`
}

type Repository interface {
	ListSignatureType() ([]SignatureTypeItem, error)
	ListApprovalStepsByType(typeID string) ([]ApprovalStepI, error)
	ListApprovalSignsByStep(stepID string) ([]ApprovalSignI, error)
	CreateSignatureType(code, name string) (string, error)
	InsertApprovalStep(typeID string, step int, condition string) (string, error)
	InsertApprovalSign(stepID, userID string) error
	UpdateSignatureType(id, code, name string) error
	DeleteApprovalStepsByType(typeID string) error
	DeleteApprovalSignsByStep(stepID string) error
	ListUser(companyID, search string) ([]UserListItem, error)
}

type UseCase interface {
	ListSignatureType() ([]SignatureTypeItem, error)
	CreateSignatureType(req TypeCreateRequest) error
	UpdateSignatureType(id string, req TypeCreateRequest) error
	ListUser(companyID, search string) ([]UserListItem, error)
}
