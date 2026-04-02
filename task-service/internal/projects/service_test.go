package projects

import (
	"context"
	"testing"

	"github.com/task-manager/task-service/pkg/apperror"
	"github.com/task-manager/task-service/pkg/models"
)

const (
	testOwnerID   = "owner-1"
	testMemberID  = "member-1"
	testProjectID = "project-1"
)

// mockRepository implements RepositoryInterface for testing
type mockRepository struct {
	projects map[string]*models.Project
	lastID   int
}

func newMockRepository() *mockRepository {
	return &mockRepository{projects: make(map[string]*models.Project)}
}

func (m *mockRepository) GetByUser(ctx context.Context, userID string) ([]models.Project, error) {
	var result []models.Project
	for _, p := range m.projects {
		for _, member := range p.Members {
			if member == userID {
				result = append(result, *p)
				break
			}
		}
	}
	return result, nil
}

func (m *mockRepository) GetByID(ctx context.Context, projectID string) (*models.Project, error) {
	p, ok := m.projects[projectID]
	if !ok {
		return nil, apperror.NotFound("project")
	}
	return p, nil
}

func (m *mockRepository) Create(ctx context.Context, project *models.Project) (string, error) {
	m.lastID++
	id := "proj-" + string(rune('0'+m.lastID))
	m.projects[id] = project
	return id, nil
}

func (m *mockRepository) Update(ctx context.Context, projectID string, updates map[string]interface{}) error {
	p, ok := m.projects[projectID]
	if !ok {
		return apperror.NotFound("project")
	}
	for key, val := range updates {
		switch key {
		case "name":
			p.Name = val.(string)
		case "description":
			p.Description = val.(string)
		}
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, projectID string) error {
	delete(m.projects, projectID)
	return nil
}

func (m *mockRepository) AddMemberTx(ctx context.Context, projectID string, memberID string) error {
	p, ok := m.projects[projectID]
	if !ok {
		return apperror.NotFound("project")
	}
	for _, member := range p.Members {
		if member == memberID {
			return apperror.Conflict("user is already a member")
		}
	}
	p.Members = append(p.Members, memberID)
	return nil
}

func (m *mockRepository) RemoveMemberTx(ctx context.Context, projectID string, memberID string) error {
	p, ok := m.projects[projectID]
	if !ok {
		return apperror.NotFound("project")
	}
	if memberID == p.OwnerID {
		return apperror.BadRequest("cannot remove the project owner")
	}
	for i, member := range p.Members {
		if member == memberID {
			p.Members = append(p.Members[:i], p.Members[i+1:]...)
			return nil
		}
	}
	return nil
}

func TestCreateProject(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)

	project, err := svc.Create(context.Background(), CreateProjectRequest{
		Name:        "Test Project",
		Description: "Description",
	}, testOwnerID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if project.Name != "Test Project" {
		t.Errorf("expected name 'Test Project', got '%s'", project.Name)
	}
	if project.OwnerID != testOwnerID {
		t.Errorf("expected ownerID '%s', got '%s'", testOwnerID, project.OwnerID)
	}
	if len(project.Members) != 1 || project.Members[0] != testOwnerID {
		t.Error("expected owner to be in members list")
	}
}

func TestCreateProjectEmptyName(t *testing.T) {
	repo := newMockRepository()
	svc := NewService(repo)

	_, err := svc.Create(context.Background(), CreateProjectRequest{Name: ""}, testOwnerID)

	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestUpdateProjectByOwner(t *testing.T) {
	repo := newMockRepository()
	repo.projects[testProjectID] = &models.Project{
		ID: testProjectID, Name: "Old", OwnerID: testOwnerID, Members: []string{testOwnerID},
	}
	svc := NewService(repo)

	newName := "New Name"
	project, err := svc.Update(context.Background(), testProjectID, UpdateProjectRequest{Name: &newName}, testOwnerID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if project.Name != "New Name" {
		t.Errorf("expected name 'New Name', got '%s'", project.Name)
	}
}

func TestUpdateProjectByNonOwner(t *testing.T) {
	repo := newMockRepository()
	repo.projects[testProjectID] = &models.Project{
		ID: testProjectID, OwnerID: testOwnerID, Members: []string{testOwnerID},
	}
	svc := NewService(repo)

	newName := "Hack"
	_, err := svc.Update(context.Background(), testProjectID, UpdateProjectRequest{Name: &newName}, "other-user")

	if err == nil {
		t.Fatal("expected forbidden error")
	}
}

func TestDeleteProjectByNonOwner(t *testing.T) {
	repo := newMockRepository()
	repo.projects[testProjectID] = &models.Project{
		ID: testProjectID, OwnerID: testOwnerID,
	}
	svc := NewService(repo)

	err := svc.Delete(context.Background(), testProjectID, "other-user")

	if err == nil {
		t.Fatal("expected forbidden error")
	}
}

func TestAddMember(t *testing.T) {
	repo := newMockRepository()
	repo.projects[testProjectID] = &models.Project{
		ID: testProjectID, OwnerID: testOwnerID, Members: []string{testOwnerID},
	}
	svc := NewService(repo)

	err := svc.AddMember(context.Background(), testProjectID, testMemberID, testOwnerID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(repo.projects[testProjectID].Members) != 2 {
		t.Error("expected 2 members after adding")
	}
}

func TestAddMemberDuplicate(t *testing.T) {
	repo := newMockRepository()
	repo.projects[testProjectID] = &models.Project{
		ID: testProjectID, OwnerID: testOwnerID, Members: []string{testOwnerID},
	}
	svc := NewService(repo)

	err := svc.AddMember(context.Background(), testProjectID, testOwnerID, testOwnerID)

	if err == nil {
		t.Fatal("expected conflict error for duplicate member")
	}
}

func TestAddMemberByNonOwner(t *testing.T) {
	repo := newMockRepository()
	repo.projects[testProjectID] = &models.Project{
		ID: testProjectID, OwnerID: testOwnerID, Members: []string{testOwnerID},
	}
	svc := NewService(repo)

	err := svc.AddMember(context.Background(), testProjectID, testMemberID, "other-user")

	if err == nil {
		t.Fatal("expected forbidden error")
	}
}

func TestRemoveMemberOwner(t *testing.T) {
	repo := newMockRepository()
	repo.projects[testProjectID] = &models.Project{
		ID: testProjectID, OwnerID: testOwnerID, Members: []string{testOwnerID, testMemberID},
	}
	svc := NewService(repo)

	err := svc.RemoveMember(context.Background(), testProjectID, testOwnerID, testOwnerID)

	if err == nil {
		t.Fatal("expected error when removing owner")
	}
}
