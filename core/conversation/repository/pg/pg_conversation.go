package pg

import (
	"context"
	"database/sql"
	"log"
	r "message/domain/repository"

)

type conversationRepo struct {
	Conn *sql.DB
}

func NewRepository(conn *sql.DB) r.ConversationRepository {
	return &conversationRepo{
		Conn: conn,
	}
}

func (p *conversationRepo) GetOrCreateConversation(ctx context.Context, id int, profileId int) (conversationId int, err error) {
	var query string
	query = `select id from chat where second_parent_id = $1 and parent_id = $2`
	err = p.Conn.QueryRowContext(ctx, query, profileId, id).Scan(&conversationId)
	if err != nil {
		log.Println("Creando conversation")
		query = `insert into chat(second_parent_id,parent_id,type_chat) values($1,$2,$3) returning id`
		err = p.Conn.QueryRowContext(ctx, query, profileId, id, r.TypeChatInboxEstablecimiento).Scan(&conversationId)
		if err != nil {
			return
		}
	}
	return
}

func (p *conversationRepo) SaveMessage(ctx context.Context,d *r.Message) (err error) {
	// query := `insert into conversation_message (id,conversation_id,sender_id,content,created_at,reply_to)
	// values($1,$2,$3,$4,$5,$6) returning id,created_at`
	// err = p.Conn.QueryRowContext(ctx, query,d.Id ,d.ConversationId, d.SenderId, d.Content,d.CreatedAt, d.ReplyTo)
	// .Scan(&d.Id, &d.CreatedAt)
	log.Println(d.CreatedAt, "CreatedAt Message")
	query := `insert into conversation_message (chat_id,profile_id,content,created_at,reply_to,
	type_message,data,establecimiento_id,is_user) 
	values($1,$2,$3,current_timestamp,$4,$5,$6,$7,$8) returning id,created_at`
	err =p.Conn.QueryRowContext(ctx, query, d.ChatId, d.ProfileId, d.Content, d.ReplyTo,
		d.TypeMessage, d.Data, d.ParentId,d.IsUser).Scan(&d.Id, &d.CreatedAt)
	if err != nil {
		log.Println(err, "FAIL TO SAVE MESSAGE")
	}
	return
}


func (p *conversationRepo) GetConversations(ctx context.Context, id int) (res []r.Conversation, err error) {
	query := `select c.conversation_id,e.establecimiento_id,e.name,e.photo from conversations as c
	inner join establecimientos as e on e.establecimiento_id = c.establecimiento_id where c.profile_id = $1`
	res, err = p.fetchConversations(ctx, query, id)
	return
}
func(p *conversationRepo)UpdateMessageToReaded(ctx context.Context,id int)(err error){
	query := `update conversation_message  set is_read = true where id = $1`
	_,err = p.Conn.ExecContext(ctx,query,id)
	return 
}
func (p *conversationRepo)DeleteMessage(ctx context.Context,id int )(err error){
	query := `update conversation_message set is_deleted = true where id = $1`
	_,err = p.Conn.ExecContext(ctx,query,id)
	return
}

func (p *conversationRepo) GetChatUnreadMessages(ctx context.Context, chatId int, lastUpdated string) (res []r.Message, err error) {
	query := `select m.id,m.chat_id,m.profile_id,m.content,m.data,m.created_at,m.reply_to,
	m.type_message,m.is_user,m.is_deleted
	from conversation_message as m where m.chat_id = $1 and m.created_at > $2 limit 100`
	res, err = p.fetchMessagesGrupo(ctx, query, chatId, lastUpdated)
	return
}



func (m *conversationRepo) fetchMessagesGrupo(ctx context.Context, query string, args ...interface{}) (res []r.Message, err error) {
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
			&t.IsDeleted,
		)
		res = append(res, t)
	}
	return res, nil
}

func (p *conversationRepo) fetchConversations(ctx context.Context, query string,
args ...interface{}) (res []r.Conversation, err error) {
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
	res = make([]r.Conversation, 0)
	for rows.Next() {
		t := r.Conversation{}
		err = rows.Scan(
			&t.Id,
			&t.EstablecimientoId,
			&t.EstablecimientoName,
			&t.EstablecimientoPhoto,
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
