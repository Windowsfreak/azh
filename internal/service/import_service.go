package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"azh/internal/model"
	"azh/internal/repository"
	"gorm.io/gorm"
)

// ImportService handles CSV import logic
type ImportService struct {
	db                *gorm.DB
	courseRepo        *repository.CourseRepository
	memberRepo        *repository.MemberRepository
	memberCourseRepo  *repository.MemberCourseRepository
	participationRepo *repository.ParticipationRepository
}

// NewImportService creates a new ImportService
func NewImportService(
	db *gorm.DB,
	courseRepo *repository.CourseRepository,
	memberRepo *repository.MemberRepository,
	memberCourseRepo *repository.MemberCourseRepository,
	participationRepo *repository.ParticipationRepository,
) *ImportService {
	return &ImportService{
		db:                db,
		courseRepo:        courseRepo,
		memberRepo:        memberRepo,
		memberCourseRepo:  memberCourseRepo,
		participationRepo: participationRepo,
	}
}

// ProcessCSV processes a CSV file based on detected type
func (s *ImportService) ProcessCSV(filePath, fileName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	// Read header to detect file type
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("unable to read header: %v", err)
	}

	// Detect file type based on column presence
	isParticipants := false
	isTrainings := false
	for _, h := range header {
		h = strings.TrimSpace(h)
		if strings.EqualFold(h, "Alter") {
			isParticipants = true
		} else if strings.EqualFold(h, "Trainer") {
			isTrainings = true
		}
	}

	if isParticipants && !isTrainings {
		return s.importParticipants(filePath, header, reader)
	} else if isTrainings && !isParticipants {
		return s.importCourses(filePath, header, reader)
	} else {
		// Fallback to filename hint
		if strings.Contains(strings.ToLower(fileName), "trainingsstatistik") {
			return s.importCourses(filePath, header, reader)
		} else if strings.Contains(strings.ToLower(fileName), "trainingsanmeldungen") {
			return s.importParticipants(filePath, header, reader)
		}
		return fmt.Errorf("unable to determine file type for: %s", fileName)
	}
}

// importCourses imports course data from TrainingsStatistik.csv
func (s *ImportService) importCourses(filePath string, header []string, reader *csv.Reader) error {
	// Map header to column indices
	idIdx := 0   // First column assumed as ID
	nameIdx := 1 // Second column assumed as Name
	ortIdx := -1
	trainerIdx := -1
	sparteIdx := -1
	wochentagIdx := -1
	startIdx := -1
	endeIdx := -1
	//letzterTerminIdx := -1
	for i, h := range header {
		h = strings.TrimSpace(h)
		switch {
		case strings.EqualFold(h, "Ort"):
			ortIdx = i
		case strings.EqualFold(h, "Trainer"):
			trainerIdx = i
		case strings.EqualFold(h, "Sparte"):
			sparteIdx = i
		case strings.EqualFold(h, "Wochentag"):
			wochentagIdx = i
		case strings.EqualFold(h, "Start"):
			startIdx = i
		case strings.EqualFold(h, "Ende"):
			endeIdx = i
			//case strings.EqualFold(h, "letzter Termin"):
			//	letzterTerminIdx = i
		}
	}

	// Read and process rows
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading row: %v", err)
		}

		// Skip empty or summary rows (e.g., "Gesamt")
		if len(row) < idIdx+1 || strings.TrimSpace(row[idIdx]) == "" || strings.Contains(strings.ToLower(row[idIdx]), "gesamt") {
			continue
		}

		// Extract data
		var courseID uint
		if _, err := fmt.Sscanf(row[idIdx], "%d", &courseID); err != nil {
			continue // Skip rows with invalid course ID
		}
		name := row[nameIdx]
		location := safeGet(row, ortIdx)
		trainerNames := safeGet(row, trainerIdx)
		trainingType := safeGet(row, sparteIdx)
		weekday := safeGet(row, wochentagIdx)
		startTime := safeGet(row, startIdx)
		endTime := safeGet(row, endeIdx)
		//lastScheduleStr := safeGet(row, letzterTerminIdx)

		// Parse last schedule date if provided (format DD.MM.YYYY)
		var lastSchedule = time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)
		//if lastScheduleStr != "" {
		//	if parsed, err := time.Parse("02.01.2006", lastScheduleStr); err == nil {
		//		lastSchedule = parsed
		//	}
		//}

		// Upsert course
		course := model.Course{
			ID:           courseID,
			Name:         name,
			Location:     location,
			TrainingType: trainingType,
			Weekday:      weekday,
			StartTime:    startTime,
			EndTime:      endTime,
			LastSchedule: lastSchedule,
			TrainerNames: trainerNames,
		}
		if err := s.db.Save(&course).Error; err != nil {
			return fmt.Errorf("error saving course %d: %v", courseID, err)
		}
	}
	return nil
}

