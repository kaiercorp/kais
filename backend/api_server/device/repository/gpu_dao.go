package repository

import (
	"api_server/ent"
	"api_server/ent/device"
	"api_server/ent/gpu"
	"api_server/logger"
	"api_server/utils"
	"context"
	"fmt"
	"strconv"
	"sync"

	"entgo.io/ent/dialect/sql"
)

type IGPUDAO interface {
	UpsertMany(ctx context.Context, req EngineInfoDTO) error
	SelectIsUse(ctx context.Context) ([]*ent.Device, error)
	SelectIdle(ctx context.Context) ([]*ent.Gpu, error)
	SelectManyByDeviceID(ctx context.Context, deviceID int) ([]*ent.Gpu, error)
	UpdateAllDisUse(ctx context.Context) error
	UpdateManyState(ctx context.Context, ids []int, state string) error
	SelectGpuByIndex(ctx context.Context, gpuIndex string) (*ent.Gpu, error)
}

type GPUDAO struct {
	dbms *ent.Client
}

var onceGPU sync.Once
var instanceGPU *GPUDAO

func NewGPUDAO() *GPUDAO {
	onceGPU.Do(func() {
		logger.Debug("GPU DAO instance")
		instanceGPU = &GPUDAO{
			dbms: utils.GetEntClient(),
		}
	})

	return instanceGPU
}

func (dao *GPUDAO) UpsertMany(ctx context.Context, req EngineInfoDTO) error {
	logger.Debug(fmt.Sprintf("%+v", req))
	return dao.dbms.Gpu.MapCreateBulk(
		req.GPUs,
		func(c *ent.GpuCreate, i int) {
			c.
				SetDeviceID(req.DeviceID).
				SetIndex(req.GPUs[i].Index).
				SetName(req.GPUs[i].Name).
				SetUUID(req.GPUs[i].UUID).
				SetState(utils.GPU_STATE_IDLE)
		},
	).
		OnConflict(
			sql.ConflictColumns(gpu.FieldUUID),
		).
		UpdateIsUse().
		Exec(ctx)
}

func (dao *GPUDAO) SelectIsUse(ctx context.Context) ([]*ent.Device, error) {
	logger.Debug("Select Usable GPUs")
	return dao.dbms.Device.Query().
		Where(device.IsUse(true)).
		WithGpu(func(g *ent.GpuQuery) {
			g.Where(gpu.IsUse(true))
		}).
		Order(device.ByID(sql.OrderAsc())).
		All(ctx)
}

func (dao *GPUDAO) SelectIdle(ctx context.Context) ([]*ent.Gpu, error) {
	logger.Debug("Select Idle GPUs")
	return dao.dbms.Gpu.Query().
		Where(
			gpu.And(
				gpu.IsUse(true),
				gpu.State(utils.GPU_STATE_IDLE),
			),
		).
		All(ctx)
}

func (dao *GPUDAO) UpdateAllDisUse(ctx context.Context) error {
	logger.Debug("Update all gpus to disuse")
	return dao.dbms.Gpu.Update().
		SetState(utils.GPU_STATE_IDLE).
		SetIsUse(false).
		Exec(ctx)
}

func (dao *GPUDAO) UpdateOneState(ctx context.Context, gpu_id int, state string) error {
	logger.Debug(fmt.Sprintf(`{"gpu_id": %d, "state": %s}`, gpu_id, state))
	return dao.dbms.Gpu.UpdateOneID(gpu_id).
		SetState(state).
		Exec(ctx)
}

func (dao *GPUDAO) UpdateManyState(ctx context.Context, ids []int, state string) error {
	logger.Debug(fmt.Sprintf(`{"ids": %+v, "state": %s}`, ids, state))
	return dao.dbms.Gpu.Update().
		SetState(state).
		Where(gpu.IDIn(ids...)).
		Exec(ctx)
}

func (dao *GPUDAO) SelectManyByDeviceID(ctx context.Context, deviceID int) ([]*ent.Gpu, error) {
	logger.Debug(fmt.Sprintf("SelectManyByDeviceID: %d", deviceID))

	return dao.dbms.Gpu.Query().
		Where(gpu.DeviceID(deviceID)).
		Order(gpu.ByIndex(sql.OrderAsc())).
		All(ctx)
}

func (dao *GPUDAO) SelectGpuByIndex(ctx context.Context, gpuIndex string) (*ent.Gpu, error) {
	logger.Debug(fmt.Sprintf(`{"index": %v}`, gpuIndex))
	gpuIndexInt, err := strconv.Atoi(gpuIndex)
	if err != nil {
		logger.Error("failed to convert string to int")
		return nil, err
	}

	return dao.dbms.Gpu.Query().
		Where(gpu.Index(gpuIndexInt), gpu.IsUse(true)).
		Only(ctx)
}
