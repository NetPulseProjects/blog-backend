package postgresql

import (
	"app/pkg/domain/entity"
	"app/pkg/lib/ers"
	"database/sql"
	"github.com/google/uuid"
	"log"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (repo UserRepo) Subscribe(subscription entity.UserSubscription) error {
	const op = "postgresql.UserRepo.Subscribe"

	subscribeStmt, err := repo.db.Prepare(
		`INSERT INTO user_subscription(owner_id, subscribed_user_id, created_at) VALUES ($1, $2, $3)`,
	)
	if err != nil {
		return ers.ThrowMessage(op, err)
	}

	_, err = subscribeStmt.Exec(subscription.OwnerId, subscription.SubscribedUserId, subscription.CreatedAt)
	if err != nil {
		return ers.ThrowMessage(op, err)
	}

	return nil
}

func (repo UserRepo) Unsubscribe(ownerId, unsubscribedId uuid.UUID) error {
	const op = "postgresql.UserRepo.Unsubscribe"

	unsubscribeStmt, err := repo.db.Prepare(
		`DELETE FROM user_subscription WHERE owner_id = $1 AND subscribed_user_id = $2`,
	)
	if err != nil {
		return ers.ThrowMessage(op, err)
	}

	_, err = unsubscribeStmt.Exec(ownerId, unsubscribedId)
	if err != nil {
		return ers.ThrowMessage(op, err)
	}

	return nil
}

func (repo UserRepo) GetSubSiteBarItems() (*[]entity.UserSubSiteBarItem, error) {
	const op = "postgresql.UserRepo.GetSubSiteBarItems"

	rows, err := repo.db.Query(`SELECT id, avatar_url, name FROM "user" WHERE role = 'sub_site'`)
	if err != nil {
		return nil, ers.ThrowMessage(op, err)
	}

	defer rows.Close()

	var items []entity.UserSubSiteBarItem

	for rows.Next() {
		var row entity.UserSubSiteBarItem
		err = rows.Scan(&row.Id, &row.AvatarUrl, &row.Name)

		if err != nil {
			log.Println(err)
			return nil, ers.ThrowMessage(op, err)
		}

		items = append(items, row)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
	}

	return &items, nil
}

func (repo UserRepo) FindById(id uuid.UUID) (*entity.User, error) {
	const op = "postgresql.UserRepo.FindById"

	fUserSettings := entity.UserSettings{}
	fUser := entity.User{}

	row := repo.db.QueryRow(`
	   SELECT 
			u.*,
			(SELECT COUNT(*) FROM user_subscription WHERE owner_id = u.id) AS subscriptions_count,
			(SELECT COUNT(*) FROM user_subscription WHERE subscribed_user_id = u.id) AS subscribers_count,
			us.news_line_default, us.news_line_sort
		FROM 
			"user" u
		INNER JOIN 
			"user_settings" us ON u.id = us.user_id
		WHERE 
			u.id = $1
    `, id)
	err := row.Scan(
		&fUser.Id,
		&fUser.EncryptedPassword,
		&fUser.Salt,
		&fUser.CreatedAt,
		&fUser.UpdatedAt,
		&fUser.Role,
		&fUser.Email,
		&fUser.Name,
		&fUser.Description,
		&fUser.AvatarUrl,
		&fUser.CoverUrl,
		&fUser.SubscriptionsCount,
		&fUser.SubscribersCount,
		&fUserSettings.NewsLineDefault,
		&fUserSettings.NewsLineSort,
	)
	if err != nil {
		return &fUser, ers.ThrowMessage(op, err)
	}

	fUserSettings.UserId = fUser.Id
	fUser.Settings = fUserSettings

	return &fUser, nil
}

func (repo UserRepo) UpdateSettings(uSettings *entity.UserSettings) error {
	const op = "postgresql.UserRepo.Subscribe"

	updateSettingsStmt, err := repo.db.Prepare(
		`UPDATE user_settings SET news_line_default = $2, news_line_sort = $3 WHERE user_id = $1`,
	)
	if err != nil {
		return ers.ThrowMessage(op, err)
	}

	_, err = updateSettingsStmt.Exec(
		uSettings.UserId,
		uSettings.NewsLineDefault,
		uSettings.NewsLineSort,
	)
	if err != nil {
		return ers.ThrowMessage(op, err)
	}

	return nil
}

func (repo UserRepo) GetByAuthId(id uuid.UUID) (*entity.User, error) {
	const op = "postgresql.UserRepo.GetByAuthId"

	fUser := entity.User{
		Settings: entity.UserSettings{},
	}

	row := repo.db.QueryRow(`
	   SELECT 
			u.*,
			(SELECT COUNT(*) FROM user_subscription WHERE owner_id = u.id) AS subscriptions_count,
			(SELECT COUNT(*) FROM user_subscription WHERE subscribed_user_id = u.id) AS subscribers_count,
			us.news_line_default, us.news_line_sort
		FROM 
			"user" u
		INNER JOIN 
			"user_settings" us ON u.id = us.user_id
		INNER JOIN 
			"user_auth" ua ON u.id = ua.user_id
		WHERE 
			ua.id = $1
    `, id)
	err := row.Scan(
		&fUser.Id,
		&fUser.EncryptedPassword,
		&fUser.Salt,
		&fUser.CreatedAt,
		&fUser.UpdatedAt,
		&fUser.Role,
		&fUser.Email,
		&fUser.Name,
		&fUser.Description,
		&fUser.AvatarUrl,
		&fUser.CoverUrl,
		&fUser.SubscriptionsCount,
		&fUser.SubscribersCount,
		&fUser.Settings.NewsLineDefault,
		&fUser.Settings.NewsLineSort,
	)
	if err != nil {
		return &fUser, ers.ThrowMessage(op, err)
	}

	fUser.Settings.UserId = fUser.Id

	return &fUser, nil
}

func (repo UserRepo) FindByEmail(email string) (*entity.User, error) {
	const op = "postgresql.UserRepo.FindByEmail"

	fUser := entity.User{
		Settings: entity.UserSettings{},
	}

	row := repo.db.QueryRow(`
	   SELECT 
			u.*,
			(SELECT COUNT(*) FROM user_subscription WHERE owner_id = u.id) AS subscriptions_count,
			(SELECT COUNT(*) FROM user_subscription WHERE subscribed_user_id = u.id) AS subscribers_count,
			us.news_line_default, us.news_line_sort
		FROM 
			"user" u
		INNER JOIN 
			"user_settings" us ON u.id = us.user_id
		WHERE 
			u.email = $1
    `, email)
	err := row.Scan(
		&fUser.Id,
		&fUser.EncryptedPassword,
		&fUser.Salt,
		&fUser.CreatedAt,
		&fUser.UpdatedAt,
		&fUser.Role,
		&fUser.Email,
		&fUser.Name,
		&fUser.Description,
		&fUser.AvatarUrl,
		&fUser.CoverUrl,
		&fUser.SubscriptionsCount,
		&fUser.SubscribersCount,
		&fUser.Settings.NewsLineDefault,
		&fUser.Settings.NewsLineSort,
	)
	if err != nil {
		return &fUser, ers.ThrowMessage(op, err)
	}

	fUser.Settings.UserId = fUser.Id

	return &fUser, nil
}

func (repo UserRepo) CreatePersonal(newUser *entity.User) error {
	const op = "postgresql.UserRepo.CreatePersonal"

	tx, err := repo.db.Begin()
	if err != nil {
		return ers.ThrowMessage(op, err)
	}

	newUserStmt, err := tx.Prepare(
		`INSERT INTO "user"(id, encrypted_password, salt, created_at, updated_at, role, email, name, description, avatar_url, cover_url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
	)
	if err != nil {
		return ers.ThrowMessage(op, err)
	}

	_, err = newUserStmt.Exec(
		newUser.Id,
		newUser.EncryptedPassword,
		newUser.Salt,
		newUser.CreatedAt,
		newUser.UpdatedAt,
		newUser.Role,
		newUser.Email,
		newUser.Name,
		newUser.Description,
		newUser.AvatarUrl,
		newUser.CoverUrl,
	)
	if err != nil {
		return ers.ThrowMessage(op, err)
	}

	userSettingsStmt, err := tx.Prepare(
		`INSERT INTO "user_settings"(user_id, news_line_default, news_line_sort) VALUES ($1, $2, $3)`,
	)
	if err != nil {
		tx.Rollback()
		return ers.ThrowMessage(op, err)
	}

	_, err = userSettingsStmt.Exec(
		newUser.Id,
		newUser.Settings.NewsLineDefault,
		newUser.Settings.NewsLineSort,
	)
	if err != nil {
		tx.Rollback()
		return ers.ThrowMessage(op, err)
	}

	err = tx.Commit()
	if err != nil {
		return ers.ThrowMessage(op, err)
	}

	return nil
}
