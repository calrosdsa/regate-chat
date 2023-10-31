package pg

import (
	"context"
	"database/sql"
	"log"

	// "log"
	r "message/domain/repository"
)

type salaRepo struct {
	Conn *sql.DB
}

func NewRepository(conn *sql.DB) r.SalaRepository {
	return &salaRepo{
		Conn: conn,
	}
}


func (p *salaRepo) SaveMessage(ctx context.Context, d *r.Message) (err error) {
	log.Println(d.CreatedAt, "CreatedAt Message")
	query := `insert into sala_message (chat_id,profile_id,content,created_at,reply_to,type_message,data,sala_id) 
	values($1,$2,$3,current_timestamp,$4,$5,$6,$7) returning id,created_at`
	err = p.Conn.QueryRowContext(ctx, query, d.ChatId, d.ProfileId, d.Content, d.ReplyTo,
		d.TypeMessage, d.Data, d.ParentId).Scan(&d.Id, &d.CreatedAt)
	if err != nil {
		log.Println(err, "FAIL TO SAVE MESSAGE")
	}
	return
}

func (p *salaRepo)DeleteMessage(ctx context.Context,id int)(err error){
	query := `update sala_message set is_deleted = true where id = $1`
	_,err = p.Conn.ExecContext(ctx,query,id)
	return
}

func (p *salaRepo)GetUsers(ctx context.Context,d r.RequestUsersGroupOrRoom)(res []r.UsersGroupOrRoom,err error){
	var (
		query string
		activesCount int
		inactivesCount int
	)
	query = `select count(*) FILTER(WHERE is_out = false) as actives,
	count(*) FILTER(WHERE is_out = true) as inactives from users_sala where sala_id = $1`
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
	from users_sala as u inner join profiles as p on p.profile_id = u.profile_id
	where sala_id = $1 order by u.updated_at desc limit $2`
	res,err = p.fetchUsers(ctx,query,d.ParentId,diff)
	return
}
func (p *salaRepo) GetChatUnreadMessages(ctx context.Context, chatId int, lastUpdated string) (res []r.Message, err error) {
	query := `select m.id,m.chat_id,m.profile_id,m.content,m.data,m.created_at,
	m.reply_to,m.type_message,is_deleted
	from sala_message as m where chat_id = $1 and m.created_at > $2`
	res, err = p.fetchMessagesGrupo(ctx, query, chatId, lastUpdated)
	return
}
func (m *salaRepo) fetchMessagesGrupo(ctx context.Context, query string, args ...interface{}) (res []r.Message, err error) {
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


func (m *salaRepo) fetchUsers(ctx context.Context, query string, args ...interface{}) (res []r.UsersGroupOrRoom, err error) {
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

