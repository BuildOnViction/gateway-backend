package entity

type AddressDao struct {
	*DAO
}

// NewAddressDao returns a new instance of AddressDao.
func NewAddressDao(dbname string) *AddressDao {
	dbName := dbname
	collectionName := "issued_addresses"
	return &AddressDao{
		DAO: &DAO{
			collectionName: collectionName,
			dbName:         dbName,
		},
	}
}

// func (dao *AddressDao) GetAccountIndex(coin string) (uint64, error) {
// 	indexCounter := db.client.Database(db.config.Database.Name).Collection("account_index")
// 	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
// 	filter := bson.M{"coin": coin}

// 	var result bson.M
// 	err := indexCounter.FindOne(ctx, filter).Decode(&result)

// 	if err != nil {
// 		return 0, err
// 	}

// 	maxIndex, _ := strconv.ParseUint(fmt.Sprintf("%v", result["total"]), 10, 32)
// 	return maxIndex, err
// }

// func (dao *AddressDao) UpdateAccountIndex(coin string, index uint64) error {
// 	indexCounter := db.client.Database(db.config.Database.Name).Collection("account_index")
// 	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
// 	filter := bson.M{"coin": coin}

// 	var result bson.M
// 	err := indexCounter.FindOne(ctx, filter).Decode(&result)

// 	if err != nil || len(result) == 0 {
// 		_, err := indexCounter.InsertOne(ctx, bson.M{"coin": coin, "total": index})
// 		return err
// 	}

// 	indexCounter.FindOneAndUpdate(ctx, filter,
// 		bson.M{
// 			"$set": bson.M{"coin": coin, "total": index},
// 		},
// 	)

// 	return nil
// }

// func (dao *AddressDao) GetAddress(coin string, tomoAddress string) string {
// 	collection := db.client.Database(db.config.Database.Name).Collection("registered_addrs")
// 	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

// 	var result bson.M

// 	collection.FindOne(ctx, bson.M{
// 		"coin": strings.ToUpper(coin),
// 		"tomo": strings.ToLower(tomoAddress),
// 	}).Decode(&result)

// 	return fmt.Sprintf("%v", result["address"])
// }

// func (dao *AddressDao) getTomoAddress(coin string, address string) string {
// 	collection := db.client.Database(db.config.Database.Name).Collection("registered_addrs")
// 	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

// 	var result bson.M

// 	collection.FindOne(ctx, bson.M{
// 		"coin":    strings.ToUpper(coin),
// 		"address": bson.M{"$regex": primitive.Regex{Pattern: address, Options: "i"}},
// 	}).Decode(&result)

// 	return fmt.Sprintf("%v", result["tomo"])
// }
