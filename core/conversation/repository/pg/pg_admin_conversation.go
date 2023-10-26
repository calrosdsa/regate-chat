package pg

import (
	"context"
	"database/sql"
	"log"
	r "message/domain/repository"
)


type conversationAdminRepo struct {
	Conn *sql.DB
}

func NewAdminRepository(conn *sql.DB) r.ConversationAdminRepository {
	return &conversationAdminRepo{
		Conn: conn,
	}
}

func (p *conversationAdminRepo) GetConversationsEstablecimiento(ctx context.Context,uuid string) (res []r.EstablecimientoConversation,err error) {
	query := `select p.nombre,p.apellido,p.profile_photo,c.id,p.profile_id,c.parent_id
	from chat as c 
	inner join profiles as p on c.second_parent_id = p.profile_id
	where c.parent_id = (select establecimiento_id from establecimientos where uuid = $1)`
	res,err = p.fetchConversationsAdmin(ctx,query,uuid)
	// query := `insert into conversation_message (id,conversation_id,sender_id,content,created_at,reply_to) 
	// values($1,$2,$3,$4,$5,$6) returning id,created_at`
	return
}

func (p *conversationAdminRepo) fetchConversationsAdmin(ctx context.Context, query string, args ...interface{}) (res []r.EstablecimientoConversation, err error) {
	rows, err := p.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			log.Println(errRow)
		}
	}()
	res = make([]r.EstablecimientoConversation, 0)
	for rows.Next() {
		t := r.EstablecimientoConversation{}
		err = rows.Scan(
			&t.Name,
			&t.Apellido,
			&t.Photo,
			&t.ConversationId,
			&t.ProfileId,
			&t.ParentId,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		res = append(res, t)
	}
	return res, nil
}