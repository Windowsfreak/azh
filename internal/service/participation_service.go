package service

import (
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"azh/internal/model"
	"azh/internal/repository"
)

// ParticipantDTO represents the data transfer object for participants
type ParticipantDTO struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Notes     string `json:"notes"`
	Present   bool   `json:"present"`
}

// ParticipationService handles business logic for participations
type ParticipationService struct {
	courseRepo        *repository.CourseRepository
	memberCourseRepo  *repository.MemberCourseRepository
	participationRepo *repository.ParticipationRepository
	memberRepo        *repository.MemberRepository
}

// NewParticipationService creates a new ParticipationService
func NewParticipationService(
	courseRepo *repository.CourseRepository,
	memberCourseRepo *repository.MemberCourseRepository,
	participationRepo *repository.ParticipationRepository,
	memberRepo *repository.MemberRepository,
) *ParticipationService {
	return &ParticipationService{
		courseRepo:        courseRepo,
		memberCourseRepo:  memberCourseRepo,
		participationRepo: participationRepo,
		memberRepo:        memberRepo,
	}
}

// GetParticipants retrieves participants for a specific course and date
func (s *ParticipationService) GetParticipants(courseID, date string) ([]ParticipantDTO, error) {
	// Parse the selected date
	selectedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %v", err)
	}

	// Get member IDs associated with the course
	memberIDs, err := s.memberCourseRepo.GetMembersByCourseID(courseID)
	if err != nil {
		return nil, err
	}

	// Fetch member details, filtered by sign_up_date and cancellation_date
	members, err := s.memberRepo.GetByIDsAndDate(memberIDs, selectedDate)
	if err != nil {
		return nil, err
	}

	// Get participation records for the course and date
	participations, err := s.participationRepo.GetByCourseAndDate(courseID, date)
	if err != nil {
		return nil, err
	}

	// Map participations by member ID for quick lookup
	participationMap := make(map[uint]struct{})
	for _, p := range participations {
		participationMap[p.MemberID] = struct{}{}
	}

	// Build participant DTOs
	participants := make([]ParticipantDTO, 0, len(members))
	for _, member := range members {
		_, present := participationMap[member.ID]
		participants = append(participants, ParticipantDTO{
			ID:        member.ID,
			FirstName: member.FirstName,
			LastName:  member.LastName,
			Phone:     member.Phone,
			Notes:     member.Notes,
			Present:   present,
		})
	}
	return participants, nil
}

// SetAttendance updates the attendance status for a participant
func (s *ParticipationService) SetAttendance(courseID uint, date time.Time, memberID uint, present bool) error {
	if present {
		participation := &model.Participation{
			MemberID: memberID,
			CourseID: courseID,
			Date:     date,
		}
		return s.participationRepo.Upsert(participation)
	} else {
		return s.participationRepo.Delete(memberID, courseID, date)
	}
}

// ExportData exports participation data within a date range as CSV
func (s *ParticipationService) ExportData(minDate, maxDate string) (string, error) {
	participations, err := s.participationRepo.GetExportData(minDate, maxDate)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	writer := csv.NewWriter(&builder)

	// Write CSV header
	header := []string{"Date", "CourseID", "MemberID"}
	if err := writer.Write(header); err != nil {
		return "", err
	}

	// Write data rows
	for _, p := range participations {
		row := []string{
			p.Date.Format("2006-01-02"),
			fmt.Sprintf("%d", p.CourseID),
			fmt.Sprintf("%d", p.MemberID),
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	writer.Flush()
	return builder.String(), nil
}
