package SM04

import (
	"obx_rest/mechanic"
	"strings"
)

type useCase struct {
	repo Repository
}

func NCase(r Repository) UseCase {
	return &useCase{repo: r}
}

func (u *useCase) ListSignatureType() ([]SignatureTypeItem, error) {
	list, err := u.repo.ListSignatureType()
	if err != nil {
		return nil, mechanic.InternalError("Failed to list signature types", err)
	}
	return list, nil
}

func (u *useCase) CreateSignatureType(req TypeCreateRequest) error {
	if req.Code == "" {
		return mechanic.ValidationError("Code is required")
	}
	if req.Name == "" {
		return mechanic.ValidationError("Name is required")
	}
	if len(req.Steps) == 0 {
		return mechanic.ValidationError("At least one approval step is required")
	}
	typeID, err := u.repo.CreateSignatureType(req.Code, req.Name)
	if err != nil {
		return mechanic.InternalError("failed to create signature type", err)
	}
	for _, step := range req.Steps {
		stepID, err := u.repo.InsertApprovalStep(typeID, step.Step, step.Condition)
		if err != nil {
			return mechanic.InternalError("failed to create approval step", err)
		}
		for _, uid := range step.UserIDs {
			if err := u.repo.InsertApprovalSign(stepID, uid); err != nil {
				return mechanic.InternalError("failed to assign signer", err)
			}
		}
	}
	return nil
}

func (u *useCase) UpdateSignatureType(id string, req TypeCreateRequest) error {
	if id == "" {
		return mechanic.ValidationError("Signature type ID is required")
	}
	if req.Code == "" {
		return mechanic.ValidationError("Code is required")
	}
	if req.Name == "" {
		return mechanic.ValidationError("Name is required")
	}
	if len(req.Steps) == 0 {
		return mechanic.ValidationError("At least one approval step is required")
	}
	existing, err := u.repo.ListApprovalStepsByType(id)
	if err != nil {
		return mechanic.InternalError("failed to load existing approval steps", err)
	}
	for _, step := range existing {
		_ = u.repo.DeleteApprovalSignsByStep(step.ID)
	}
	if err := u.repo.DeleteApprovalStepsByType(id); err != nil {
		return mechanic.InternalError("failed to remove existing approval steps", err)
	}
	if err := u.repo.UpdateSignatureType(id, req.Code, req.Name); err != nil {
		return mechanic.InternalError("failed to update signature type", err)
	}
	for _, step := range req.Steps {
		stepID, err := u.repo.InsertApprovalStep(id, step.Step, step.Condition)
		if err != nil {
			return mechanic.InternalError("failed to create approval step", err)
		}
		for _, uid := range step.UserIDs {
			if err := u.repo.InsertApprovalSign(stepID, uid); err != nil {
				return mechanic.InternalError("failed to assign signer", err)
			}
		}
	}
	return nil
}

func (u *useCase) ListUser(companyID, search string) ([]UserListItem, error) {
	list, err := u.repo.ListUser(companyID, search)
	if err != nil {
		return nil, mechanic.InternalError("Failed to list signer users", err)
	}
	return list, nil
}

func (u *useCase) ListRequest(companyID string, meta mechanic.ActionMeta) (
	[]RequestItem, mechanic.GridMeta, error) {
	page, size, err := mechanic.CheckMeta(meta.Page, meta.Size)
	if err != nil {
		return nil, mechanic.GridMeta{}, err
	}
	total, err := u.repo.CountRequest(companyID, meta.Search)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to count requests", err)
	}
	list, err := u.repo.ListRequest(companyID, meta.Search, page, size, meta.SortBy, meta.SortOrder)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to list requests", err)
	}
	return list, mechanic.BuildMeta(page, size, total), nil
}

func (u *useCase) CreateRequest(companyID string, req RequestCreateRequest) error {
	if strings.TrimSpace(companyID) == "" {
		companyID = req.CompanyID
	}
	if strings.TrimSpace(companyID) == "" {
		return mechanic.ValidationError("Company ID is required")
	}
	if strings.TrimSpace(req.TypeID) == "" {
		return mechanic.ValidationError("Signature type is required")
	}
	if strings.TrimSpace(req.RequesterID) == "" {
		return mechanic.ValidationError("Requester is required")
	}
	if strings.TrimSpace(req.Code) == "" {
		return mechanic.ValidationError("Request code is required")
	}
	if strings.TrimSpace(req.Title) == "" {
		return mechanic.ValidationError("Title is required")
	}
	if req.Priority == "" {
		req.Priority = "NORMAL"
	}
	if err := u.repo.CreateRequest(companyID, req); err != nil {
		return mechanic.InternalError("Failed to create request", err)
	}
	return nil
}

