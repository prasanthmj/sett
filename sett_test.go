package sett_test

import (
	"encoding/gob"
	"github.com/prasanthmj/sett"
	"os"
	"syreclabs.com/go/faker"
	"testing"
	"time"
)

// instance for testing
var s *sett.Sett

func TestMain(m *testing.M) {
	// set up database for tests

	s = sett.Open(sett.DefaultOptions("./data/jobsdb7"))
	defer s.Close()

	// clean up
	os.RemoveAll("./data/jobsdb7")

	os.Exit(m.Run())
}

func TestSet(t *testing.T) {
	k := faker.RandomString(8)
	v := faker.RandomString(8)
	// should be able to add a key and value
	err := s.SetStr(k, v)
	if err != nil {
		t.Error("Set operation failed:", err)
		return
	}
	vr, err := s.GetStr(k)
	if err != nil {
		t.Error("Get operation failed:", err)
		return
	}
	if vr != v {
		t.Error("Value does not match")
		return
	}

	v2 := faker.RandomString(12)
	err = s.SetStr(k, v2)
	if err != nil {
		t.Error("Set operation failed:", err)
		return
	}
	vr2, err := s.GetStr(k)
	if err != nil {
		t.Error("Get operation failed:", err)
		return
	}
	if vr2 != v2 {
		t.Error("Second retrieval Value does not match")
		return
	}
}

func TestDelete(t *testing.T) {
	k := faker.RandomString(8)
	v := faker.RandomString(8)

	err := s.SetStr(k, v)
	if err != nil {
		t.Error("Set operation failed:", err)
	}

	k2 := faker.RandomString(8)
	v2 := faker.RandomString(8)

	err = s.SetStr(k2, v2)
	if err != nil {
		t.Error("Set operation failed:", err)
		return
	}

	// should be able to delete key
	if err := s.Delete(k); err != nil {
		t.Error("Delete operation failed:", err)
		return
	}

	// key should be gone
	_, err = s.GetStr(k)
	if err == nil {
		t.Error("Key \"key\" should not be found, but it was")
		return
	}

	vr2, err := s.GetStr(k2)
	if err != nil {
		t.Error("Error getting second key", err)
		return
	}
	if vr2 != v2 {
		t.Error("Error second value does not match", err)
		return
	}

}

func TestTableGet(t *testing.T) {
	table := faker.RandomString(8)
	k := faker.RandomString(8)
	v := faker.RandomString(8)

	err := s.Table(table).SetStr(k, v)
	if err != nil {
		t.Error("TableSet operation failed:", err)
		return
	}
	// should be able to retrieve the value of a key
	vr, err := s.Table(table).GetStr(k)
	if err != nil {
		t.Error("TableGet operation failed:", err)
		return
	}
	// val should be "value"
	if v != vr {
		t.Errorf("TableGet: Expected value %s, got %s", v, vr)
	}
}

func TestTableDelete(t *testing.T) {
	table := faker.RandomString(8)
	k := faker.RandomString(8)
	v := faker.RandomString(8)

	err := s.Table(table).SetStr(k, v)
	if err != nil {
		t.Error("TableSet operation failed:", err)
	}
	// should be able to delete key from table
	if err := s.Table(table).Delete(k); err != nil {
		t.Error("Delete operation failed:", err)
	}

	// key should be gone
	_, err = s.Table(table).GetStr(k)
	if err == nil {
		t.Error("Key in table \"table\" should not be found, but it was")
	}
}

func TestKeysFilter(t *testing.T) {
	//Add some random key values first
	for i := 0; i < 15; i++ {
		k := faker.RandomString(12)
		//Make sure the key is unique
		for c := 0; c < 100; c++ {
			if !s.HasKey(k) {
				break
			}
		}
		v := faker.RandomString(22)

		err := s.SetStr(k, v)
		if err != nil {
			t.Error("Set operation failed:", err)
			return
		}
	}

	//Add some keys with specific prefix
	prefix := "prefix_"
	for i := 0; i < 15; i++ {
		k := prefix + faker.RandomString(8)
		//Make sure the key is unique
		for c := 0; c < 100; c++ {
			if !s.HasKey(k) {
				break
			}
		}
		v := faker.RandomString(8)

		err := s.SetStr(k, v)
		if err != nil {
			t.Error("Set operation failed:", err)
			return
		}
	}

	keys, _ := s.Keys(prefix)
	l := len(keys)
	if l != 15 {
		t.Error("Keys expected 15 keys, got", l)
	}
}

