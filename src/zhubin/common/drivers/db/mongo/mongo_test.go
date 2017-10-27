package mongo

import (
	"datamesh.com/common/utils/randgen"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

var collection = "TestCollection"

func initMongoDB() *MongoDB {
	var address = "127.0.0.1:27017"
	var database = "TestDatabase"
	conf := Mongo{address, "", "", database}
	mdb := NewMongo(conf)
	return mdb
}

type Person struct {
	ID          string   `bson:"_id"`
	Name        string   `bson:"name"`
	Age         int      `bson:"age" json:"age"`
	GirlFriends []string `bson:"girlFriends"`
	ThirdGirls  []Girl   `bson:"girl"`
}

type Girl struct {
	Id   int64  `bson:"_id"`
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func TestMongoDB_Upsert(t *testing.T) {
	mdb := initMongoDB()
	defer mdb.shutDown()

	p := Girl{
		Id:   1,
		Name: "world",
	}

	err := mdb.Upsert(collection, p.Id, &p)
	assert.Nil(t, err)

}

// insert/get/find/update/replace/delete doc using given _id
func TestMongoDocStore_flow(t *testing.T) {
	mdb := initMongoDB()
	defer mdb.shutDown()
	var id = randgen.GenUniqueString(8)
	var name = randgen.GenRandString(8)
	p := Person{
		ID:          id,
		Name:        name,
		Age:         18,
		GirlFriends: []string{"xiao", "da", "zhi", "ha"},
	}

	// test insert
	err := mdb.Insert(collection, &p)
	assert.Nil(t, err, err)
	exist, err := mdb.CheckDoc(collection, id)
	assert.Nil(t, err, err)
	assert.True(t, exist, "should exist")

	// test get
	pp := Person{}
	err = mdb.Get(collection, id, &pp)
	assert.Nil(t, err, err)
	assert.Equal(t, p.Name, pp.Name, "name should be the same.")

	// test find
	ps := []Person{}
	m := map[string]string{
		"name": name,
	}
	err = mdb.Find(collection, m, &ps)
	assert.Nil(t, err, err)
	assert.True(t, len(ps) == 1, "only one person should be found")
	assert.True(t, ps[0].Name == p.Name)

	// test update
	name_new := randgen.GenRandString(8)
	change := bson.M{"$set": bson.M{"name": name_new, "age": 24}}
	err = mdb.UpdateBy(collection, "name", name, change)
	assert.Nil(t, err, err)
	np := Person{}
	err = mdb.Get(collection, id, &np)
	assert.Nil(t, err, err)
	assert.Equal(t, name_new, np.Name, "name should be updated.")

	// test update embedded doc
	em_change := bson.M{"$set": bson.M{"girlFriends.0": "biubiu"}}
	err = mdb.Update(collection, id, em_change)
	assert.Nil(t, err, err)
	npp := Person{}
	err = mdb.Get(collection, id, &npp)
	assert.Nil(t, err, err)
	assert.Equal(t, "biubiu", npp.GirlFriends[0], "sub filed should be updated.")

	// test replace
	rp_name := randgen.GenRandString(16)
	pers := Person{ID: id, Name: rp_name}
	err = mdb.Replace(collection, id, pers)
	assert.Nil(t, err, err)
	assert.Equal(t, rp_name, pers.Name, "doc should be replaced.")

	// test delete
	err = mdb.Delete(collection, id)
	assert.Nil(t, err, err)
	still_exist, err := mdb.CheckDoc(collection, id)
	assert.Nil(t, err, err)
	assert.True(t, still_exist == false, "should be deleted.")
}

// test insert duplicate doc
func TestMongoDocStore_insert_duplicate(t *testing.T) {
	mdb := initMongoDB()
	defer mdb.shutDown()
	// first insert
	var id = randgen.GenUniqueString(16)
	var name = randgen.GenRandString(16)
	p := Person{
		ID:   id,
		Name: name,
	}
	err := mdb.Insert(collection, &p)
	assert.Nil(t, err, err)
	exist, err := mdb.CheckDoc(collection, id)
	assert.Nil(t, err, err)
	assert.True(t, exist, "should exist")
	// insert again
	err = mdb.Insert(collection, &p)
	assert.True(t, mgo.IsDup(err), "should complain duplicate error")
	// delete test data
	err = mdb.Delete(collection, id)
	assert.Nil(t, err, err)
	still_exist, err := mdb.CheckDoc(collection, id)
	assert.Nil(t, err, err)
	assert.True(t, still_exist == false, "should be deleted.")
}

// test update non-exist doc
func TestMongoDocStore_update_nonexist_doc(t *testing.T) {
	mdb := initMongoDB()
	defer mdb.shutDown()
	nonexit_id := randgen.GenRandString(20)
	change := bson.M{"$set": bson.M{"name": "aaa", "age": 24}}
	err := mdb.Update(collection, nonexit_id, change)
	assert.True(t, err == mgo.ErrNotFound, "should be ErrNotFound")
}

// test doc replacement
func TestMongoDB_get_flow(t *testing.T) {
	mdb := initMongoDB()
	defer mdb.shutDown()
	id := randgen.GenUniqueString(16)
	name := randgen.GenUniqueString(16)
	age := 18
	ageNew := 19

	p := Person{
		ID:   id,
		Name: name,
		Age:  age,
	}
	err := mdb.Insert(collection, &p)
	assert.Nil(t, err, err)

	exist, err := mdb.CheckDoc(collection, id)
	assert.Nil(t, err, err)
	assert.True(t, exist, "should exist")

	// test get by
	gp := Person{}
	err = mdb.GetBy(collection, "name", name, &gp)
	assert.Nil(t, err, err)
	assert.Equal(t, name, p.Name, "name should be the same.")

	// test update by
	err = mdb.UpdateBy(collection, "name", name, bson.M{"$set": bson.M{"age": ageNew}})
	assert.Nil(t, err, err)

	// test get fields
	fp := Person{}

	err = mdb.GetFields(collection, id, []string{"age"}, &fp)
	assert.Nil(t, err, err)
	assert.True(t, fp.Age == ageNew && fp.Name == "", "should be true.")
}

func TestRedis_Save(t *testing.T) {
	mdb := initMongoDB()
	defer mdb.shutDown()
	for i := 0; i < 20; i++ {
		id := randgen.GenUniqueString(16)
		name := randgen.GenUniqueString(16)
		age := 18
		p := Person{
			ID:   id,
			Name: name,
			Age:  age,
		}
		err := mdb.Insert(collection, &p)
		assert.Nil(t, err, err)
	}
}

func TestMongoDB_FindAll(t *testing.T) {
	mdb := initMongoDB()
	defer mdb.shutDown()
	p := []Person{}

	mdb.Find(collection, nil, &p)
	fmt.Println(len(p))

	for i := range p {
		fmt.Println(p[i].Name)
	}
}

func TestMongoDB_FindFieldIn(t *testing.T) {
	mdb := initMongoDB()
	defer mdb.shutDown()
	p := []Person{}

	m := bson.M{"age": bson.M{"$nin": []int{8}}}

	mdb.FindFieldIn(collection, m, &p)
	fmt.Println(len(p))
}

func TestMongoDB_PagingFind(t *testing.T) {
	mdb := initMongoDB()
	defer mdb.shutDown()
	p := []Person{}

	mdb.PagingFind(collection, nil, 31, 5, &p)

	for i := range p {
		fmt.Println(p[i].Name)
	}
}

func TestMongoDB_Count(t *testing.T) {
	mdb := initMongoDB()
	defer mdb.shutDown()
	fmt.Println(mdb.Count(collection, nil))
}