func (u *useCase) UpdateRequest(id string, req RequestUpdateRequest) error {
	if strings.TrimSpace(id) == "" {
		return mechanic.ValidationError("Request ID is required")
	}
	if strings.TrimSpace(req.TypeID) == "" {
		return mechanic.ValidationError("Signature type is required")
	}
	if strings.TrimSpace(req.RequesterID) == "" {
		return mechanic.ValidationError("Requester is required")
	}
	if strings.TrimSpace(req.Code) == "" {
		return mechanic.ValidationError("Request code is required")
	}
	if strings.TrimSpace(req.Title) == "" {
		return mechanic.ValidationError("Title is required")
	}
	if req.Priority == "" {
		req.Priority = "NORMAL"
	}
	if err := u.repo.UpdateRequest(id, req); err != nil {
		return mechanic.InternalError("Failed to update request", err)
	}
	return nil
}

func (u *useCase) DeleteRequest(id string) error {
	if strings.TrimSpace(id) == "" {
		return mechanic.ValidationError("Request ID is required")
	}
	if err := u.repo.DeleteRequest(id); err != nil {
		return mechanic.InternalError("Failed to delete request", err)
	}
	return nil
}

func (u *useCase) ListDocument(companyID string, meta mechanic.ActionMeta) (
	[]DocumentItem, mechanic.GridMeta, error) {
	page, size, err := mechanic.CheckMeta(meta.Page, meta.Size)
	if err != nil {
		return nil, mechanic.GridMeta{}, err
	}
	total, err := u.repo.CountDocument(companyID, meta.Search)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to count documents", err)
	}
	list, err := u.repo.ListDocument(companyID, meta.Search, page, size, meta.SortBy, meta.SortOrder)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to list documents", err)
	}
	return list, mechanic.BuildMeta(page, size, total), nil
}

func (u *useCase) CreateDocument(companyID, createdBy string, req DocumentCreateRequest) error {
	if strings.TrimSpace(companyID) == "" {
		companyID = req.CompanyID
	}
	if strings.TrimSpace(companyID) == "" {
		return mechanic.ValidationError("Company ID is required")
	}
	if strings.TrimSpace(req.Code) == "" {
		return mechanic.ValidationError("Document code is required")
	}
	if strings.TrimSpace(req.Title) == "" {
		return mechanic.ValidationError("Title is required")
	}
	if req.ContentType == "" {
		req.ContentType = "markdown"
	}
	if req.Version == "" {
		req.Version = "1.0"
	}
	if req.Status == "" {
		req.Status = "DRAFT"
	}
	if err := u.repo.CreateDocument(companyID, createdBy, req); err != nil {
		return mechanic.InternalError("Failed to create document", err)
	}
	return nil
}

func (u *useCase) UpdateDocument(id string, req DocumentUpdateRequest) error {
	if strings.TrimSpace(id) == "" {
		return mechanic.ValidationError("Document ID is required")
	}
	if strings.TrimSpace(req.Code) == "" {
		return mechanic.ValidationError("Document code is required")
	}
	if strings.TrimSpace(req.Title) == "" {
		return mechanic.ValidationError("Title is required")
	}
	if req.ContentType == "" {
		req.ContentType = "markdown"
	}
	if req.Version == "" {
		req.Version = "1.0"
	}
	if err := u.repo.UpdateDocument(id, req); err != nil {
		return mechanic.InternalError("Failed to update document", err)
	}
	return nil
}

func (u *useCase) DeleteDocument(id string) error {
	if strings.TrimSpace(id) == "" {
		return mechanic.ValidationError("Document ID is required")
	}
	if err := u.repo.DeleteDocument(id); err != nil {
		return mechanic.InternalError("Failed to delete document", err)
	}
	return nil
}

func (u *useCase) ListDocumentCategory(companyID string) ([]DocumentCategoryItem, error) {
	list, err := u.repo.ListDocumentCategory(companyID)
	if err != nil {
		return nil, mechanic.InternalError("Failed to list document categories", err)
	}
	return list, nil
}

func (u *useCase) ListDocumentRevision(documentID, companyID string, meta mechanic.ActionMeta) (
	[]DocumentRevisionItem, mechanic.GridMeta, error) {
	if strings.TrimSpace(documentID) == "" {
		return nil, mechanic.GridMeta{}, mechanic.ValidationError("Document ID is required")
	}
	page, size, err := mechanic.CheckMeta(meta.Page, meta.Size)
	if err != nil {
		return nil, mechanic.GridMeta{}, err
	}
	total, err := u.repo.CountDocumentRevision(documentID, companyID, meta.Search)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to count document revisions", err)
	}
	list, err := u.repo.ListDocumentRevision(documentID, companyID, meta.Search, page, size, meta.SortBy, meta.SortOrder)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to list document revisions", err)
	}
	return list, mechanic.BuildMeta(page, size, total), nil
}