func TestDrop(t *testing.T) {
	table := faker.RandomString(8)
	var keys [15]string
	for i := 0; i < 15; i++ {
		k := faker.RandomString(8)
		v := faker.RandomString(8)
		s.Table(table).SetStr(k, v)
		keys[i] = k
	}

	// should be able to delete "table"
	if err := s.Table(table).Drop(); err != nil {
		t.Error("Table Drop, unexpected error", err)
		return
	}

	for i := 0; i < 15; i++ {
		_, err := s.Table(table).GetStr(keys[i])

		if err == nil {
			t.Errorf("Key %s in table \"batch\" should not be found as table droppped, but it was", keys[i])
		}
	}
	// check that a key should be gone
}

func TestTableNameShouldntPersist(t *testing.T) {
	table := faker.RandomString(8)
	k := faker.RandomString(8)
	v := faker.RandomString(8)

	err := s.Table(table).SetStr(k, v)
	if err != nil {
		t.Error("Error setting table value ", err)
		return
	}

	k2 := faker.RandomString(8)
	v2 := faker.RandomString(8)

	err = s.SetStr(k2, v2)
	if err != nil {
		t.Error("Error setting  value ", err)
		return
	}

	s.Table("another-table").SetStr(k, v)

	vr2, err := s.GetStr(k2)
	if err != nil {
		t.Error("Error getting value ", err)
		return
	}
	if vr2 != v2 {
		t.Error("Value does not match for key ", k2)
	}

}

func TestTTL(t *testing.T) {
	table := faker.RandomString(8)

	pk := faker.RandomString(8)
	pv := faker.RandomString(8)

	err := s.Table(table).SetStr(pk, pv)
	if err != nil {
		t.Error("Couldn't set key value", err)
		return
	}

	k := faker.RandomString(8)
	v := faker.RandomString(8)

	mytable := s.Table(table).WithTTL(100 * time.Millisecond)

	err = mytable.SetStr(k, v)
	if err != nil {
		t.Error("Couldn't set key value", err)
		return
	}

	time.Sleep(200 * time.Millisecond)

	_, err = mytable.GetStr(k)
	if err == nil {
		t.Error("Could fetch key value even after TTL expiry")
	}

	vpk, err := mytable.GetStr(pk)
	if err != nil {
		t.Error("Couldn't get key value for permanent key", err)
		return
	}
	if vpk != pv {
		t.Errorf("Couldn't get key value for permanent key. Expected %s Received %s", pv, vpk)
	}
}

type Signup struct {
	Name  string
	Email string
	Age   int
}

func TestSettingStruct(t *testing.T) {
	gob.Register(&Signup{})
	var su Signup
	su.Name = faker.Name().Name()
	su.Email = faker.Internet().SafeEmail()
	su.Age = faker.Number().NumberInt(2)

	k := faker.RandomString(8)

	err := s.Table("signups").SetStruct(k, &su)
	if err != nil {
		t.Error("Error setting struct value ", err)
		return
	}

	sur, err := s.Table("signups").GetStruct(k)
	if err != nil {
		t.Error("Error getting struct value ", err)
	}

	sur2 := sur.(*Signup)

	if sur2.Name != su.Name {
		t.Errorf("The retrieved value does not match %s vs %s", sur2.Name, su.Name)
	}
}

func TestSimpleSet(t *testing.T) {

	k := faker.RandomString(12)
	v := faker.RandomString(12)

	err := s.Set(k, v)
	if err != nil {
		t.Error("Set has thrown error ", err)
		return
	}
	vr, err := s.Get(k)
	if err != nil {
		t.Error("Get has thrown error", err)
		return
	}
	if vr != v {
		t.Errorf("The returned values does not match expected %s got %s", v, vr)
	}
}

type UserSession struct {
	ID    string
	Email string
}

func TestInsert(t *testing.T) {
	gob.Register(&UserSession{})

	session := UserSession{}
	session.ID = faker.RandomString(12)
	session.Email = faker.Internet().Email()

	id, err := s.Table("sessions").Insert(&session)
	if err != nil {
		t.Error("Error inserting a value", err)
		return
	}
	t.Logf("Inserted session ID %s ID length %d", id, len(id))
	sessret, err := s.Table("sessions").GetStruct(id)
	if err != nil {
		t.Error("Error getting inserted value", err)
		return
	}
	session2 := sessret.(*UserSession)
	if session2.ID != session.ID || session2.Email != session.Email {
		t.Error("retrieved session value does not match")
		return
	}

	id, err = s.Table("sessions").WithKeyLength(8).Insert(&session)
	if err != nil {
		t.Error("Error inserting a value", err)
		return
	}
	if len(id) != 8 {
		t.Error("The Id length has no effect")
	}
	t.Logf("Inserted session ID %s ID length %d", id, len(id))

}

