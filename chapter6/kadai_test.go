package chapter6

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// このテストはgomockを利用するサンプルです。
func TestSample(t *testing.T) {
	t.Run("サンプル1", func(t *testing.T) {
		// サブテストごとにMockのControllerを作成してください。
		ctrl := gomock.NewController(t)

		// 自動生成されたMockをNewする
		mock := NewMockIFUserItemRepository(ctrl)

		// ここからは意味がないテスト
		ctx := context.Background()
		userItem := IUserItem{
			UserID: 1,
			ItemID: 1,
			Count:  100,
		}

		mock.EXPECT().
			// ここに渡された変数は、値が一致しない場合はテストが成功しない（ポインタの場合はポインタの一致）
			// gomock.Any()は全ての値が許容される
			// *sql.TxやiUserItemはService内で生成されるため、ポインタの一致をチェックすることは難しい
			Update(ctx, gomock.Any(), gomock.Any()).
			// そのためDoAndReturnの関数内で値をチェックしてあげるとよい
			DoAndReturn(
				func(ctx1 context.Context, tx1 *sql.Tx, userItem1 *IUserItem) (ok bool, err error) {
					assert.Equal(t, userItem, *userItem1)
					// この関数の戻り値がMockを実行した時の戻り値になる
					return true, nil
				},
			)

		ok, err := mock.Update(ctx, nil, &userItem)
		assert.True(t, ok)
		assert.NoError(t, err)
	})
}

