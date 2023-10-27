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
func (p *conversationAdminRepo) GetMessages(ctx context.Context, profileId int, page int16,
	size int8) (res []r.Message, err error) {
	query := `select  m.id,m.chat_id,m.profile_id,m.content,m.data,m.created_at,m.reply_to,
	m.type_message,m.is_user,m.is_read
	from conversation_message as m where m.chat_id = $1 
	order by created_at desc limit $2 offset $3`
	res, err = p.fetchMessages(ctx, query, profileId, size, page*int16(size))
	return
}

func (p *conversationAdminRepo) GetConversationsEstablecimiento(ctx context.Context,uuid string) (res []r.ChatEstablecimiento,err error) {
	query := `select p.nombre,p.apellido,p.profile_photo,c.id,p.profile_id,c.parent_id,
	cm.content,cm.created_at,
	(select count(*) from conversation_message where c.id = chat_id and is_read = false)
	from chat as c left join lateral 
	(select m.content,m.created_at from conversation_message as m 
	where c.id = m.chat_id
	order by created_at desc limit 1 ) cm on true
	inner join profiles as p on c.second_parent_id = p.profile_id
	where c.parent_id = (select establecimiento_id from establecimientos where uuid = $1)`
	res,err = p.fetchConversationsAdmin(ctx,query,uuid)
	// query := `insert into conversation_message (id,conversation_id,sender_id,content,created_at,reply_to) 
	// values($1,$2,$3,$4,$5,$6) returning id,created_at`
	return
}

func (m *conversationAdminRepo) fetchMessages(ctx context.Context, query string, args ...interface{}) (res []r.Message, err error) {
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
			&t.IsUser,
			&t.IsRead,
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

func (p *conversationAdminRepo) fetchConversationsAdmin(ctx context.Context, query string,
	 args ...interface{}) (res []r.ChatEstablecimiento, err error) {
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
	res = make([]r.ChatEstablecimiento, 0)
	for rows.Next() {
		t := r.ChatEstablecimiento{}
		err = rows.Scan(
			&t.Chat.Name,
			&t.Chat.Apellido,
			&t.Chat.Photo,
			&t.Chat.ConversationId,
			&t.Chat.ProfileId,
			&t.Chat.ParentId,

			&t.Message.Content,
			&t.Message.CreatedAt,

			&t.CounUnreadMessages,

		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		res = append(res, t)
	}
	return res, nil
}