func TestInsertWithExpiry(t *testing.T) {
	gob.Register(&UserSession{})

	session := UserSession{}
	session.ID = faker.RandomString(12)
	session.Email = faker.Internet().Email()

	id, err := s.Table("sessions").WithTTL(200 * time.Millisecond).Insert(&session)

	time.Sleep(300 * time.Millisecond)

	_, err = s.Table("sessions").GetStruct(id)
	if err == nil {
		t.Error("Expiry is not working for Insert")
	}
}

func TestGetKeys(t *testing.T) {
	gob.Register(&UserSession{})

	table := faker.RandomString(12)

	for i := 0; i < 15; i++ {
		session := UserSession{}
		session.ID = faker.RandomString(12)
		session.Email = faker.Internet().Email()

		k, err := s.Table(table).Insert(&session)
		if err != nil {
			t.Error("Error inserting new items ", err)
			return
		}
		t.Logf("key: %s", k)
	}

	keys, err := s.Table(table).Keys()
	if err != nil {
		t.Error("Error Getting item keys ", err)
		return
	}
	if len(keys) != 15 {
		t.Errorf("Expected 15 keys got %d", len(keys))
	}
	//t.Logf("Received keys %v ", keys)
	for _, k := range keys {
		t.Logf("key %s ", k)
		it, err := s.Table(table).GetStruct(k)
		if err != nil {
			t.Errorf("Error getting item with key %s : %v ", k, err)
			return
		}
		sess := it.(*UserSession)
		t.Logf("retrieved session obj %v ", sess)
	}
}

func TestCutting(t *testing.T) {
	gob.Register(&Signup{})
	var su Signup
	su.Name = faker.Name().Name()
	su.Email = faker.Internet().SafeEmail()
	su.Age = faker.Number().NumberInt(2)

	k := faker.RandomString(8)

	err := s.Table("signups").SetStruct(k, &su)
	if err != nil {
		t.Error("Error setting struct value ", err)
		return
	}
	sur, err := s.Table("signups").Cut(k)
	if err != nil {
		t.Error("Error cutting struct value ", err)
		return
	}

	sur2 := sur.(*Signup)

	if sur2.Name != su.Name {
		t.Errorf("The retrieved value does not match %s vs %s", sur2.Name, su.Name)
	}

	_, err = s.Table("signups").GetStruct(k)
	if err == nil {
		t.Error("The item can be retrieved even after cutting it")
	}
}

func TestCuttingWithInsert(t *testing.T) {
	gob.Register(&Signup{})
	table := faker.RandomString(8)
	var su Signup
	su.Name = faker.Name().Name()
	su.Email = faker.Internet().SafeEmail()
	su.Age = faker.Number().NumberInt(2)

	k, err := s.Table(table).Insert(&su)
	if err != nil {
		t.Error("Error inserting struct value ", err)
		return
	}

	sur, err := s.Table(table).Cut(k)
	if err != nil {
		t.Error("Error cutting struct value ", err)
		return
	}

	sur2 := sur.(*Signup)

	if sur2.Name != su.Name {
		t.Errorf("The retrieved value does not match %s vs %s", sur2.Name, su.Name)
	}

	_, err = s.Table(table).GetStruct(k)
	if err == nil {
		t.Error("The item can be retrieved even after cutting it")
	}
}

type Item struct {
	Color string
	Name  string
}

func TestFilterFunc(t *testing.T) {
	gob.Register(&Item{})
	table := faker.RandomString(8)
	var itm1 Item
	itm1.Color = "green"
	itm1.Name = faker.RandomString(12)
	s.Table(table).Insert(&itm1)

	var itm2 Item
	itm2.Color = "red"
	itm2.Name = faker.RandomString(12)
	s.Table(table).Insert(&itm2)

	var itm3 Item
	itm3.Color = "green"
	itm3.Name = faker.RandomString(12)
	s.Table(table).Insert(&itm3)

	keys, err := s.Table(table).Filter(func(k string, i interface{}) bool {
		it := i.(*Item)
		if it.Color == "red" {
			return true
		}
		return false
	})
	if err != nil {
		t.Errorf("Error running filter %v", err)
		return
	}
	if len(keys) != 1 {
		t.Errorf("Filter didn't find the right keys. got this: %v", keys)
		return
	}
	i2, err := s.Table(table).GetStruct(keys[0])
	if err != nil {
		t.Errorf("Error getting item %s ", keys[0])
		return
	}
	it2 := i2.(*Item)

	if it2.Name != itm2.Name {
		t.Errorf("Filter retrieval is incorrect expected %s received %s", itm2.Name, it2.Name)
	}
}
