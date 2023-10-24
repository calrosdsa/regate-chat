package pg

import (
	"context"
	"database/sql"
	"log"
	r "message/domain/repository"
)

type chatRepo struct {
	Conn *sql.DB
}

func NewRepository(conn *sql.DB) r.ChatRepository {
	return &chatRepo{
		Conn: conn,
	}
}

func (p *chatRepo)GetChatsUser(ctx context.Context,profileId int,page int16,size int8)(res []r.Chat,
err error){
	// query := `select g.grupo_id,g.name,g.photo,gm.content,gm.created_at,
	// (select count(*) from grupo_message as gmc where gmc.grupo_id = g.grupo_id) as count
	// from user_grupo as ug left join lateral 
	// (select m.content,m.created_at from grupo_message as m where ug.grupo_id = m.grupo_id
	// order by created_at desc limit 1 ) gm on true
	// inner join grupos as g on g.grupo_id = ug.grupo_id where  ug.profile_id = $1
	// union all 
	// select c.conversation_id,e.name,e.photo,cm.content,cm.created_at,
	// (select count(*) from conversation_message as cmc where cmc.conversation_id = c.conversation_id) as count
	// from conversations as c left join lateral 
	// (select m.content,m.created_at from conversation_message as m 
	// where m.conversation_id = c.conversation_id order by created_at desc limit 1 ) cm on true
	// left join establecimientos as e on e.establecimiento_id = c.establecimiento_id
	// where c.profile_id = 43 
	// order by created_at desc limit $2 offset $3`
	query := `select c.id,g.name,g.photo,($4),ug.grupo_id from user_grupo as ug
	inner join grupos as g on g.grupo_id = ug.grupo_id
	inner join chat as c on c.parent_id = ug.grupo_id
	where  ug.profile_id = $1
	union all 
	select c.id,e.name,e.photo,($5),e.establecimiento_id from chat as c 
	left join establecimientos as e on e.establecimiento_id = c.parent_id
	where c.second_parent_id = $1
	limit $2 offset $3`
	res,err = p.fetchChats(ctx,query,profileId,size, page * int16(size),r.TypeChatGrupo,r.TypeChatInboxEstablecimiento)
	return
}


func (p *chatRepo) fetchChats(ctx context.Context, query string, args ...interface{}) (res []r.Chat, err error) {
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
	res = make([]r.Chat, 0)
	for rows.Next() {
		t := r.Chat{}
		err = rows.Scan(
			&t.Id,
			&t.Name,
			&t.Photo,
			// &t.LastMessage,
			// &t.LastMessageCreated,
			// &t.MessagesCount,
			&t.TypeChat,
			&t.ParentId,
			// &t.ProfileId,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		res = append(res, t)
	}
	return res, nil
}


