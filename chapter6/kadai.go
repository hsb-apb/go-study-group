package chapter6

//go:generate mockgen -source=$GOFILE -destination=kadai_mock.go -package=$GOPACKAGE -self_package=github.com/apbgo/go-study-group/$GOPACKAGE

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// [課題内容]
// 以下の2つのInterface(IFUserItemService, IFUserItemRepository)を満たすstructを実装してください。
// 今回の課題ではトランザクション境界はProvide()内で構いません。
// gomockを利用してRepositoryのMockファイルを自動生成しています。
// テストではIFUserItemRepositoryをモックしたUserItemServiceを使ってみましょう。
// TransactionのBegin, Commit, Rollbackも本来Mockしたいですが、実装が多くなってしまうので
// 今回はする必要がありません。

// Reward 報酬モデル
type Reward struct {
	ItemID int64
	Count  int64
}

// IFUserItemService 報酬の付与の機能を表すインターフェイス
type IFUserItemService interface {
	// 対象のUserIDに引数で渡された報酬を付与します.
	Provide(ctx context.Context, userID int64, rewards ...Reward) error
}

// IFUserItemRepository i_user_itemテーブルへの操作を行うインターフェイス
type IFUserItemRepository interface {
	// FindByUserIdAndItemIDs 一致するモデルを複数返却する.
	FindByUserIdAndItemIDs(
		ctx context.Context,
		tx *sql.Tx,
		userID int64,
		itemIDs []int64,
	) (iUserItems []*IUserItem, err error)

	// Insert 対象のモデルから1件Insertを実行する
	Insert(
		ctx context.Context,
		tx *sql.Tx,
		iUserItem *IUserItem,
	) error

	// Update対象のモデルから1件Updateを実行する
	// Update対象レコードが0件の場合、okはfalseになる
	Update(
		ctx context.Context,
		tx *sql.Tx,
		iUserItem *IUserItem,
	) (ok bool, err error)
}

// UserItemService [実装対象]
type UserItemService struct {
	db                 *sql.DB
	userItemRepository IFUserItemRepository
}

func (u *UserItemService) Provide(ctx context.Context, userID int64, rewards ...Reward) error {

	// Transactionを開始
	tx, err := u.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return err
	}

	defer func() {
		if err := recover(); err != nil {
			// panic時のロールバック
			tx.Rollback()
		}
	}()

	items := make([]int64, len(rewards))
	for i, reward := range rewards {
		items[i] = reward.ItemID
	}
	userItems, err := u.userItemRepository.FindByUserIdAndItemIDs(ctx, tx, userID, items)
	if err != nil {
		tx.Rollback()
		return err
	}

	itemMap := make(map[int64]*IUserItem)
	for _, userItem := range userItems {
		itemMap[userItem.ItemID] = userItem
	}

	for _, reward := range rewards {
		if userItem, ok := itemMap[reward.ItemID]; ok {
			// 更新する
			userItem.Count += reward.Count
			_, err := u.userItemRepository.Update(ctx, tx, userItem)
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// 持っていないので登録する
			userItem := IUserItem{
				UserID: userID,
				ItemID: reward.ItemID,
				Count:  reward.Count,
			}
			err := u.userItemRepository.Insert(ctx, tx, &userItem)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

	}
	// すべてうまく行ったらコミット
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// NewUserItemService コンストラクタ [実装対象]
func NewUserItemService(
	db *sql.DB,
	userItemRepository *UserItemRepository,
) *UserItemService {
	return &UserItemService{db, userItemRepository}
}

// UserItemRepository [実装対象]
type UserItemRepository struct {
}

func (userItemRepository *UserItemRepository) FindByUserIdAndItemIDs(
	ctx context.Context,
	tx *sql.Tx,
	userID int64,
	itemIDs []int64,
) (iUserItems []*IUserItem, err error) {

	query, args, err := sqlx.In("SELECT * FROM i_user_item WHERE item_id IN (?);", itemIDs)
	// SQLを実行
	rows, err := tx.QueryContext(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 1レコードずつ、処理する
	for rows.Next() {
		// モデルを作成して、カラムへのポインタをScan()に渡す
		record := IUserItem{}
		if err = rows.Scan(
			&record.UserID,
			&record.ItemID,
			&record.Count,
			&record.CreatedAt,
			&record.UpdatedAt,
			&record.DeletedAt,
		); err != nil {
			return nil, err
		}
		iUserItems = append(iUserItems, &record)
	}

	// 処理中にエラーが発生する場合もあるのでここの処理を忘れずに
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return iUserItems, err
}

func (userItemRepository *UserItemRepository) Insert(
	ctx context.Context,
	tx *sql.Tx,
	iUserItem *IUserItem,
) error {

	_, err := tx.ExecContext(ctx, "INSERT INTO i_user_item VALUES (?, ?, ?, now(), now(), NULL)", iUserItem.UserID, iUserItem.ItemID, iUserItem.Count)
	if err != nil {
		return err
	}

	return nil
}

func (userItemRepository *UserItemRepository) Update(
	ctx context.Context,
	tx *sql.Tx,
	iUserItem *IUserItem,
) (ok bool, err error) {

	result, err := tx.ExecContext(ctx, "UPDATE i_user_item SET count = ?, updated_at = now() WHERE user_id = ? and item_id = ?", iUserItem.Count, iUserItem.UserID, iUserItem.ItemID)
	if err != nil {
		return false, err
	}
	updateCount, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if updateCount != 1 {
		return false, fmt.Errorf("update count is false")
	}

	return true, nil
}

// NewUserItemRepository コンストラクタ [実装対象]
func NewUserItemRepository() *UserItemRepository {
	return &UserItemRepository{}
}