func TestUserItemService_Provide(t *testing.T) {

	t.Run("正常系 1個づつ update / insert", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockUserItemRepository := NewMockIFUserItemRepository(ctrl)
		mockUserItemRepository.EXPECT().
			FindByUserIdAndItemIDs(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx1 context.Context, tx1 *sql.Tx, userID int64, itemIDs []int64) (iUserItems []*IUserItem, err error) {
					// itemIDが1のものだけ持っている想定
					record := IUserItem{}
					record.UserID = userID
					record.ItemID = itemIDs[0]
					record.Count = 1
					iUserItems = append(iUserItems, &record)
					return iUserItems, nil
				},
			)

		mockUserItemRepository.EXPECT().
			Update(ctx, gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx1 context.Context, tx1 *sql.Tx, userItem1 *IUserItem) (ok bool, err error) {
					// カウントが2になってアップデートされるはず
					assert.Equal(t, int64(1), userItem1.UserID)
					assert.Equal(t, int64(1), userItem1.ItemID)
					assert.Equal(t, int64(2), userItem1.Count)
					return true, nil
				},
			)

		mockUserItemRepository.EXPECT().
			Insert(ctx, gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx1 context.Context, tx1 *sql.Tx, userItem1 *IUserItem) (err error) {
					// itemIdが2のアイテムが1つinsertされるはず
					assert.Equal(t, int64(1), userItem1.UserID)
					assert.Equal(t, int64(2), userItem1.ItemID)
					assert.Equal(t, int64(1), userItem1.Count)
					return nil
				},
			)

		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:5446)/chapter6?parseTime=true")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		userItemService := UserItemService{
			db:                 db,
			userItemRepository: mockUserItemRepository,
		}

		// ------ テストデータ -------- //
		reward1 := Reward{
			ItemID: 1,
			Count:  1,
		}
		reward2 := Reward{
			ItemID: 2,
			Count:  1,
		}

		assert.NoError(t, userItemService.Provide(ctx, 1, reward1, reward2))

	})

	t.Run("正常系 複数個 update / insert", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockUserItemRepository := NewMockIFUserItemRepository(ctrl)
		mockUserItemRepository.EXPECT().
			FindByUserIdAndItemIDs(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx1 context.Context, tx1 *sql.Tx, userID int64, itemIDs []int64) (iUserItems []*IUserItem, err error) {
					// itemIDが1のものと3のものを持っている
					record := IUserItem{}
					record.UserID = userID
					record.ItemID = itemIDs[0]
					record.Count = 1
					iUserItems = append(iUserItems, &record)

					record2 := IUserItem{}
					record2.UserID = userID
					record2.ItemID = itemIDs[2]
					record2.Count = 2
					iUserItems = append(iUserItems, &record2)
					return iUserItems, nil
				},
			)

		mockUserItemRepository.EXPECT().
			Update(ctx, gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx1 context.Context, tx1 *sql.Tx, userItem1 *IUserItem) (ok bool, err error) {

					switch userItem1.ItemID {
					case 1:
						// カウントが2になってアップデートされるはず
						assert.Equal(t, int64(1), userItem1.UserID)
						assert.Equal(t, int64(1), userItem1.ItemID)
						assert.Equal(t, int64(2), userItem1.Count)
					case 3:
						// カウントが2になってアップデートされるはず
						assert.Equal(t, int64(1), userItem1.UserID)
						assert.Equal(t, int64(3), userItem1.ItemID)
						assert.Equal(t, int64(12), userItem1.Count)
					}
					return true, nil
				},
			).AnyTimes()

		mockUserItemRepository.EXPECT().
			Insert(ctx, gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx1 context.Context, tx1 *sql.Tx, userItem1 *IUserItem) (err error) {

					switch userItem1.ItemID {
					case 2:
						// itemIdが2のアイテムが1つinsertされるはず
						assert.Equal(t, int64(1), userItem1.UserID)
						assert.Equal(t, int64(2), userItem1.ItemID)
						assert.Equal(t, int64(1), userItem1.Count)
					case 4:
						// itemIdが4のアイテムが15こinsertされるはず
						assert.Equal(t, int64(1), userItem1.UserID)
						assert.Equal(t, int64(4), userItem1.ItemID)
						assert.Equal(t, int64(15), userItem1.Count)

					}
					return nil
				},
			).AnyTimes()

		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:5446)/chapter6?parseTime=true")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		userItemService := UserItemService{
			db:                 db,
			userItemRepository: mockUserItemRepository,
		}

		// ------ テストデータ -------- //
		reward1 := Reward{
			ItemID: 1,
			Count:  1,
		}
		reward2 := Reward{
			ItemID: 2,
			Count:  1,
		}
		reward3 := Reward{
			ItemID: 3,
			Count:  10,
		}
		reward4 := Reward{
			ItemID: 4,
			Count:  15,
		}

		assert.NoError(t, userItemService.Provide(ctx, 1, reward1, reward2, reward3, reward4))

	})

	t.Run("異常系 select", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ctx := context.Background()

		mockUserItemRepository := NewMockIFUserItemRepository(ctrl)
		mockUserItemRepository.EXPECT().
			FindByUserIdAndItemIDs(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx1 context.Context, tx1 *sql.Tx, userID int64, itemIDs []int64) (iUserItems []*IUserItem, err error) {
					return iUserItems, fmt.Errorf("select error")
				},
			)

		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:5446)/chapter6?parseTime=true")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		userItemService := UserItemService{
			db:                 db,
			userItemRepository: mockUserItemRepository,
		}
		reward1 := Reward{
			ItemID: 1,
			Count:  1,
		}
		reward2 := Reward{
			ItemID: 2,
			Count:  1,
		}
		assert.Error(t, userItemService.Provide(ctx, 1, reward1, reward2))
	})

	t.Run("異常系 update", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ctx := context.Background()

		mockUserItemRepository := NewMockIFUserItemRepository(ctrl)
		mockUserItemRepository.EXPECT().
			FindByUserIdAndItemIDs(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx1 context.Context, tx1 *sql.Tx, userID int64, itemIDs []int64) (iUserItems []*IUserItem, err error) {
					record := IUserItem{}
					record.UserID = userID
					record.ItemID = itemIDs[0]
					record.Count = 1
					iUserItems = append(iUserItems, &record)
					return iUserItems, nil
				},
			)

		mockUserItemRepository.EXPECT().
			Update(ctx, gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx1 context.Context, tx1 *sql.Tx, userItem1 *IUserItem) (ok bool, err error) {
					// updateでなにがしかのエラーが発生
					return false, fmt.Errorf("update error")
				},
			)

		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:5446)/chapter6?parseTime=true")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		userItemService := UserItemService{
			db:                 db,
			userItemRepository: mockUserItemRepository,
		}
		reward1 := Reward{
			ItemID: 1,
			Count:  1,
		}
		reward2 := Reward{
			ItemID: 2,
			Count:  1,
		}
		assert.Error(t, userItemService.Provide(ctx, 1, reward1, reward2))
	})

	t.Run("異常系 insert", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		mockUserItemRepository := NewMockIFUserItemRepository(ctrl)
		mockUserItemRepository.EXPECT().
			FindByUserIdAndItemIDs(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx1 context.Context, tx1 *sql.Tx, userID int64, itemIDs []int64) (iUserItems []*IUserItem, err error) {
					return iUserItems, nil
				},
			)

		mockUserItemRepository.EXPECT().
			Insert(ctx, gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx1 context.Context, tx1 *sql.Tx, userItem1 *IUserItem) (err error) {
					// insertでなにがしかのエラーが発生
					return fmt.Errorf("insert error")
				},
			)

		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:5446)/chapter6?parseTime=true")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		userItemService := UserItemService{
			db:                 db,
			userItemRepository: mockUserItemRepository,
		}
		reward1 := Reward{
			ItemID: 1,
			Count:  1,
		}
		reward2 := Reward{
			ItemID: 2,
			Count:  1,
		}
		assert.Error(t, userItemService.Provide(ctx, 1, reward1, reward2))
	})
}
