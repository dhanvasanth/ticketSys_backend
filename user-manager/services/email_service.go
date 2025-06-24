package services

import (
	"ticket-service/config"
	"ticket-service/models"
)

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) GetEmailTemplates(filters map[string]string) ([]*models.EmailTemplate, error) {
	query := "SELECT * FROM email_templates WHERE 1=1"
	args := []interface{}{}
	
	if templateType := filters["template_type"]; templateType != "" {
		query += " AND template_type = ?"
		args = append(args, templateType)
	}
	if organizationID := filters["organization_id"]; organizationID != "" {
		query += " AND organization_id = ?"
		args = append(args, organizationID)
	}
	if isActive := filters["is_active"]; isActive != "" {
		query += " AND is_active = ?"
		args = append(args, isActive == "true")
	}
	
	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var templates []*models.EmailTemplate
	for rows.Next() {
		template := &models.EmailTemplate{}
		err := rows.Scan(&template.ID, &template.Name, &template.Subject,
			&template.Body, &template.TemplateType, &template.OrganizationID,
			&template.IsActive, &template.CreatedAt, &template.UpdatedAt)
		if err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}
	
	return templates, nil
}

func (s *EmailService) CreateEmailTemplate(req *models.CreateEmailTemplateRequest) (*models.EmailTemplate, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	
	query := `INSERT INTO email_templates (name, subject, body, template_type, organization_id, is_active) 
		VALUES (?, ?, ?, ?, ?, ?)`
	
	result, err := config.DB.Exec(query, req.Name, req.Subject, req.Body,
		req.TemplateType, req.OrganizationID, isActive)
	if err != nil {
		return nil, err
	}
	
	id, _ := result.LastInsertId()
	
	// Get created template
	template := &models.EmailTemplate{}
	selectQuery := "SELECT * FROM email_templates WHERE id = ?"
	config.DB.QueryRow(selectQuery, id).Scan(&template.ID, &template.Name, &template.Subject,
		&template.Body, &template.TemplateType, &template.OrganizationID,
		&template.IsActive, &template.CreatedAt, &template.UpdatedAt)
	
	return template, nil
}