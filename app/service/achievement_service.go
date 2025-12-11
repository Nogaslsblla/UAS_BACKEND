wpackage service

import (
	"context"
	"errors"
	"time"
	mongodb "uas_backend/model/MongoDB"
	model "uas_backend/model/Postgresql"
	"uas_backend/repository"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementService interface {
	SubmitPrestasi(ctx context.Context, userID uuid.UUID, req mongodb.Achievement) (*model.AchievementReference, error)
	SubmitForVerification(ctx context.Context, userID uuid.UUID, achievementID uuid.UUID) (*model.AchievementReference, error)
}

type achievementService struct {
	repo repository.AchievementRepository
}

func NewAchievementService(repo repository.AchievementRepository) AchievementService {
	return &achievementService{repo: repo}
}

func (s *achievementService) SubmitPrestasi(ctx context.Context, userID uuid.UUID, req mongodb.Achievement) (*model.AchievementReference, error) {
	// 1. Cari data Student berdasarkan User ID yang login
	student, err := s.repo.GetStudentByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("student data not found for this user")
	}

	// 2. Setup Data untuk MongoDB
	req.ID = primitive.NewObjectID()
	req.StudentID = student.ID // Link ke UUID Student di Postgres
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	if req.CustomFields == nil {
		req.CustomFields = make(map[string]interface{})
	}

	// 3. Simpan ke MongoDB
	mongoID, err := s.repo.SaveAchievementMongo(ctx, req)
	if err != nil {
		return nil, err
	}

	// 4. Setup Data untuk Postgres (Reference)
	// Sesuai SRS Flow 4: Status awal 'draft'
	ref := model.AchievementReference{
		ID:                 uuid.New(),
		StudentID:          student.ID,
		MongoAchievementID: mongoID,
		Status:             "draft",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// 5. Simpan ke Postgres
	err = s.repo.SaveAchievementReference(ctx, ref)
	if err != nil {
		return nil, err
	}

	return &ref, nil
}

// SubmitForVerification - FR-004: Submit untuk Verifikasi
func (s *achievementService) SubmitForVerification(ctx context.Context, userID uuid.UUID, achievementID uuid.UUID) (*model.AchievementReference, error) {
	// 1. Cari data Student berdasarkan User ID yang login
	student, err := s.repo.GetStudentByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("student data not found for this user")
	}

	// 2. Get achievement reference by ID
	ref, err := s.repo.GetAchievementReferenceByID(ctx, achievementID)
	if err != nil {
		return nil, errors.New("achievement not found")
	}

	// 3. Validasi: Pastikan achievement milik student yang login
	if ref.StudentID != student.ID {
		return nil, errors.New("unauthorized: achievement does not belong to this student")
	}

	// 4. Validasi: Pastikan status adalah 'draft'
	if ref.Status != "draft" {
		return nil, errors.New("achievement must be in 'draft' status to submit")
	}

	// 5. Update status menjadi 'submitted'
	err = s.repo.UpdateAchievementStatusToSubmitted(ctx, achievementID)
	if err != nil {
		return nil, errors.New("failed to update achievement status")
	}

	// 6. Get updated achievement reference
	updatedRef, err := s.repo.GetAchievementReferenceByID(ctx, achievementID)
	if err != nil {
		return nil, err
	}

	// 7. Get advisor ID untuk notifikasi (opsional, bisa digunakan untuk logging)
	advisorID, err := s.repo.GetAdvisorIDByStudentID(ctx, student.ID)
	if err == nil {
		// Log notifikasi ke dosen wali (bisa dikembangkan lebih lanjut)
		// Untuk saat ini, advisor_id sudah tersimpan di student record
		_ = advisorID // Placeholder untuk future notification system
	}

	return updatedRef, nil
}