// importParticipants imports participant and enrollment data from Trainingsanmeldungen.csv
func (s *ImportService) importParticipants(filePath string, header []string, reader *csv.Reader) error {
	// Map header to column indices
	datumIdx := -1
	kundigungsdatumIdx := -1
	mitgliedsnummerIdx := -1
	vornameIdx := -1
	nachnameIdx := -1
	telefonIdx := -1
	mitteilungIdx := -1
	notizenIdx := -1
	alterIdx := -1
	emailIdx := -1
	kursIdIdx := -1
	for i, h := range header {
		h = strings.TrimSpace(h)
		switch {
		case strings.EqualFold(h, "Datum"):
			datumIdx = i
		case strings.EqualFold(h, "Kündigungsdatum"):
			kundigungsdatumIdx = i
		case strings.EqualFold(h, "Mitgliedsnummer"):
			mitgliedsnummerIdx = i
		case strings.Contains(strings.ToLower(h), "vorname"):
			vornameIdx = i
		case strings.Contains(strings.ToLower(h), "nachname"):
			nachnameIdx = i
		case strings.EqualFold(h, "Erreichbarkeit per Telefon"):
			telefonIdx = i
		case strings.EqualFold(h, "Mitteilung"):
			mitteilungIdx = i
		case strings.EqualFold(h, "Notizen Büro"):
			notizenIdx = i
		case strings.EqualFold(h, "Alter"):
			alterIdx = i
		case strings.EqualFold(h, "E-Mail-Adresse"):
			emailIdx = i
		case strings.EqualFold(h, "Kurs Id"):
			kursIdIdx = i
		}
	}

	// Collect data for batch processing
	membersMap := make(map[uint]model.Member)
	memberCoursesSet := make(map[string]model.MemberCourse)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading row: %v", err)
		}

		// Skip empty rows
		if len(row) < mitgliedsnummerIdx+1 || strings.TrimSpace(row[mitgliedsnummerIdx]) == "" {
			continue
		}

		// Parse member ID
		var memberID uint
		if _, err := fmt.Sscanf(row[mitgliedsnummerIdx], "%d", &memberID); err != nil {
			continue // Skip rows with invalid member ID
		}

		// Parse dates (format may vary, try DD.MM.YYYY or full timestamp)
		signUpDateStr := safeGet(row, datumIdx)
		cancellationDateStr := safeGet(row, kundigungsdatumIdx)
		var signUpDate, cancellationDate time.Time
		signUpDate = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
		cancellationDate = time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)
		if signUpDateStr != "" {
			if parsed, err := parseDate(signUpDateStr); err == nil {
				signUpDate = parsed
			}
		}
		if cancellationDateStr != "" {
			if parsed, err := parseDate(cancellationDateStr); err == nil {
				cancellationDate = parsed
			}
		}

		// Parse age
		var age int
		ageStr := safeGet(row, alterIdx)
		if ageStr != "" {
			fmt.Sscanf(ageStr, "%d", &age)
		}

		// Combine notes fields
		notes := strings.TrimSpace(safeGet(row, mitteilungIdx) + " " + safeGet(row, notizenIdx))

		// Update members map (latest data wins for duplicates)
		membersMap[memberID] = model.Member{
			ID:               memberID,
			FirstName:        safeGet(row, vornameIdx),
			LastName:         safeGet(row, nachnameIdx),
			Email:            safeGet(row, emailIdx),
			Phone:            safeGet(row, telefonIdx),
			SignUpDate:       signUpDate,
			CancellationDate: cancellationDate,
			Age:              age,
			Notes:            notes,
		}

		// Parse course ID directly from Kurs Id column
		var courseID uint
		if _, err := fmt.Sscanf(row[kursIdIdx], "%d", &courseID); err != nil {
			continue // Skip rows with invalid course ID
		}

		// Add to member_courses set (deduplicate)
		key := fmt.Sprintf("%d-%d", memberID, courseID)
		memberCoursesSet[key] = model.MemberCourse{
			MemberID: memberID,
			CourseID: courseID,
		}
	}

	// Batch update database
	// 1. Upsert members
	for _, member := range membersMap {
		if err := s.db.Save(&member).Error; err != nil {
			return fmt.Errorf("error saving member %d: %v", member.ID, err)
		}
	}

	// 2. Truncate and reimport member_courses
	if err := s.db.Exec("TRUNCATE TABLE member_courses").Error; err != nil {
		return fmt.Errorf("error truncating member_courses: %v", err)
	}
	for _, mc := range memberCoursesSet {
		if err := s.db.Save(&mc).Error; err != nil {
			return fmt.Errorf("error saving member_course %d-%s: %v", mc.MemberID, mc.CourseID, err)
		}
	}

	return nil
}

// safeGet retrieves a value from a slice safely
func safeGet(row []string, index int) string {
	if index >= 0 && index < len(row) {
		return strings.TrimSpace(row[index])
	}
	return ""
}

// parseDate attempts to parse a date string in multiple formats
func parseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"02.01.2006",
		"02.01.2006 15:04:05",
	}
	for _, format := range formats {
		if parsed, err := time.Parse(format, dateStr); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}
