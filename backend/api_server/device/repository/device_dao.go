package repository

import (
	"context"
	"fmt"
	"sync"

	"api_server/ent"
	"api_server/ent/device"
	"api_server/ent/gpu"
	"api_server/logger"
	"api_server/utils"

	"entgo.io/ent/dialect/sql"
)

type IDeviceDAO interface {
	InsertOne(ctx context.Context, d DeviceDTO) (*ent.Device, error)
	SelectAll(ctx context.Context) ([]*ent.Device, error)
	SelectActive(ctx context.Context) ([]*ent.Device, error)
	SelectOne(ctx context.Context, device_id int) (*ent.Device, error)
	SelectIdleByGPU(ctx context.Context, gpu_ids []int) ([]*ent.Device, error)
	SelectByGPU(ctx context.Context, gpu_ids []int) ([]*ent.Device, error)

	// SelectByIPAndPort는 지정된 IP와 Port 값을 가지는 Device를 조회합니다.
	//
	// 매개변수:
	//   - ctx: context.Context
	//   - ip: 조회할 IP 주소
	//   - port: 조회할 포트 번호
	//
	// 반환값:
	//   - []*ent.Device: 조회된 Device 리스트
	//   - error: 오류 정보
	SelectByIPAndPort(ctx context.Context, ip string, port int) ([]*ent.Device, error)
	UpdateOne(ctx context.Context, d DeviceDTO) (*ent.Device, error)
	DeleteMany(ctx context.Context, ids []int) (int, error)
	DeleteOne(ctx context.Context, devie_id int) error
}

type DeviceDAO struct {
	dbms *ent.Client
}

var once sync.Once
var instance *DeviceDAO

func New() *DeviceDAO {
	once.Do(func() {
		logger.Debug("Device DAO instance")
		instance = &DeviceDAO{
			dbms: utils.GetEntClient(),
		}
	})

	return instance
}

func (dao *DeviceDAO) InsertOne(ctx context.Context, req DeviceDTO) (*ent.Device, error) {
	logger.Debug(fmt.Sprintf("%+v", req))
	return dao.dbms.Device.Create().
		SetName(req.Name).
		SetIP(req.IP).
		SetPort(*req.Port).
		SetIsUse(true).
		SetType("engine").
		Save(ctx)
}

func (dao *DeviceDAO) SelectAll(ctx context.Context) ([]*ent.Device, error) {
	logger.Debug("Select All devices")
	return dao.dbms.Device.Query().
		Order(device.ByName(sql.OrderAsc())).
		All(ctx)
}

func (dao *DeviceDAO) SelectActive(ctx context.Context) ([]*ent.Device, error) {
	logger.Debug("Select Active devices")
	return dao.dbms.Device.
		Query().
		Where(device.IsUse(true)).
		Order(device.ByName(sql.OrderAsc())).
		All(ctx)
}

func (dao *DeviceDAO) SelectOne(ctx context.Context, device_id int) (*ent.Device, error) {
	logger.Debug(fmt.Sprintf(`{"device_id": %d}`, device_id))
	return dao.dbms.Device.Query().
		Where(device.ID(device_id)).
		Only(ctx)
}

func (dao *DeviceDAO) SelectIdleByGPU(ctx context.Context, gpu_ids []int) ([]*ent.Device, error) {
	logger.Debug(fmt.Sprintf(`{"gpu_ids": %d}`, gpu_ids))
	/*
		select distinct s.*
		from device s
		join gpu as t on t.device_id = s.id
		where t.id in (${gpu_ids})
	*/
	return dao.dbms.Device.Query().
		Where(func(s *sql.Selector) {
			t := sql.Table(gpu.Table)
			s.Join(t).On(s.C(device.FieldID), t.C(gpu.FieldDeviceID))
			s.Where(sql.And(
				sql.InInts(t.C(gpu.FieldID), gpu_ids...),
				sql.EQ(t.C(gpu.FieldState), utils.GPU_STATE_IDLE),
			))
		}).Unique(true).
		// WithGpu(func(g *ent.GpuQuery) {
		// 	g.Select(gpu.FieldID)
		// }).
		All(ctx)
}

func (dao *DeviceDAO) SelectByGPU(ctx context.Context, gpu_ids []int) ([]*ent.Device, error) {
	logger.Debug(fmt.Sprintf(`{"gpu_ids": %d}`, gpu_ids))
	return dao.dbms.Device.Query().
		Where(func(s *sql.Selector) {
			t := sql.Table(gpu.Table)
			s.Join(t).On(s.C(device.FieldID), t.C(gpu.FieldDeviceID))
			s.Where(sql.InInts(t.C(gpu.FieldID), gpu_ids...))
		}).Unique(true).
		All(ctx)
}

func (dao *DeviceDAO) UpdateOne(ctx context.Context, req DeviceDTO) (*ent.Device, error) {
	logger.Debug(fmt.Sprintf("%+v", req))
	return dao.dbms.Device.UpdateOneID(req.ID).
		SetNillableName(&req.Name).
		SetNillableIP(&req.IP).
		SetNillablePort(req.Port).
		SetNillableIsUse(req.IsUse).
		SetType(req.Type).
		SetConnection(req.Connection).
		Save(ctx)
}

func (dao *DeviceDAO) DeleteMany(ctx context.Context, ids []int) (int, error) {
	logger.Debug(fmt.Sprintf("%+v", ids))
	return dao.dbms.Device.Delete().
		Where(device.IDIn(ids...)).
		Exec(ctx)
}

func (dao *DeviceDAO) DeleteOne(ctx context.Context, devie_id int) error {
	logger.Debug(fmt.Sprintf(`{"devie_id": %d}`, devie_id))
	return dao.dbms.Device.
		DeleteOneID(devie_id).
		Exec(ctx)
}

// SelectByIPAndPort는 지정된 IP와 Port 값을 가지는 Device를 조회합니다.
func (dao *DeviceDAO) SelectByIPAndPort(ctx context.Context, ip string, port int) ([]*ent.Device, error) {
	logger.Debug(fmt.Sprintf("Select devices by IP: %s and Port: %d", ip, port))
	return dao.dbms.Device.
		Query().
		Where(
			device.IP(ip),
			device.Port(port),
		).
		Order(device.ByName(sql.OrderAsc())).
		All(ctx)
}
