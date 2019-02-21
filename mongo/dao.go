package mongo

import (
	"log"

	//. "github.com/mlabouardy/movies-restapi/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/prithvisagarrao/restaurant-api/models"
)

type MoviesDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "movies"
)

// Establish a connection to database
func (m *MoviesDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}


// Find list of movies
func (m *MoviesDAO) FindAll() ([]models.Recipe, error) {
	var movies []models.Recipe
	err := db.C(COLLECTION).Find(bson.M{}).All(&movies)
	return movies, err
}

// Find a movie by its id
func (m *MoviesDAO) FindById(id string) (models.Recipe, error) {
	var movie models.Recipe
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&movie)
	return movie, err
}

// Insert a movie into database
func (m *MoviesDAO) Insert(movie models.Recipe) error {
	err := db.C(COLLECTION).Insert(&movie)
	return err
}

// Delete an existing movie
func (m *MoviesDAO) Delete(movie models.Recipe) error {
	err := db.C(COLLECTION).Remove(&movie)
	return err
}

// Update an existing movie
func (m *MoviesDAO) Update(movie models.Recipe) error {
	err := db.C(COLLECTION).UpdateId(movie.RecipeId, &movie)
	return err
}
