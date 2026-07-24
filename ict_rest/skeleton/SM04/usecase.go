package SM04

import (
	"ict_rest/mechanic"
)

type useCase struct {
	repo Repository
}

func NCase(r Repository) UseCase {
	return &useCase{repo: r}
}

func (u *useCase) ListSignatureType() ([]SignatureTypeItem, error) {
	return u.repo.ListSignatureType()
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
	return u.repo.ListUser(companyID, search)
}
