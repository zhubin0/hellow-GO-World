package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type MongoDB struct {
	Database string
	Sn       *mgo.Session
}

func NewMongo(conf Mongo) *MongoDB {
	var url string = "mongodb://"
	if strings.TrimSpace(conf.Username) != "" {
		url += conf.Username + ":" + conf.Password + "@"
	}
	url += conf.Address
	if strings.TrimSpace(conf.Dbname) != "" {
		url += "/" + conf.Dbname
	}
	//logrus.Error(url)
	sn, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	// enable safe mode and journal sync
	// NOTE do not forget to enable journaling on the mongodb server
	sn.SetSafe(&mgo.Safe{J: true})
	return &MongoDB{
		Database: conf.Dbname,
		Sn:       sn,
	}
}

func (md *MongoDB) shutDown() {
	md.Sn.Close()
}

func (md *MongoDB) Insert(collection string, doc interface{}) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	return c.Insert(doc)
}

func (md *MongoDB) Upsert(collection string, id, doc interface{}) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	_, err := c.Upsert(bson.M{"_id": id}, doc)
	return err
}

// upsert on doc by selector on unique keys
func (md *MongoDB) UpsertOne(collection string, selector interface{}, doc interface{}) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	_, err := c.Upsert(selector, doc)
	return err
}

// Delete a doc.
// NOTE if err == mgo.ErrNotFound, the deletion would fail and irrecoverable
func (md *MongoDB) Delete(collection string, id interface{}) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	return c.RemoveId(id)
}

// Removes all matching docs with selector.
// NOTE be careful here, otherwise you may delete docs unexpectedly
func (md *MongoDB) Remove(collection string, selector interface{}) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	_, err := c.RemoveAll(selector)
	return err
}

// Replace a doc.
// NOTE if err == mgo.ErrNotFound, the replace would faild and irrecoverable
// NOTE the doc to replace must be the document, not the BSON selector.
func (md *MongoDB) Replace(collection, id string, doc interface{}) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	return c.UpdateId(id, doc)
}

// UpdatePassword doc by the given field and it's value.
// NOTE this would only update one doc if exists. If not exists, mgo.ErrNotFound err emerges.
// NOTE if you need to update some of the fields, use "$set" instead of passing the whole struct.
//      if you pass the whole struct, the whole document would be replaced.
func (md *MongoDB) UpdateBy(collection string, field string, fieldValue, change interface{}) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	return c.Update(bson.M{field: fieldValue}, change)
}

// Update all doc by the given field and it's value.
// NOTE this would only update one doc if exists. If not exists, mgo.ErrNotFound err emerges.
// NOTE if you need to update some of the fields, use "$set" instead of passing the whole struct.
//      if you pass the whole struct, the whole document would be replaced.
func (md *MongoDB) UpdateAllBy(collection string, field string, fieldValue, change interface{}) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	_, err := c.UpdateAll(bson.M{field: fieldValue}, change)
	return err
}

// UpdatePassword a doc by id.
// NOTE this would only update one doc if exists. If not exists, mgo.ErrNotFound err emerges.
// NOTE if you need to update some of the fields, use "$set" instead of passing the whole struct.
//      if you pass the whole struct, the whole document would be replaced.
func (md *MongoDB) Update(collection string, id interface{}, change interface{}) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	return c.UpdateId(id, change)
}

// UpdatePassword a doc if your update is complicated
// NOTE make sure use correct selector before you call this function
func (md *MongoDB) UpdateSelfDefined(collection string, selector interface{}, update interface{}) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	return c.Update(selector, update)
}

// Ensure index create an index if not exists.
func (md *MongoDB) EnsureIndex(collection string, index mgo.Index) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	return c.EnsureIndex(index)
}

// Check if the doc exists
func (md *MongoDB) CheckDoc(collection string, id string) (bool, error) {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	var doc = struct {
		Id string `bson:"_id"`
	}{}
	if err := c.Find(bson.M{"_id": id}).Select(bson.M{"_id": 1}).Limit(1).One(&doc); err != nil {
		if err == mgo.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Check if the docs exist
func (md *MongoDB) CheckDocs(collection string, ids []string) ([]bool, error) {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	var docs = []struct {
		Id string `bson:"_id"`
	}{}
	uniq := helperUniqueStrings(ids)
	if err := c.Find(bson.M{"_id": bson.M{"$in": uniq}}).Select(bson.M{"_id": 1}).All(&docs); err != nil {
		return nil, err
	}
	ret := make([]bool, len(ids))
	mp := map[string]bool{}
	for _, doc := range docs {
		mp[doc.Id] = true
	}
	for i, id := range ids {
		if _, ok := mp[id]; ok {
			ret[i] = true
		}
	}
	return ret, nil

}

func helperUniqueStrings(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)
	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}

/*
Example:
	type Person struct {
		ID   string `bson:"_id"`
		Name string `bson:"name"`
		Age  string `bson:"age"`
	}
	p:=Person{}
	if err:=mds.Get("id-to-get",&p);err!=ni;{
		log.Error(err)
	}
	fmt.Println(p)
*/
// Fetch the whole doc according to _id
// NOTE you can use mgo.ErrNotFound to determine where the error is doc-not-found
func (md *MongoDB) Get(collection string, id interface{}, result interface{}) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	return c.FindId(id).One(result)
}

