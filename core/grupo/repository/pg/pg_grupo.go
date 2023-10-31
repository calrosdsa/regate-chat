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
func (p *grupoRepo)DeleteMessage(ctx context.Context,id int)(err error){
	query := `update grupo_message set is_deleted = true where id = $1`
	_,err = p.Conn.ExecContext(ctx,query,id)
	return 
}
func (p *grupoRepo) GetUnreadMessages(ctx context.Context, profileId int, page int16,
	size int8) (res []r.Message, err error) {
	query := `select  gm.id,gm.chat_id,gm.profile_id,gm.content,gm.data,gm.created_at,gm.reply_to,gm.type_message,
	gm.is_deleted from user_grupo as ug inner join grupo_message as gm on gm.grupo_id = ug.grupo_id 
	and ug.last_update_messages <= gm.created_at where ug.profile_id = $1 
	limit $2 offset $3`
	res, err = p.fetchMessagesGrupo(ctx, query, profileId, size, page*int16(size))
	return
}
func (p *grupoRepo) GetChatUnreadMessage(ctx context.Context, chatId int, lastUpdated string) (res []r.Message, err error) {
	query := `select gm.id,gm.chat_id,gm.profile_id,gm.content,gm.data,gm.created_at,gm.reply_to,gm.type_message,gm.is_deleted
	 from grupo_message as gm where chat_id = $1 and gm.created_at > $2 limit 100`
	res, err = p.fetchMessagesGrupo(ctx, query, chatId, lastUpdated)
	return

	// select gm.id,gm.chat_id,gm.profile_id,gm.content,gm.data,gm.created_at,gm.reply_to,gm.type_message
	//  from grupo_message as gm where chat_id = 1 and gm.created_at >= $2
}

func (p *grupoRepo) UpdateUserGrupoLastTimeUpdateMessage(ctx context.Context, profileId int) (err error) {
	query := `update user_grupo set last_update_messages = current_timestamp where profile_id = $1`
	_, err = p.Conn.ExecContext(ctx, query, profileId)
	return
}

func (p *grupoRepo) SaveGrupoMessage(ctx context.Context, d *r.Message) (err error) {
	log.Println(d.CreatedAt, "CreatedAt Message")
	query := `insert into grupo_message (chat_id,profile_id,content,created_at,reply_to,type_message,data,grupo_id) 
	values($1,$2,$3,current_timestamp,$4,$5,$6,$7) returning id,created_at`
	err = p.Conn.QueryRowContext(ctx, query, d.ChatId, d.ProfileId, d.Content, d.ReplyTo,
		d.TypeMessage, d.Data, d.ParentId).Scan(&d.Id, &d.CreatedAt)
	if err != nil {
		log.Println(err, "FAIL TO SAVE MESSAGE")
	}
	return
}
func (p *grupoRepo)GetUsers(ctx context.Context,d r.RequestUsersGroupOrRoom)(res []r.UsersGroupOrRoom,err error){
	var (
		query string
		activesCount int
		inactivesCount int
	)
	query = `select count(*) FILTER(WHERE is_out = false) as actives,
	count(*) FILTER(WHERE is_out = true) as inactives from user_grupo where grupo_id = $1`
	err = p.Conn.QueryRowContext(ctx,query,d.ParentId).Scan(&activesCount,&inactivesCount)
	if err != nil {
		return
	}
	if activesCount == d.ActiveUsersCount && inactivesCount == d.InactiveUsersCount {
		return
	}
	sumUsers := activesCount + inactivesCount
	sumUsersLocal := d.ActiveUsersCount + d.InactiveUsersCount
	if sumUsers == sumUsersLocal {
		return
	}
	diff := sumUsers - sumUsersLocal
	query = `select u.id,u.is_admin,u.is_out,p.profile_id,p.nombre,p.apellido,p.profile_photo
	from user_grupo as u inner join profiles as p on p.profile_id = u.profile_id
	where grupo_id = $1 order by u.updated_at desc limit $2 `
	res,err = p.fetchUsers(ctx,query,d.ParentId,diff)
	return
}

func (m *grupoRepo) fetchUsers(ctx context.Context, query string, args ...interface{}) (res []r.UsersGroupOrRoom, err error) {
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
	res = make([]r.UsersGroupOrRoom, 0)
	for rows.Next() {
		t := r.UsersGroupOrRoom{}
		err = rows.Scan(
			&t.Id,
			&t.IsAdmin,
			&t.IsOut,
			&t.ProfileId,
			&t.ProfileName,
			&t.ProfileApellido,
			&t.ProfilePhoto,
		)
		res = append(res, t)
	}
	return res, nil
}

func (m *grupoRepo) fetchMessagesGrupo(ctx context.Context, query string, args ...interface{}) (res []r.Message, err error) {
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
	res = make([]r.Message, 0)
	for rows.Next() {
		t := r.Message{}
		err = rows.Scan(
			&t.Id,
			&t.ChatId,
			&t.ProfileId,
			&t.Content,
			&t.Data,
			&t.CreatedAt,
			&t.ReplyTo,
			&t.TypeMessage,
			&t.IsDeleted,
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
