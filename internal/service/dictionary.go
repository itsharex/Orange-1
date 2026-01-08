package service

import (
	"github.com/FruitsAI/Orange/internal/models"
	"github.com/FruitsAI/Orange/internal/repository"
)

// DictionaryService 字典服务
type DictionaryService struct {
	dictRepo *repository.DictionaryRepository
}

// NewDictionaryService 创建字典服务
func NewDictionaryService() *DictionaryService {
	return &DictionaryService{
		dictRepo: repository.NewDictionaryRepository(),
	}
}

// List 获取字典列表
func (s *DictionaryService) List() ([]models.Dictionary, error) {
	return s.dictRepo.List()
}

// GetItems 获取字典项
func (s *DictionaryService) GetItems(code string) ([]models.DictionaryItem, error) {
	return s.dictRepo.GetItemsByCode(code)
}

// CreateItem 创建字典项
func (s *DictionaryService) CreateItem(code, label, value string, sort int) (*models.DictionaryItem, error) {
	dict, err := s.dictRepo.FindByCode(code)
	if err != nil {
		return nil, err
	}

	item := &models.DictionaryItem{
		DictionaryID: dict.ID,
		Label:        label,
		Value:        value,
		Sort:         sort,
		Status:       1,
	}

	if err := s.dictRepo.CreateItem(item); err != nil {
		return nil, err
	}

	return item, nil
}

// UpdateItem 更新字典项
func (s *DictionaryService) UpdateItem(id int64, label, value string, sort int) (*models.DictionaryItem, error) {
	// First find the item to ensure it exists and get other fields if necessary
	// Since we don't have a direct FindItemByID exposed in repo for just one item easily without code context (repo.GetItems returns list),
	// but we can trust GORM's Save or use a direct update.
	// For better practice, let's fetch it first if we want to be safe, but repo doesn't expose FindItemByID.
	// We can add FindItemByID to repo or just use an update with map or struct.
	// Let's implement FindItemByID in repo if needed, but for now let's assume we can pass a struct with ID.
	// Actually, wait, GORM Save update all fields.
	// Let's modify repo to add FindItemByID, or better, just use db.First in repo's UpdateItem if check is needed.
	// But let's assume valid ID for now or user repository's update.

	// A better approach: We need to find the item first to make sure we don't overwrite other fields like DictionaryID or CreateTime with zero values if we create a new struct.
	// However, the repo UpdateItem takes a *models.DictionaryItem.
	// Let's rely on the handler to pass valid data or adds a GetItem method in Repo.
	// Checking existing repo methods... `GetItems` returns a list. `DeleteItem` deletes by ID.
	// I should probably add `FindItemByID` to repo for safety, but to keep it simple and minimal changes as per request:
	// I can perform an update with specific fields using map in repo, OR
	// I'll assume I can just update specific columns if I structured it that way, but `r.db.Save` replaces everything.
	// I'll add `FindItemByID` to repo in a separate step if I was being very strict, but `s.dictRepo.UpdateItem` expects a model.
	// Let's check `internal/repository/dictionary.go` again.
	// It has `DeleteItem(id int64)`.
	// I will add a simple FindItemByID to repo in the same file first? Or just implement a "UpdateItemFields" in repo.
	// Let's try to do it cleanly. I'll stick to adding `UpdateItem` that takes a struct, and in service I should ideally fetch it first.
	// LIMITATION: I don't have `FindItemByID` in repo.
	// I'll assume valid input for now or use a raw update in repo if I could change it.
	// Let's add `FindItemByID` in repo? No, user didn't ask explicitly.
	// I'll assume I can construct the object. But wait, if I construct a new object with only ID, Label, Value, Sort, other fields might be zeroed out if I use Save.
	// `r.db.Save` is `INSERT ON DUPLICATE KEY UPDATE` or `UPDATE` depending on ID.
	// If I only populate ID, Label, Value, Sort, then DictionaryID might become 0! This is bad.
	// SO I NEED `FindItemByID` in repo.
	// I will update the plan/execution to include adding `FindItemByID` to repo or specific update.
	// Let's look at `internal/repository/dictionary.go` again.
	// It does not have `FindItemByID`.
	// I will modify `internal/repository/dictionary.go` to include `FindItemByID` as well in the next step or this one if I can.
	// No, I'll allow `UpdateItem` in service to take ID and call a new Repo method that updates specific columns.

	// Revised plan for backend:
	// 1. Repo: Add `UpdateItem(item)` (Done). Wait, I used `Save`. I should have used `Updates` or verified.
	// Let's stick thereto. I will need to fetch the item first.
	// I will add `FindItemByID` to repository in the NEXT step (since I already sent one replace_file_content).
	// Actually, I can do it in the Service if I had access to DB, but I don't.
	// So:
	// Step 1 (Done): Added `UpdateItem` (Save).
	// Step 2: Add `FindItemByID` to Repo (I missed this in plan, but it's needed).
	// Step 3: Service uses `FindItemByID` then `UpdateItem`.

	// Let's proceed with adding `UpdateItem` to Service first, assuming I'll fix the repo dependency in a sec.
	// Wait, simpler: I'll modify Repo to `UpdateItem` using `db.Model(&item).Updates(...)`. using a map or struct.
	// If I use `Updates`, I don't need to fetch full object if I just want to update fields.
	// But `user.id` checks? DictionaryID logic?
	// Let's keep it simple. Service will just update.

	// Correct logic for Service UpdateItem:
	// It needs to retrieve the item to ensure it exists and preserve DictionaryID.
	// Since I cannot modify repo in the same tool call easily without conflict or complex regex.
	// I'll skip adding `FindItemByID` and just add a direct `UpdateItemColumns` in repo?
	// Or simpler: I will add `FindItemByID` to repo now.

	// Actually, I already sent the tool call for `UpdateItem` in repo.
	// Let's edit `internal/repository/dictionary.go` AGAIN to add `FindItemByID`.

	// For now, let's implement Service UpdateItem and assume I'll add the necessary Repo method.
	// Wait, I can't write code that calls a non-existent method.
	// I'll write the service code to call `dictRepo.FindItemByID` which I will implement.

	item, err := s.dictRepo.FindItemByID(id)
	if err != nil {
		return nil, err
	}

	item.Label = label
	item.Value = value
	item.Sort = sort

	if err := s.dictRepo.UpdateItem(item); err != nil {
		return nil, err
	}
	return item, nil
}

// DeleteItem 删除字典项
func (s *DictionaryService) DeleteItem(id int64) error {
	return s.dictRepo.DeleteItem(id)
}