// Get by the given field and fieldValue.
// NOTE you can use mgo.ErrNotFound to determine where the error is doc-not-found
func (md *MongoDB) GetBy(collection string, field string, fieldValue string, result interface{}) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	return c.Find(bson.M{field: fieldValue}).One(result)
}

// Fetch some fields of the doc according to _id.
// NOTE you can use mgo.ErrNotFound to determine where the error is doc-not-found
func (md *MongoDB) GetFields(collection, id string, fields []string, result interface{}) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	return c.FindId(id).Select(selectFields(fields)).One(result)
}

// Find and filter docs. Only support simple filter.
func (md *MongoDB) Find(collection string, filters map[string]string, result interface{}) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	return c.Find(filters).All(result)
}

// Call this method, if fields of document is unique
func (md *MongoDB) FindOne(collection string, filters interface{}, result interface{}) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	return c.Find(filters).One(result)
}

// Call this method, if fields of document is unique
func (md *MongoDB) SortFindOne(collection string, filters interface{}, result interface{}) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	return c.Find(filters).Sort("-_id").One(result)
}

func (md *MongoDB) FindFieldIn(collection string, query interface{}, results interface{}) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	return c.Find(query).All(results)
}

// Find and filter some fields of docs. Only support simple filter.
func (md *MongoDB) FindFields(collection string, filters map[string]string, fields []string, results interface{}) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	return c.Find(filters).Select(selectFields(fields)).All(results)
}

func (md *MongoDB) FindPartFieldIn(collection string, query interface{}, selectFields interface{}, results interface{}) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	return c.Find(query).Select(selectFields).All(results)
}

// Find all and paging.
func (md *MongoDB) PagingFind(collection string, filters interface{}, skip int, limit int, results interface{}) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	return c.Find(filters).Skip(skip).Limit(limit).All(results)
}

// Explain Find all and paging.
func (md *MongoDB) ExplainFind(collection string, filters interface{}, skip int, limit int) (*bson.M, error) {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	res := &bson.M{}
	err := c.Find(filters).Skip(skip).Limit(limit).Explain(res)
	return res, err
}

// Explain Find all and paging, returning only part of the doc
func (md *MongoDB) ExplainFindPartDoc(collection string, filters interface{}, selectFields interface{}, skip int, limit int) (*bson.M, error) {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	res := &bson.M{}
	err := c.Find(filters).Select(selectFields).Skip(skip).Limit(limit).Explain(res)
	return res, err
}

func (md *MongoDB) FindAndSort(collection string, filters interface{}, results interface{}, sortField1 string, sortField2 string) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	return c.Find(filters).Sort(sortField1, sortField2).All(results)
}

func (md *MongoDB) PagingFindAndSort(collection string, filters interface{}, skip int, limit int, results interface{}, sortField string) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	return c.Find(filters).Skip(skip).Limit(limit).Sort(sortField).All(results)
}

func (md *MongoDB) PagingFindAndSortMulti(collection string, filters interface{}, skip int, limit int, results interface{}, sortField1 string, sortField2 string) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	return c.Find(filters).Skip(skip).Limit(limit).Sort(sortField1, sortField2).All(results)
}

// Get current doc count of the given collection
func (md *MongoDB) Count(collection string, filters interface{}) (int, error) {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	return c.Find(filters).Count()
}

func selectFields(q []string) (r bson.M) {
	r = make(bson.M, len(q))
	for _, s := range q {
		r[s] = 1
	}
	return
}

// A thin wrapper for execute pipeline in mgo.
// An example would be to do a 'lookup' operation.
// See: https://gist.github.com/sindbach/efeda4b6b614574ee08229ce4a183882 for example
func (md *MongoDB) Pipe(collection string, pipeline interface{}, results interface{}) error {
	session := md.Sn.Clone()
	defer session.Close()
	c := session.DB(md.Database).C(collection)
	return c.Pipe(pipeline).All(results)
}
