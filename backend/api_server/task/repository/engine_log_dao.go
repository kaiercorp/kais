package repository

import (
	"context"
	"fmt"

	"api_server/ent"
	"api_server/logger"
	"api_server/utils"
)

type EngineLogDAO struct {
	dbms *ent.Client
	ctx  context.Context
}

type IEngineLogDAO interface {
	selectEngineLog(modeling_id int, model_name string) (*EngineLogResponse, error)
}

func NewEngineLogDAO() *EngineLogDAO {
	return &EngineLogDAO{
		dbms: utils.GetEntClient(),
		ctx:  context.Background(),
	}
}

func (dao *EngineLogDAO) SelectEngineLog(modeling_id int) (*[]EngineLogResponse, error) {
	logger.Debug(fmt.Sprintf(`{"modeling_id": %d}`, modeling_id))

	rows, err := dao.dbms.QueryContext(
		dao.ctx,
		fmt.Sprintf(`select  el.created_at, el.level, el.line, el.message, m.modeling_step from enginelog el 
		join modeling m on m.id = el.modeling_id
		where modeling_id = %d ;
		`,
			modeling_id),
	)

	if err != nil {
		return nil, err
	}

	results := []EngineLogResponse{}
	for rows.Next() {
		result := EngineLogResponse{}
		if err := rows.Scan(
			&result.CreatedAt, &result.Level, &result.Line, &result.Message, &result.ModelingStep,
		); err != nil {
			fmt.Println(err)
			continue
		}

		results = append(results, result)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("empty engine log results")
	}
	return &results, nil
}
