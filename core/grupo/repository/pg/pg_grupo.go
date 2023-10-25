package pg

import (
	"context"
	"database/sql"
	"log"

	// "log"
	r "message/domain/repository"
)

type grupoRepo struct {
	Conn *sql.DB
}

func NewRepository(conn *sql.DB) r.GrupoRepository {
	return &grupoRepo{
		Conn: conn,
	}
}

func (p *grupoRepo)GetUnreadMessages(ctx context.Context,profileId int,page int16,
	size int8)(res []r.MessageGrupo,err error){
	query := `select  gm.id,gm.chat_id,gm.profile_id,gm.content,gm.data,gm.created_at,gm.reply_to,gm.type_message
	from user_grupo as ug inner join grupo_message as gm on gm.grupo_id = ug.grupo_id 
	and ug.last_update_messages <= gm.created_at where ug.profile_id = $1 
	limit $2 offset $3`
	res,err = p.fetchMessagesGrupo(ctx,query,profileId,size,page *int16(size))
	return
	
}
func (p *grupoRepo)GetChatUnreadMessage(ctx context.Context,chatId int64,lastUpdated string )(res []r.MessageGrupo,err error){
	query := `select gm.id,gm.chat_id,gm.profile_id,gm.content,gm.data,gm.created_at,gm.reply_to,gm.type_message
	 from grupo_message as gm where chat_id = $1 and gm.created_at >= '2023-10-25T02:37:06.264Z'`
	res,err = p.fetchMessagesGrupo(ctx,query,chatId,lastUpdated)
	return

	select gm.id,gm.chat_id,gm.profile_id,gm.content,gm.data,gm.created_at,gm.reply_to,gm.type_message
	 from grupo_message as gm where chat_id = 1 and gm.created_at >= $2
}

func(p *grupoRepo)UpdateUserGrupoLastTimeUpdateMessage(ctx context.Context,profileId int)(err error){
	query :=`update user_grupo set last_update_messages = current_timestamp where profile_id = $1`
	_,err = p.Conn.ExecContext(ctx,query,profileId)
	return 
}

func (p *grupoRepo) SaveGrupoMessage(ctx context.Context, d *r.MessageGrupo) (err error) {
	log.Println(d.CreatedAt,"CreatedAt Message")
	query := `insert into grupo_message (id,chat_id,profile_id,content,created_at,reply_to,type_message,data,grupo_id) 
	values($1,$2,$3,$4,current_timestamp,$5,$6,$7,$8) returning id,created_at`
	err = p.Conn.QueryRowContext(ctx, query, d.Id, d.ChatId, d.ProfileId, d.Content, d.ReplyTo,
		d.TypeMessage, d.Data,d.ParentId).Scan(&d.Id, &d.CreatedAt)
	if err != nil {
		log.Println(err, "FAIL TO SAVE MESSAGE")
	}
	return
}


func (m *grupoRepo) fetchMessagesGrupo(ctx context.Context, query string, args ...interface{}) (res []r.MessageGrupo, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			log.Println(errRow)
		}
	}()
	res = make([]r.MessageGrupo, 0)
	for rows.Next() {
		t := r.MessageGrupo{}
		err = rows.Scan(
			&t.Id,
			&t.ChatId,
			&t.ProfileId,
			&t.Content,
			&t.Data,
			&t.CreatedAt,
			&t.ReplyTo,
			&t.TypeMessage,
			// &t.ReplyMessage.Id,
			// &t.ReplyMessage.GrupoId,
			// &t.ReplyMessage.ProfileId,
			// &t.ReplyMessage.Content,
			// &t.ReplyMessage.CreatedAt,
			// &t.ReplyMessage.TypeMessage,
			// &t.ReplyMessage.Data,
		)
		res = append(res, t)
	}
	return res, nil
}




// func (m *grupoRepo) fetchMessagesGrupo(ctx context.Context, query string, args ...interface{}) (res []r.MessageGrupo, err error) {
// 	rows, err := m.Conn.QueryContext(ctx, query, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer func() {
// 		errRow := rows.Close()
// 		if errRow != nil {
// 			log.Println(errRow)
// 		}
// 	}()
// 	res = make([]r.MessageGrupo, 0)
// 	for rows.Next() {
// 		t := r.MessageGrupo{}
// 		err = rows.Scan(
// 			&t.Id,
// 			&t.GrupoId,
// 			&t.ProfileId,
// 			&t.Content,
// 			&t.CreatedAt,
// 			&t.ReplyTo,
// 			&t.ReplyMessage.Id,
// 			&t.ReplyMessage.GrupoId,
// 			&t.ReplyMessage.ProfileId,
// 			&t.ReplyMessage.Content,
// 			&t.ReplyMessage.CreatedAt,
// 		)
// 		res = append(res, t)
// 	}
// 	return res, nil
// }
