package SM04

import "obx_rest/mechanic"

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

type RequestItem struct {
	ID             string `json:"id"`
	CompanyID      string `json:"company_id"`
	CompanyName    string `json:"company_name"`
	TypeID         string `json:"type_id"`
	TypeName       string `json:"type_name"`
	RequesterID    string `json:"requester_id"`
	RequesterName  string `json:"requester_name"`
	Code           string `json:"code"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Priority       string `json:"priority"`
	Status         string `json:"status"`
	CurrentStep    int    `json:"current_step"`
	CompletionNote string `json:"completion_note"`
	CompletedBy    string `json:"completed_by"`
	CompletedAt    string `json:"completed_at"`
	CreatedAt      string `json:"created_at"`
}

type RequestCreateRequest struct {
	CompanyID   string `json:"company_id"`
	TypeID      string `json:"type_id" binding:"required"`
	RequesterID string `json:"requester_id" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

type RequestUpdateRequest struct {
	TypeID      string `json:"type_id"`
	RequesterID string `json:"requester_id"`
	Code        string `json:"code"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

type DocumentItem struct {
	ID           string `json:"id"`
	CompanyID    string `json:"company_id"`
	CompanyName  string `json:"company_name"`
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	Code         string `json:"code"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ContentType  string `json:"content_type"`
	Content      string `json:"content"`
	FileName     string `json:"file_name"`
	FileSize     int64  `json:"file_size"`
	Version      string `json:"version"`
	Status       string `json:"status"`
	CreatedBy    string `json:"created_by"`
	ApprovedBy   string `json:"approved_by"`
	ApprovedAt   string `json:"approved_at"`
	IsActive     bool   `json:"is_active"`
	CreatedAt    string `json:"created_at"`
}

type DocumentCreateRequest struct {
	CompanyID   string `json:"company_id"`
	CategoryID  string `json:"category_id"`
	Code        string `json:"code" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	ContentType string `json:"content_type"`
	Content     string `json:"content"`
	FileName    string `json:"file_name"`
	FileSize    int64  `json:"file_size"`
	Version     string `json:"version"`
	Status      string `json:"status"`
}

type DocumentUpdateRequest struct {
	CategoryID  string `json:"category_id"`
	Code        string `json:"code"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ContentType string `json:"content_type"`
	Content     string `json:"content"`
	FileName    string `json:"file_name"`
	FileSize    int64  `json:"file_size"`
	Version     string `json:"version"`
	Status      string `json:"status"`
}

type DocumentCategoryItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DocumentRevisionItem struct {
	ID         string `json:"id"`
	DocumentID string `json:"document_id"`
	Version    string `json:"version"`
	Content    string `json:"content"`
	FilePath   string `json:"file_path"`
	Status     string `json:"status"`
	Note       string `json:"note"`
	CreatedBy  string `json:"created_by"`
	CreatedAt  string `json:"created_at"`
}

type DocumentRevisionCreateRequest struct {
	CompanyID string `json:"company_id"`
	Version   string `json:"version" binding:"required"`
	Content   string `json:"content"`
	FilePath  string `json:"file_path"`
	Status    string `json:"status"`
	Note      string `json:"note"`
}

type DocumentRevisionUpdateRequest struct {
	Version  string `json:"version"`
	Content  string `json:"content"`
	FilePath string `json:"file_path"`
	Status   string `json:"status"`
	Note     string `json:"note"`
}

type DocumentEvidenceItem struct {
	ID         string `json:"id"`
	DocumentID string `json:"document_id"`
	Action     string `json:"action"`
	Note       string `json:"note"`
	UserID     string `json:"user_id"`
	UserName   string `json:"user_name"`
	CreatedAt  string `json:"created_at"`
}

type DocumentEvidenceCreateRequest struct {
	CompanyID string `json:"company_id"`
	Action    string `json:"action" binding:"required"`
	Note      string `json:"note"`
}

type DocumentEvidenceUpdateRequest struct {
	Action string `json:"action"`
	Note   string `json:"note"`
}

type SignatureFormItem struct {
	ID        string              `json:"id"`
	Step      int                 `json:"step"`
	RequestID string              `json:"request_id"`
	Condition string              `json:"condition"`
	Status    string              `json:"status"`
	CreatedAt string              `json:"created_at"`
	Flags     []SignatureFlagItem `json:"flags"`
}

type SignatureFlagItem struct {
	ID        string `json:"id"`
	FormID    string `json:"form_id"`
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	Status    string `json:"status"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
}

type FlagActionRequest struct {
	Status  string `json:"status" binding:"required"`
	Comment string `json:"comment"`
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
	CountRequest(companyID, search string) (int, error)
	ListRequest(companyID, search string, page, size int, sortBy, sortOrder string) ([]RequestItem, error)
	CreateRequest(companyID string, req RequestCreateRequest) error
	UpdateRequest(id string, req RequestUpdateRequest) error
	DeleteRequest(id string) error
	CountDocument(companyID, search string) (int, error)
	ListDocument(companyID, search string, page, size int, sortBy, sortOrder string) ([]DocumentItem, error)
	CreateDocument(companyID, createdBy string, req DocumentCreateRequest) error
	UpdateDocument(id string, req DocumentUpdateRequest) error
	DeleteDocument(id string) error
	ListDocumentCategory(companyID string) ([]DocumentCategoryItem, error)
	CountDocumentRevision(documentID, companyID, search string) (int, error)
	ListDocumentRevision(documentID, companyID, search string, page, size int, sortBy, sortOrder string) ([]DocumentRevisionItem, error)
	CreateDocumentRevision(documentID, companyID, createdBy string, req DocumentRevisionCreateRequest) error
	UpdateDocumentRevision(id string, req DocumentRevisionUpdateRequest) error
	DeleteDocumentRevision(id string) error
	CountDocumentEvidence(documentID, companyID, search string) (int, error)
	ListDocumentEvidence(documentID, companyID, search string, page, size int, sortBy, sortOrder string) ([]DocumentEvidenceItem, error)
	CreateDocumentEvidence(documentID, companyID, userID string, req DocumentEvidenceCreateRequest) error
	UpdateDocumentEvidence(id string, req DocumentEvidenceUpdateRequest) error
	DeleteDocumentEvidence(id string) error

	GetRequestMeta(requestID string) (typeID, status string, currentStep int, err error)
	CountSignatureForms(requestID string) (int, error)
	ListSignatureForms(requestID, companyID string) ([]SignatureFormItem, error)
	ListSignatureFlags(formID string) ([]SignatureFlagItem, error)
	GetSignatureForm(formID string) (requestID, condition string, step int, err error)
	GetSignatureFlag(flagID string) (formID, status string, err error)
	InsertSignatureForm(requestID string, step int, condition string) (string, error)
	InsertSignatureFlag(formID, userID string) error
	UpdateSignatureFlag(flagID, status, comment string) error
	UpdateSignatureFormStatus(formID, status string) error
	UpdateRequestStatus(id, status string) error
	UpdateRequestStep(id string, step int) error
	UpdateRequestCompletion(id, status, completedBy, completionNote string) error
}

type UseCase interface {
	ListSignatureType() ([]SignatureTypeItem, error)
	CreateSignatureType(req TypeCreateRequest) error
	UpdateSignatureType(id string, req TypeCreateRequest) error
	ListUser(companyID, search string) ([]UserListItem, error)
	ListRequest(companyID string, meta mechanic.ActionMeta) ([]RequestItem, mechanic.GridMeta, error)
	CreateRequest(companyID string, req RequestCreateRequest) error
	UpdateRequest(id string, req RequestUpdateRequest) error
	DeleteRequest(id string) error
	ListDocument(companyID string, meta mechanic.ActionMeta) ([]DocumentItem, mechanic.GridMeta, error)
	CreateDocument(companyID, createdBy string, req DocumentCreateRequest) error
	UpdateDocument(id string, req DocumentUpdateRequest) error
	DeleteDocument(id string) error
	ListDocumentCategory(companyID string) ([]DocumentCategoryItem, error)
	ListDocumentRevision(documentID, companyID string, meta mechanic.ActionMeta) ([]DocumentRevisionItem, mechanic.GridMeta, error)
	CreateDocumentRevision(documentID, companyID, createdBy string, req DocumentRevisionCreateRequest) error
	UpdateDocumentRevision(id string, req DocumentRevisionUpdateRequest) error
	DeleteDocumentRevision(id string) error
	ListDocumentEvidence(documentID, companyID string, meta mechanic.ActionMeta) ([]DocumentEvidenceItem, mechanic.GridMeta, error)
	CreateDocumentEvidence(documentID, companyID, userID string, req DocumentEvidenceCreateRequest) error
	UpdateDocumentEvidence(id string, req DocumentEvidenceUpdateRequest) error
	DeleteDocumentEvidence(id string) error

	ListForm(requestID, companyID string) ([]SignatureFormItem, error)
	GenerateForms(requestID string) (int, error)
	FlagAction(flagID, userID string, req FlagActionRequest) error
}
