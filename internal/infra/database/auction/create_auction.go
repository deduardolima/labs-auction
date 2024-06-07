package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
	EndTime     int64                           `bson:"end_time"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {

	durationStr := os.Getenv("AUCTION_DURATION")
	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		logger.Error("Invalid auction duration", err)
		return internal_error.NewInternalServerError("Invalid auction duration")
	}
	endTime := auctionEntity.Timestamp.Add(time.Duration(duration) * time.Second).Unix()

	auctionEntity.Id = uuid.New().String()

	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
		EndTime:     endTime,
	}
	_, err = ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	go ar.checkAndCloseAuction(ctx, auctionEntityMongo)

	return nil
}

func (ar *AuctionRepository) checkAndCloseAuction(ctx context.Context, auctionEntityMongo *AuctionEntityMongo) {
	time.Sleep(time.Duration(auctionEntityMongo.EndTime-time.Now().Unix()) * time.Second)

	filter := bson.M{"_id": auctionEntityMongo.Id}
	update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}

	_, err := ar.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Error updating auction status to completed", err)
	}
}

func (ar *AuctionRepository) CloseExpiredAuctions() {
	for {
		time.Sleep(1 * time.Minute)

		now := time.Now().Unix()
		filter := bson.M{"status": auction_entity.Active, "end_time": bson.M{"$lt": now}}
		update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}

		_, err := ar.Collection.UpdateMany(context.Background(), filter, update)
		if err != nil {
			logger.Error("Error updating expired auctions", err)
		}
	}
}
