package job

import (
	"context"
	"fmt"

	"github.com/swaindhruti/pharmastock-backend/internal/inventory"
	"github.com/swaindhruti/pharmastock-backend/internal/medicine"
)

type Processor interface {
	Process(ctx context.Context, job *Job) error
}

type processor struct {
	medicineRepo  medicine.Repository
	inventoryRepo inventory.Repository
}

func NewProcessor(medicineRepo medicine.Repository, inventoryRepo inventory.Repository) Processor {
	return &processor{
		medicineRepo:  medicineRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (p *processor) Process(ctx context.Context, job *Job) error {
	parsed, err := medicine.ParseFile(job.FilePath)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	if len(parsed) == 0 {
		return nil
	}

	medicines := make([]*medicine.Medicine, len(parsed))
	for i, entry := range parsed {
		medicines[i] = &medicine.Medicine{Name: entry.Name}
	}

	if err := p.medicineRepo.BatchInsert(ctx, medicines); err != nil {
		return fmt.Errorf("failed to batch insert medicines: %w", err)
	}

	names := make([]string, len(parsed))
	for i, entry := range parsed {
		names[i] = entry.Name
	}

	saved, err := p.medicineRepo.GetMedicinesByNames(ctx, names)
	if err != nil {
		return fmt.Errorf("failed to get medicines by names: %w", err)
	}

	nameToID := make(map[string]int64, len(saved))
	for _, m := range saved {
		nameToID[m.Name] = m.ID
	}

	var entries []inventory.Entry
	for _, entry := range parsed {
		if id, ok := nameToID[entry.Name]; ok {
			entries = append(entries, inventory.Entry{
				StockistID: job.StockistID,
				MedicineID: id,
			})
		}
	}

	if len(entries) > 0 {
		if err := p.inventoryRepo.BulkCreate(ctx, entries); err != nil {
			return fmt.Errorf("failed to create inventory entries: %w", err)
		}
	}

	return nil
}
