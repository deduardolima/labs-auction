package auction_test

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/infra/database/auction"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	"github.com/stretchr/testify/assert"
)

func TestCreateAuction(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful auction creation", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		repo := auction.NewAuctionRepository(mt.DB)
		ctx := context.Background()

		auctionEntity := &auction_entity.Auction{
			Id:          "auction-id-1",
			ProductName: "Product 1",
			Category:    "Category 1",
			Description: "Description 1",
			Condition:   auction_entity.New,
			Status:      auction_entity.Active,
			Timestamp:   time.Now(),
		}

		os.Setenv("AUCTION_DURATION", "60")

		err := repo.CreateAuction(ctx, auctionEntity)
		assert.Nil(t, err)
	})

	mt.Run("invalid auction duration", func(mt *mtest.T) {
		repo := auction.NewAuctionRepository(mt.DB)
		ctx := context.Background()

		auctionEntity := &auction_entity.Auction{
			Id:          "auction-id-1",
			ProductName: "Product 1",
			Category:    "Category 1",
			Description: "Description 1",
			Condition:   auction_entity.New,
			Status:      auction_entity.Active,
			Timestamp:   time.Now(),
		}

		os.Setenv("AUCTION_DURATION", "invalid-duration")

		err := repo.CreateAuction(ctx, auctionEntity)
		assert.NotNil(t, err)
		assert.Equal(t, internal_error.NewInternalServerError("Invalid auction duration").Message, err.Message)
	})

	mt.Run("error inserting auction", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Code:    11000,
			Message: "duplicate key error",
		}))

		repo := auction.NewAuctionRepository(mt.DB)
		ctx := context.Background()

		auctionEntity := &auction_entity.Auction{
			Id:          "auction-id-1",
			ProductName: "Product 1",
			Category:    "Category 1",
			Description: "Description 1",
			Condition:   auction_entity.New,
			Status:      auction_entity.Active,
			Timestamp:   time.Now(),
		}

		os.Setenv("AUCTION_DURATION", "60")

		err := repo.CreateAuction(ctx, auctionEntity)
		assert.NotNil(t, err)
		assert.Equal(t, internal_error.NewInternalServerError("Error trying to insert auction").Message, err.Message)
	})
}