func (u *useCase) CreateDocumentRevision(documentID, companyID, createdBy string, req DocumentRevisionCreateRequest) error {
	if strings.TrimSpace(documentID) == "" {
		return mechanic.ValidationError("Document ID is required")
	}
	if strings.TrimSpace(companyID) == "" {
		companyID = req.CompanyID
	}
	if strings.TrimSpace(companyID) == "" {
		return mechanic.ValidationError("Company ID is required")
	}
	if strings.TrimSpace(req.Version) == "" {
		return mechanic.ValidationError("Version is required")
	}
	if req.Status == "" {
		req.Status = "DRAFT"
	}
	if err := u.repo.CreateDocumentRevision(documentID, companyID, createdBy, req); err != nil {
		return mechanic.InternalError("Failed to create document revision", err)
	}
	return nil
}

func (u *useCase) UpdateDocumentRevision(id string, req DocumentRevisionUpdateRequest) error {
	if strings.TrimSpace(id) == "" {
		return mechanic.ValidationError("Document revision ID is required")
	}
	if strings.TrimSpace(req.Version) == "" {
		return mechanic.ValidationError("Version is required")
	}
	if req.Status == "" {
		req.Status = "DRAFT"
	}
	if err := u.repo.UpdateDocumentRevision(id, req); err != nil {
		return mechanic.InternalError("Failed to update document revision", err)
	}
	return nil
}

func (u *useCase) DeleteDocumentRevision(id string) error {
	if strings.TrimSpace(id) == "" {
		return mechanic.ValidationError("Document revision ID is required")
	}
	if err := u.repo.DeleteDocumentRevision(id); err != nil {
		return mechanic.InternalError("Failed to delete document revision", err)
	}
	return nil
}

func (u *useCase) ListDocumentEvidence(documentID, companyID string, meta mechanic.ActionMeta) (
	[]DocumentEvidenceItem, mechanic.GridMeta, error) {
	if strings.TrimSpace(documentID) == "" {
		return nil, mechanic.GridMeta{}, mechanic.ValidationError("Document ID is required")
	}
	page, size, err := mechanic.CheckMeta(meta.Page, meta.Size)
	if err != nil {
		return nil, mechanic.GridMeta{}, err
	}
	total, err := u.repo.CountDocumentEvidence(documentID, companyID, meta.Search)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to count document evidence", err)
	}
	list, err := u.repo.ListDocumentEvidence(documentID, companyID, meta.Search, page, size, meta.SortBy, meta.SortOrder)
	if err != nil {
		return nil, mechanic.GridMeta{}, mechanic.InternalError("Failed to list document evidence", err)
	}
	return list, mechanic.BuildMeta(page, size, total), nil
}

func (u *useCase) CreateDocumentEvidence(documentID, companyID, userID string, req DocumentEvidenceCreateRequest) error {
	if strings.TrimSpace(documentID) == "" {
		return mechanic.ValidationError("Document ID is required")
	}
	if strings.TrimSpace(companyID) == "" {
		companyID = req.CompanyID
	}
	if strings.TrimSpace(companyID) == "" {
		return mechanic.ValidationError("Company ID is required")
	}
	if strings.TrimSpace(req.Action) == "" {
		return mechanic.ValidationError("Action is required")
	}
	if err := u.repo.CreateDocumentEvidence(documentID, companyID, userID, req); err != nil {
		return mechanic.InternalError("Failed to create document evidence", err)
	}
	return nil
}

func (u *useCase) UpdateDocumentEvidence(id string, req DocumentEvidenceUpdateRequest) error {
	if strings.TrimSpace(id) == "" {
		return mechanic.ValidationError("Document evidence ID is required")
	}
	if strings.TrimSpace(req.Action) == "" {
		return mechanic.ValidationError("Action is required")
	}
	if err := u.repo.UpdateDocumentEvidence(id, req); err != nil {
		return mechanic.InternalError("Failed to update document evidence", err)
	}
	return nil
}

func (u *useCase) DeleteDocumentEvidence(id string) error {
	if strings.TrimSpace(id) == "" {
		return mechanic.ValidationError("Document evidence ID is required")
	}
	if err := u.repo.DeleteDocumentEvidence(id); err != nil {
		return mechanic.InternalError("Failed to delete document evidence", err)
	}
	return nil
}

func (u *useCase) ListForm(requestID, companyID string) ([]SignatureFormItem, error) {
	if strings.TrimSpace(requestID) == "" {
		return nil, mechanic.ValidationError("Request ID is required")
	}
	list, err := u.repo.ListSignatureForms(requestID, companyID)
	if err != nil {
		return nil, mechanic.InternalError("Failed to list signature forms", err)
	}
	return list, nil
}

