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

func (p *salaRepo) GetChatUnreadMessage(ctx context.Context, chatId int64, lastUpdated string) (res []r.Message, err error) {
	query := `select gm.id,gm.chat_id,gm.profile_id,gm.content,gm.data,gm.created_at,gm.reply_to,gm.type_message
	 from grupo_message as gm where chat_id = $1 and gm.created_at >= $2`
	res, err = p.fetchMessagesGrupo(ctx, query, chatId, lastUpdated)
	return

	// select gm.id,gm.chat_id,gm.profile_id,gm.content,gm.data,gm.created_at,gm.reply_to,gm.type_message
	//  from grupo_message as gm where chat_id = 1 and gm.created_at >= $2
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
		)
		res = append(res, t)
	}
	return res, nil
}
