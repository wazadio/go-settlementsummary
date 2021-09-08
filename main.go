package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// You will be using this Trainer type later in the program
type Trainer struct {
	Name string
	Age  int
	City string
}

func koneksiDB() *mongo.Client {
	// Rest of the code will go here
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	//fmt.Printf("%T", client)

	return client
}

func create(collection *mongo.Collection, nama string) {
	newData := Trainer{nama, 10, "Pallet Town"}

	insertResult, err := collection.InsertOne(context.TODO(), newData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

func removeDuplicateValues(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func read(collection *mongo.Collection) {
	ctx := context.Background()

	// filterCursor, err := collection.Find(ctx, bson.M{"name": "a"})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer filterCursor.Close(ctx)
	// var episodesFiltered []bson.M
	// if err = filterCursor.All(ctx, &episodesFiltered); err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(filterCursor)
	// //fmt.Println(episodesFiltered[3]["age"])
	// fmt.Printf("%T", episodesFiltered[0])

	// type traceNum struct {
	// 	Tracenum string `json:"age"`
	// }

	// for filterCursor.Next(ctx) {
	// 	var data traceNum
	// 	if err = filterCursor.Decode(&data); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(data)
	// }

	var list_tracenum []string

	//membaca semua dokumen
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	//iterasi dari semua dokumen yang terbaca
	for cursor.Next(ctx) {
		var data bson.D
		if err = cursor.Decode(&data); err != nil {
			log.Fatal(err)
		}
		tanggal := fmt.Sprint(data.Map()["timeStamp"])
		if tanggal[6] == '0' {
			tanggal = string(tanggal[7])
		} else {
			tanggal = tanggal[6:]
		}
		// fmt.Println(tanggal)
		// fmt.Println(fmt.Sprint(time.Now().Day()))
		if tanggal == fmt.Sprint(time.Now().Day()) {
			traceNum := fmt.Sprint(data.Map()["traceNum"])
			list_tracenum = append(list_tracenum, string(traceNum))
		}

		// fmt.Println(traceNum)
		// fmt.Printf("%T\n", traceNum)

		// fmt.Println(data.Map())
		// fmt.Println(data)
		// fmt.Printf("%T", data)
		// fmt.Printf("%T", data.Map()["age"])
		// fmt.Println()
		// fmt.Print(data.Map()["name"], " ", data.Map()["age"], " ", data.Map()["city"])
		// fmt.Println()
	}
	// fmt.Println(list_tracenum)

	// menghilangkan duplicate tracenum dalam list_tracenum
	list_tracenum = removeDuplicateValues(list_tracenum)
	// fmt.Println(list_tracenum)
	for _, i := range list_tracenum {
		fmt.Print(i, " = ")
		cursor, err := collection.Find(ctx, bson.M{"traceNum": i})
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(ctx)
		var pertraceNum []bson.M
		if err = cursor.All(ctx, &pertraceNum); err != nil {
			log.Fatal(err)
		}
		statusIn, err := strconv.ParseInt(fmt.Sprint(pertraceNum[0]["amount"]), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		statusOut, err := strconv.ParseInt(fmt.Sprint(pertraceNum[1]["amount"]), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(statusIn - statusOut)
		// fmt.Printf("%T", pertraceNum[0]["amount"])

		// fmt.Println(pertraceNum[0])
		// for cursor.Next(ctx) {
		// 	var data bson.D
		// 	if err = cursor.Decode(&data); err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	fmt.Println("tracenum= ", i)
		// 	fmt.Println(data)
		// }
	}
}

func update(collection *mongo.Collection, nama string) {
	filter := bson.D{{"name", nama}}

	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}

	/*
		update := bson.D{
			{"$set", bson.D{{nama, nama}, {"age", 15}}},
		}
	*/

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}

func delete(collection *mongo.Collection, nama string) {
	deleteResult, err := collection.DeleteOne(context.TODO(), bson.D{{"name", nama}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
}

func main() {

	pilihan := os.Args[1]

	collection := koneksiDB().Database("testdb").Collection("settlement")

	if pilihan == "c" {
		if len(os.Args) < 3 {
			fmt.Println("Tidak cukup argumen")
		} else {
			create(collection, os.Args[2])
		}
	} else if pilihan == "r" {
		read(collection)
	} else if pilihan == "u" {
		if len(os.Args) < 3 {
			fmt.Println("Tidak cukup argumen")
		} else {
			update(collection, os.Args[2])
		}
	} else if pilihan == "d" {
		if len(os.Args) < 3 {
			fmt.Println("Tidak cukup argumen")
		} else {
			delete(collection, os.Args[2])
		}
	} else {
		println("Perintah tidak diketahui")
	}

	err := koneksiDB().Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

}