func (u *useCase) GenerateForms(requestID string) (int, error) {
	if strings.TrimSpace(requestID) == "" {
		return 0, mechanic.ValidationError("Request ID is required")
	}
	typeID, _, _, err := u.repo.GetRequestMeta(requestID)
	if err != nil {
		return 0, mechanic.NotFound("Request not found")
	}
	exists, err := u.repo.CountSignatureForms(requestID)
	if err != nil {
		return 0, mechanic.InternalError("Failed to check existing signature forms", err)
	}
	if exists > 0 {
		return 0, mechanic.Conflict("Signature forms already generated for this request")
	}
	steps, err := u.repo.ListApprovalStepsByType(typeID)
	if err != nil {
		return 0, mechanic.InternalError("Failed to load approval steps", err)
	}
	if len(steps) == 0 {
		return 0, mechanic.ValidationError("Signature type has no approval step")
	}
	for _, step := range steps {
		formID, err := u.repo.InsertSignatureForm(requestID, step.Step, step.Condition)
		if err != nil {
			return 0, mechanic.InternalError("Failed to create signature form", err)
		}
		for _, signer := range step.Signers {
			if err := u.repo.InsertSignatureFlag(formID, signer.UserID); err != nil {
				return 0, mechanic.InternalError("Failed to create signature flag", err)
			}
		}
	}
	return len(steps), nil
}

func (u *useCase) FlagAction(flagID, userID string, req FlagActionRequest) error {
	if strings.TrimSpace(flagID) == "" {
		return mechanic.ValidationError("Signature flag ID is required")
	}
	status := strings.ToUpper(strings.TrimSpace(req.Status))
	if status != "APPROVED" && status != "REJECTED" {
		return mechanic.ValidationError("Flag status must be APPROVED or REJECTED")
	}
	formID, _, err := u.repo.GetSignatureFlag(flagID)
	if err != nil {
		return mechanic.NotFound("Signature flag not found")
	}
	requestID, condition, _, err := u.repo.GetSignatureForm(formID)
	if err != nil {
		return mechanic.NotFound("Signature form not found")
	}
	typeID, reqStatus, currentStep, err := u.repo.GetRequestMeta(requestID)
	if err != nil {
		return mechanic.NotFound("Request not found")
	}
	if reqStatus == "APPROVED" || reqStatus == "REJECTED" {
		return mechanic.ValidationError("Request is already " + reqStatus)
	}
	if err := u.repo.UpdateSignatureFlag(flagID, status, req.Comment); err != nil {
		return mechanic.InternalError("Failed to update signature flag", err)
	}
	flags, err := u.repo.ListSignatureFlags(formID)
	if err != nil {
		return mechanic.InternalError("Failed to load signature flags", err)
	}
	rejected, approved, total := false, 0, len(flags)
	for _, f := range flags {
		switch f.Status {
		case "REJECTED":
			rejected = true
		case "APPROVED":
			approved++
		}
	}
	var formStatus string
	switch {
	case rejected:
		formStatus = "REJECTED"
	case condition == "ALL_APPROVED" && approved == total:
		formStatus = "APPROVED"
	case condition == "ANY_APPROVED" && approved >= 1:
		formStatus = "APPROVED"
	case approved >= 1:
		formStatus = "IN_PROGRESS"
	default:
		formStatus = "PENDING"
	}
	if err := u.repo.UpdateSignatureFormStatus(formID, formStatus); err != nil {
		return mechanic.InternalError("Failed to update signature form status", err)
	}

	switch formStatus {
	case "REJECTED":
		if err := u.repo.UpdateRequestCompletion(requestID, "REJECTED", userID, req.Comment); err != nil {
			return mechanic.InternalError("Failed to update request status", err)
		}
	case "APPROVED":
		steps, err := u.repo.ListApprovalStepsByType(typeID)
		if err != nil {
			return mechanic.InternalError("Failed to load approval steps", err)
		}
		next := 0
		for _, s := range steps {
			if s.Step > currentStep && (next == 0 || s.Step < next) {
				next = s.Step
			}
		}
		if next == 0 {
			if err := u.repo.UpdateRequestCompletion(requestID, "APPROVED", userID, req.Comment); err != nil {
				return mechanic.InternalError("Failed to finalize request", err)
			}
		} else {
			if err := u.repo.UpdateRequestStep(requestID, next); err != nil {
				return mechanic.InternalError("Failed to advance request step", err)
			}
			if err := u.repo.UpdateRequestStatus(requestID, "IN_PROGRESS"); err != nil {
				return mechanic.InternalError("Failed to update request status", err)
			}
		}
	}
	return nil
}
